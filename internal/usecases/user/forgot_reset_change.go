package user

import (
	"context"
	repositories "doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/mailer"
	"doan/internal/services/security"
	"doan/pkg/config"
	"doan/pkg/types"
	"doan/pkg/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ForgotPassword

type ForgotPasswordInput struct {
	Email string
}

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, in ForgotPasswordInput) error
}

type forgotPasswordUseCase struct {
	userRepo repositoryinterface.UserRepository
	mailer   mailer.Mailer
	cfg      config.Manager
}

func NewForgotPasswordUseCase(repo repositoryinterface.UserRepository, mailer mailer.Mailer, cfg config.Manager) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{userRepo: repo, mailer: mailer, cfg: cfg}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, in ForgotPasswordInput) error {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" {
		return errors.New("invalid email")
	}
	cond := repositories.NewCommonCondition()
	cond.AddCondition("email", email, repositories.Equal)
	p, _ := u.userRepo.GetByCondition(ctx, cond)
	if p == nil || len(p.Data) == 0 {
		// Do not reveal existence; still return success
		return nil
	}
	user := p.Data[0]
	// Short-lived reset token via JWT to avoid new table in this phase
	appCfg := types.JWTConfig{}
	if err := u.cfg.UnmarshalKey("jwt", &appCfg); err != nil {
		return err
	}
	// Default 30 minutes for reset token
	ttl := 30 * time.Minute
	// Could also read from app.reset_token_ttl_min
	if mins := getInt(u.cfg, "app.reset_token_ttl_min", 0); mins > 0 {
		ttl = time.Duration(mins) * time.Minute
	}
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, appCfg.Secret, ttl)
	if err != nil {
		return err
	}
	baseURL := getString(u.cfg, "app.base_url", "")
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimRight(baseURL, "/"), token)
	html := fmt.Sprintf("<p>Xin chào %s,</p><p>Nhấp vào liên kết để đặt lại mật khẩu (hết hạn sau %d phút):</p><p><a href=\"%s\">Đặt lại mật khẩu</a></p>", user.FullName, int(ttl.Minutes()), resetLink)
	return u.mailer.Send(ctx, mailer.Mail{To: user.Email, Subject: "Đặt lại mật khẩu", HTML: html})
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
	userRepo repositoryinterface.UserRepository
	hasher   security.PasswordHasher
	cipher   security.PasswordCipher
	cfg      config.Manager
}

func NewResetPasswordUseCase(repo repositoryinterface.UserRepository, hasher security.PasswordHasher, cipher security.PasswordCipher, cfg config.Manager) ResetPasswordUseCase {
	return &resetPasswordUseCase{userRepo: repo, hasher: hasher, cipher: cipher, cfg: cfg}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, in ResetPasswordInput) error {
	if strings.TrimSpace(in.Token) == "" || strings.TrimSpace(in.NewPasswordEnc) == "" {
		return errors.New("invalid payload")
	}
	jwtCfg := types.JWTConfig{}
	if err := u.cfg.UnmarshalKey("jwt", &jwtCfg); err != nil {
		return err
	}
	claims, err := utils.ValidateToken(in.Token, jwtCfg.Secret)
	if err != nil {
		return errors.New("invalid or expired token")
	}
	plain, err := u.cipher.Decrypt(in.NewPasswordEnc)
	if err != nil {
		return err
	}
	hash, err := u.hasher.Hash(plain)
	if err != nil {
		return err
	}
	update := map[string]interface{}{"password": hash}
	return u.userRepo.Update(ctx, claims.UserID, update)
}

// ChangePassword

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
