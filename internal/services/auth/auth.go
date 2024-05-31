package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	"vizapSSO/internal/entity"
	"vizapSSO/internal/lib/jwt"
	"vizapSSO/internal/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log                 *slog.Logger
	usrSaver            UserSaver
	appProvider         AppProvider
	refreshTokenSaver   RefreshTokenSaver
	refreshTokenChecker RefreshTokenChecker
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
	userProvider        UserProvider
}

type UserSaver interface {
	SaveUser(phone string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	ProvideUser(phone string) (entity.User, error)
}

type AppProvider interface {
	App(appID int32) (entity.App, error)
}

type RefreshTokenSaver interface {
	SaveRefreshToken(refreshToken string, uid int64) error
}

type RefreshTokenChecker interface {
	CheckRefreshToken(refreshToken string) error
}

func New(log *slog.Logger,
	userSaver UserSaver,
	appProvider AppProvider,
	userProvider UserProvider,
	refreshTokenSaver RefreshTokenSaver,
	refreshTokenChecker RefreshTokenChecker,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:            userSaver,
		appProvider:         appProvider,
		userProvider:        userProvider,
		refreshTokenSaver:   refreshTokenSaver,
		refreshTokenChecker: refreshTokenChecker,
		accessTokenTTL:      accessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
		log:                 log,
	}
}

func (a *Auth) Login(ctx context.Context, phone, password string, appID int32,
) (accessToken, refreshToken string, err error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("login attempt")

	user, err := a.userProvider.ProvideUser(phone)
	if err != nil {
		log.Error("failed to provide user", err)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if user.ID == 0 {
		log.Error("phone not found", storage.ErrUserNotFound)
		return "", "", fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(appID)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in success")

	accessToken, err = jwt.NewAccessToken(user, app, a.accessTokenTTL)
	if err != nil {
		a.log.Error("failed to generate access token", err)

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	lastChar := accessToken[len(accessToken)-6:]

	refreshToken, err = jwt.NewRefreshToken(lastChar, app, a.refreshTokenTTL)
	if err != nil {
		a.log.Error("failed to generate refresh token", err)

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := a.refreshTokenSaver.SaveRefreshToken(refreshToken, user.ID); err != nil {
		a.log.Error("failed to save refresh token", err)
	}

	return accessToken, refreshToken, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, phone, password string,
) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op))

	log.Info("registering user")

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(phone, passwordHashed)
	if err != nil {
		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully register")

	return id, nil
}

func (a *Auth) ValidateSession(ctx context.Context, accessToken string) (isValid bool, uid int64, err error) {
	const op = "auth.ValidateSession"

	log := a.log.With(slog.String("op", op))

	log.Info("validate user token")

	isValid, uid, err = jwt.ValidateToken(accessToken)
	if err != nil {
		log.Error("failed validate token", err)
		return false, 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("token successfully validate")

	return isValid, uid, nil
}

func (a *Auth) RefreshSession(ctx context.Context, accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	const op = "auth.RefreshSession"

	log := a.log.With(slog.String("op", op))

	log.Info("refresh user token")

	err = jwt.CheckRefreshToken(refreshToken, accessToken)
	if err != nil {
		log.Error("failed to validate token pair", err)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	uid, err := jwt.IdFromJWT(accessToken)

	err = a.refreshTokenChecker.CheckRefreshToken(refreshToken)
	if err != nil {
		a.log.Error("invalid refresh token", err)
		return "", "", err
	}

	var user entity.User
	var app entity.App

	//TODO: исправить данную реализацию
	app.ID = 1
	user.ID = uid

	newAccessToken, err = jwt.NewAccessToken(user, app, a.accessTokenTTL)
	if err != nil {
		a.log.Error("failed to generate access token", err)

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	lastChar := accessToken[len(newAccessToken)-6:]

	newRefreshToken, err = jwt.NewRefreshToken(lastChar, app, a.refreshTokenTTL)
	if err != nil {
		a.log.Error("failed to generate refresh token", err)

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := a.refreshTokenSaver.SaveRefreshToken(newRefreshToken, user.ID); err != nil {
		a.log.Error("failed to save refresh token", err)
	}

	log.Info("user token successfully refreshed")

	return newAccessToken, newRefreshToken, nil
}

func (a *Auth) RequestPasswordReset(ctx context.Context, email string) (response string, err error) {
	panic("implement me")
}

func (a *Auth) PerformPasswordReset(ctx context.Context, token, newPassword string) (success bool, err error) {
	panic("implement me")
}
