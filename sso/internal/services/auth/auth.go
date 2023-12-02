package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/empfaze/golang_grpc/sso/internal/domain/models"
	"github.com/empfaze/golang_grpc/sso/internal/jwt"
	"github.com/empfaze/golang_grpc/sso/internal/storage"
)

type Auth struct {
	logger       *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int) (models.App, error)
}

const OPERATION_TRACE_REGISTER = "auth.Register"
const OPERATION_TRACE_LOGIN = "auth.Login"
const OPERATION_TRACE_IS_ADMIN = "auth.IsAdmin"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserExists = errors.New("user already exists")

func New(
	logger *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		logger:       logger,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	logger := a.logger.With(
		slog.String("op", OPERATION_TRACE_REGISTER),
		slog.String("email", email),
	)

	logger.Info("Logging user")

	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.logger.Warn("User not found", err)

			return "", fmt.Errorf("%s: %w", OPERATION_TRACE_LOGIN, ErrInvalidCredentials)
		}

		a.logger.Error("Failed to get user", err)

		return "", fmt.Errorf("%s: %w", OPERATION_TRACE_LOGIN, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logger.Info("Invalid credentials", err)

		return "", fmt.Errorf("%s: %w", OPERATION_TRACE_LOGIN, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", OPERATION_TRACE_LOGIN, err)
	}

	logger.Info("User logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.logger.Error("Failed to generate token", err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	logger := a.logger.With(
		slog.String("op", OPERATION_TRACE_REGISTER),
		slog.String("email", email),
	)

	logger.Info("Registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to generate password hash", err)

		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_REGISTER, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			logger.Warn("User already exists", err)

			return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_REGISTER, ErrUserExists)
		}

		logger.Error("Failed to save user", err)

		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_REGISTER, err)
	}

	logger.Info("User has been registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	logger := a.logger.With(
		slog.String("op", OPERATION_TRACE_IS_ADMIN),
		slog.Int64("user_id", userID),
	)

	logger.Info("Checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", OPERATION_TRACE_IS_ADMIN, err)
	}

	logger.Info("Checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
