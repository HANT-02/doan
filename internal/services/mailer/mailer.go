package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net/smtp"
	"strings"

	_interface "doan/internal/infrastructure/queue/interface"
	"doan/pkg/config"
	"doan/pkg/logger"
)

type Mail struct {
	To      string
	Subject string
	HTML    string
}

type Mailer interface {
	Send(ctx context.Context, m Mail) error
	SendOTPEmail(ctx context.Context, toEmail string, otp string) error
}

type queueMailer struct {
	q          _interface.Queue
	log        logger.Logger
	config     config.Manager
	kafkaTopic string
}

func NewMailer(q _interface.Queue, log logger.Logger, cfg config.Manager) Mailer {
	// Some defaults; in future read from cfg (smtp or mq)
	return &queueMailer{
		q:          q,
		log:        log,
		config:     cfg,
		kafkaTopic: "emails",
	}
}

func (m *queueMailer) Send(ctx context.Context, mail Mail) error {
	// Send email immediately via SMTP (queue can be reintroduced later via config)
	// Read SMTP configuration
	var host, username, password, from, appEnv string
	var port int
	var useTLS, insecureSkipVerify bool

	_ = m.config.UnmarshalKey("smtp.host", &host)
	_ = m.config.UnmarshalKey("smtp.port", &port)
	_ = m.config.UnmarshalKey("smtp.username", &username)
	_ = m.config.UnmarshalKey("smtp.password", &password)
	_ = m.config.UnmarshalKey("smtp.from", &from)
	_ = m.config.UnmarshalKey("smtp.tls", &useTLS)
	_ = m.config.UnmarshalKey("smtp.insecure_skip_verify", &insecureSkipVerify)
	_ = m.config.UnmarshalKey("app.env", &appEnv)

	if from == "" {
		from = username // fallback
	}

	// If SMTP host is not configured
	if host == "" || port == 0 {
		// In dev/local, log as if sent to avoid breaking flows
		switch strings.ToLower(strings.TrimSpace(appEnv)) {
		case "dev", "development", "local":
			m.log.Info(ctx, "DEV-MAIL SMTP not configured, printing email", "to", mail.To, "subject", mail.Subject)
			m.log.Info(ctx, "DEV-MAIL body", "html", mail.HTML)
			return nil
		}
		return fmt.Errorf("smtp is not configured: host/port missing")
	}

	// Build MIME message
	mime := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n"
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s",
		from, mail.To, encodeRFC2047(mail.Subject), mime, mail.HTML))

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	// Decide TLS strategy
	if useTLS {
		// Try STARTTLS if possible; if server requires implicit TLS (465), use TLS dial
		if port == 465 {
			// Implicit TLS
			tlsConfig := &tls.Config{ServerName: host, InsecureSkipVerify: insecureSkipVerify}
			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				return fmt.Errorf("tls dial failed: %w", err)
			}
			c, err := smtp.NewClient(conn, host)
			if err != nil {
				return fmt.Errorf("smtp new client failed: %w", err)
			}
			defer c.Quit()
			if username != "" {
				if err := c.Auth(auth); err != nil {
					return fmt.Errorf("smtp auth failed: %w", err)
				}
			}
			if err := c.Mail(from); err != nil {
				return err
			}
			if err := c.Rcpt(mail.To); err != nil {
				return err
			}
			w, err := c.Data()
			if err != nil {
				return err
			}
			if _, err := w.Write(msg); err != nil {
				_ = w.Close()
				return err
			}
			if err := w.Close(); err != nil {
				return err
			}
			m.log.Info(ctx, "Email sent via SMTP (implicit TLS)", "to", mail.To, "addr", addr)
			return nil
		}
		// STARTTLS path (e.g., port 587)
		c, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("smtp dial failed: %w", err)
		}
		defer c.Quit()
		if ok, _ := c.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{ServerName: host, InsecureSkipVerify: insecureSkipVerify}
			if err := c.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("starttls failed: %w", err)
			}
		}
		if username != "" {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}
		if err := c.Mail(from); err != nil {
			return err
		}
		if err := c.Rcpt(mail.To); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write(msg); err != nil {
			_ = w.Close()
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		m.log.Info(ctx, "Email sent via SMTP (STARTTLS)", "to", mail.To, "addr", addr)
		return nil
	}

	// Plain SMTP (no TLS) - use only for local/dev mail servers
	if err := smtp.SendMail(addr, auth, from, []string{mail.To}, msg); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}
	m.log.Info(ctx, "Email sent via SMTP (plain)", "to", mail.To, "addr", addr)
	return nil
}

func (m *queueMailer) SendOTPEmail(ctx context.Context, toEmail string, otp string) error {
	subject := "Your OTP Verification Code"

	htmlBody := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif;">
			<h3>Email Verification</h3>
			<p>Your OTP code is:</p>
			<p style="font-size: 20px; font-weight: bold;">%s</p>
			<p>This code will expire in <strong>5 minutes</strong>.</p>
			<p>If you did not request this, please ignore this email.</p>
		</body>
		</html>
	`, otp)

	mail := Mail{
		To:      toEmail,
		Subject: subject,
		HTML:    htmlBody,
	}

	return m.Send(ctx, mail)
}

// encodeRFC2047 encodes strings for use in mail headers per RFC 2047 (for non-ASCII subjects)
func encodeRFC2047(str string) string {
	// use mail's RFC 2047 to encode any string
	addr := mime.QEncoding.Encode("UTF-8", str)
	return addr
}
