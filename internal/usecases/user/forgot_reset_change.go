package user

import (
	"context"
	"doan/internal/entities"
	repositories "doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/mailer"
	"doan/internal/services/security"
	"doan/pkg/config"
	"doan/pkg/random" // New import for random token generation
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm" // New import for transaction
)

// ForgotPassword

type ForgotPasswordInput struct {
	Email string
	// Request information for audit logging (optional)
	RequestIP string
	UserAgent string
}

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, in ForgotPasswordInput) error
}

type forgotPasswordUseCase struct {
	userRepo          repositoryinterface.UserRepository
	passwordResetRepo repositoryinterface.PasswordResetRepository // New dependency
	mailer            mailer.Mailer
	cfg               config.Manager
	hasher            security.PasswordHasher // New dependency for hashing token
	db                *gorm.DB                // New dependency for transaction
}

func NewForgotPasswordUseCase(
	repo repositoryinterface.UserRepository,
	passwordResetRepo repositoryinterface.PasswordResetRepository, // New
	mailer mailer.Mailer,
	cfg config.Manager,
	hasher security.PasswordHasher, // New
	db *gorm.DB, // New
) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		userRepo:          repo,
		passwordResetRepo: passwordResetRepo,
		mailer:            mailer,
		cfg:               cfg,
		hasher:            hasher,
		db:                db,
	}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, in ForgotPasswordInput) error {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" {
		return errors.New("invalid email")
	}

	cond := repositories.NewCommonCondition()
	cond.AddCondition("email", email, repositories.Equal)
	userResult, err := u.userRepo.GetByCondition(ctx, cond)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Log actual database errors
		// u.logger.Error(ctx, "Failed to get user by email", "error", err) // Assuming logger is available
		return err
	}

	var user *entities.User
	if userResult != nil && len(userResult.Data) > 0 {
		user = userResult.Data[0]
	}

	// Always return success to avoid user enumeration attacks, but only send email if user exists
	if user == nil {
		// u.logger.Info(ctx, "Forgot password request for non-existent email", "email", email) // Assuming logger is available
		return nil
	}

	// Generate cryptographically secure random token
	plainToken := random.GenerateSixDigitOtp() // 32 chars long for good randomness
	tokenHash, err := u.hasher.Hash(plainToken)
	if err != nil {
		return fmt.Errorf("failed to hash reset token: %w", err)
	}

	// Get TTL from config, default to 15 minutes
	resetTokenTTLMinutes := getInt(u.cfg, "auth.reset_token_ttl_minutes", 15)
	expiresAt := time.Now().Add(time.Duration(resetTokenTTLMinutes) * time.Minute)

	// Save password reset record in a transaction
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Optional: Invalidate any previous pending reset tokens for this user
		// This prevents a user from having multiple active reset links
		// For simplicity, we'll just create a new one. The old one will expire.

		passwordReset := &entities.PasswordReset{
			UserID:      user.ID,
			TokenHash:   tokenHash,
			ExpiresAt:   expiresAt,
			RequestedIP: in.RequestIP,
			UserAgent:   in.UserAgent,
		}
		if err := u.passwordResetRepo.CreateTx(ctx, tx, passwordReset); err != nil {
			return fmt.Errorf("failed to create password reset record: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Get frontend reset URL from config
	frontendResetURL := getString(u.cfg, "app.frontend_reset_url", "http://localhost:5173/reset-password") // Default for local dev
	if frontendResetURL == "" {
		return errors.New("frontend reset URL not configured")
	}

	resetLink := fmt.Sprintf("%s?token=%s", strings.TrimRight(frontendResetURL, "/"), plainToken)

	// Construct and send email asynchronously
	go func(toEmail, resetLink string, ttl int) {
		subject := "Yêu cầu đặt lại mật khẩu của bạn"
		htmlBody := fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="vi">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Đặt lại mật khẩu</title>
				<style>
					body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
					.container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
					.button { display: inline-block; background-color: #007bff; color: #ffffff !important; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
					.footer { margin-top: 20px; font-size: 0.9em; color: #777; }
				</style>
			</head>
			<body>
				<div class="container">
					<h2>Yêu cầu đặt lại mật khẩu</h2>
					<p>Xin chào,</p>
					<p>Chúng tôi đã nhận được yêu cầu đặt lại mật khẩu cho tài khoản của bạn.</p>
					<p>Vui lòng nhấp vào liên kết dưới đây để đặt lại mật khẩu của bạn:</p>
					<p><a href="%s" class="button">Đặt lại mật khẩu</a></p>
					<p>Liên kết này sẽ hết hạn sau <strong>%d phút</strong>.</p>
					<p>Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này.</p>
					<div class="footer">
						<p>Trân trọng,</p>
						<p>Đội ngũ Ứng dụng của bạn</p>
					</div>
				</div>
			</body>
			</html>
		`, resetLink, ttl)

		mail := mailer.Mail{
			To:      toEmail,
			Subject: subject,
			HTML:    htmlBody,
		}

		if sendErr := u.mailer.Send(ctx, mail); sendErr != nil {
			// u.logger.Error(ctx, "Failed to send password reset email", "error", sendErr, "email", toEmail) // Assuming logger is available
			fmt.Printf("Failed to send password reset email to %s: %v\n", toEmail, sendErr) // Temp print
		}
	}(email, resetLink, resetTokenTTLMinutes)

	return nil
}

// ResetPassword

type ResetPasswordInput struct {
	Token          string
	NewPasswordEnc string
}

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, in ResetPasswordInput) error
}

type resetPasswordUseCase struct {
	userRepo          repositoryinterface.UserRepository
	passwordResetRepo repositoryinterface.PasswordResetRepository // New dependency
	hasher            security.PasswordHasher
	cipher            security.PasswordCipher
	cfg               config.Manager
	db                *gorm.DB // New dependency for transaction
}

func NewResetPasswordUseCase(
	repo repositoryinterface.UserRepository,
	passwordResetRepo repositoryinterface.PasswordResetRepository, // New
	hasher security.PasswordHasher,
	cipher security.PasswordCipher,
	cfg config.Manager,
	db *gorm.DB, // New
) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		userRepo:          repo,
		passwordResetRepo: passwordResetRepo,
		hasher:            hasher,
		cipher:            cipher,
		cfg:               cfg,
		db:                db,
	}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, in ResetPasswordInput) error {
	if strings.TrimSpace(in.Token) == "" || strings.TrimSpace(in.NewPasswordEnc) == "" {
		return errors.New("invalid payload")
	}

	// Hash the input token to compare with stored hash
	inputTokenHash, err := u.hasher.Hash(in.Token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	var foundReset *entities.PasswordReset

	// Validate and mark token as used in a transaction
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find reset record by hashed token.
		// The GetByTokenHash method already checks expiration and used_at status.
		reset, err := u.passwordResetRepo.GetByTokenHash(ctx, inputTokenHash)
		if err != nil {
			return fmt.Errorf("failed to get password reset record: %w", err)
		}
		if reset == nil || reset.UsedAt != nil { // Explicitly check used_at here as well
			return errors.New("invalid or expired token")
		}

		// Mark token as used
		if err := u.passwordResetRepo.MarkAsUsedTx(ctx, tx, reset.ID, time.Now()); err != nil {
			return fmt.Errorf("failed to mark password reset token as used: %w", err)
		}
		foundReset = reset
		return nil
	})

	if err != nil {
		return err // This will contain "invalid or expired token" or internal errors
	}

	// Decrypt and hash new password
	plainNewPassword, err := u.cipher.Decrypt(in.NewPasswordEnc)
	if err != nil {
		return errors.New("invalid new password payload")
	}

	// TODO: Add password strength validation here
	if len(plainNewPassword) < 8 { // Example minimum length
		return errors.New("new password is too short")
	}

	hashedNewPassword, err := u.hasher.Hash(plainNewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update user's password
	update := map[string]interface{}{"password": hashedNewPassword}
	if err := u.userRepo.Update(ctx, foundReset.UserID, update); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	// TODO: Consider revoking all sessions/refresh tokens for the user here if such a mechanism exists.
	// The current JWT-based system relies on token expiration.
	// If refresh tokens are stored in DB, they should be revoked.

	return nil
}

// ChangePassword (no changes needed)

type ChangePasswordInput struct {
	UserID         string
	OldPasswordEnc string
	NewPasswordEnc string
}

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, in ChangePasswordInput) error
}

type changePasswordUseCase struct {
	userRepo repositoryinterface.UserRepository
	hasher   security.PasswordHasher
	cipher   security.PasswordCipher
}

func NewChangePasswordUseCase(repo repositoryinterface.UserRepository, hasher security.PasswordHasher, cipher security.PasswordCipher) ChangePasswordUseCase {
	return &changePasswordUseCase{userRepo: repo, hasher: hasher, cipher: cipher}
}

func (u *changePasswordUseCase) Execute(ctx context.Context, in ChangePasswordInput) error {
	if strings.TrimSpace(in.UserID) == "" || in.OldPasswordEnc == "" || in.NewPasswordEnc == "" {
		return errors.New("invalid payload")
	}
	user, err := u.userRepo.GetByID(ctx, in.UserID)
	if err != nil {
		return errors.New("user not found")
	}
	oldPlain, err := u.cipher.Decrypt(in.OldPasswordEnc)
	if err != nil {
		return err
	}
	if err := u.hasher.Compare(user.Password, oldPlain); err != nil {
		return errors.New("old password incorrect")
	}
	newPlain, err := u.cipher.Decrypt(in.NewPasswordEnc)
	if err != nil {
		return err
	}
	hash, err := u.hasher.Hash(newPlain)
	if err != nil {
		return err
	}
	update := map[string]interface{}{"password": hash}
	return u.userRepo.Update(ctx, in.UserID, update)
}

// helpers to read config safely
func getString(cfg config.Manager, key, def string) string {
	var v string
	if err := cfg.UnmarshalKey(key, &v); err != nil || v == "" {
		return def
	}
	return v
}

func getInt(cfg config.Manager, key string, def int) int {
	var v int
	if err := cfg.UnmarshalKey(key, &v); err != nil || v == 0 {
		return def
	}
	return v
}
