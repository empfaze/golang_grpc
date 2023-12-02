package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/empfaze/golang_grpc/sso/internal/domain/models"
	"github.com/empfaze/golang_grpc/sso/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

const OPERATION_TRACE_NEW = "storage.sqlite.New"
const OPERATION_TRACE_SAVE_USER = "storage.sqlite.SaveUser"
const OPERATION_TRACE_GET_USER = "storage.sqlite.GetUser"
const OPERATION_TRACE_GET_APP = "storage.sqlite.App"
const OPERATION_TRACE_IS_ADMIN = "storage.sqlite.IsAdmin"

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", OPERATION_TRACE_NEW, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	query, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_SAVE_USER, err)
	}

	result, err := query.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_SAVE_USER, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_SAVE_USER, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_SAVE_USER, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	query, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", OPERATION_TRACE_GET_USER, err)
	}

	row := query.QueryRowContext(ctx, email)

	var user models.User

	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", OPERATION_TRACE_GET_USER, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", OPERATION_TRACE_GET_USER, err)
	}

	return user, nil
}

//func (s *Storage) SavePermission(ctx context.Context, userID int64, permission models.Permission, appID string) error {
//	const op = "storage.sqlite.SavePermission"
//
//	stmt, err := s.db.Prepare("INSERT INTO permissions(user_id, permission, app_id) VALUES(?, ?, ?)")
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	_, err = stmt.ExecContext(ctx, userID, permission, appID)
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	return nil
//}

func (s *Storage) GetApp(ctx context.Context, id int) (models.App, error) {
	query, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", OPERATION_TRACE_GET_APP, err)
	}

	row := query.QueryRowContext(ctx, id)

	var app models.App

	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", OPERATION_TRACE_GET_APP, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", OPERATION_TRACE_GET_APP, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	query, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", OPERATION_TRACE_IS_ADMIN, err)
	}

	row := query.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", OPERATION_TRACE_IS_ADMIN, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", OPERATION_TRACE_IS_ADMIN, err)
	}

	return isAdmin, nil
}
