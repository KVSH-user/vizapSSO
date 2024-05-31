package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"vizapSSO/internal/entity"
	"vizapSSO/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, password, dbName string) (*Storage, error) {
	const op = "postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{db: db}

	err = goose.Up(storage.db, "migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
}

func (s *Storage) SaveUser(phone string, passHash []byte) (int64, error) {
	const op = "postgres.SaveUser"

	query := `
		INSERT INTO users (phone, password_hashed) 
		VALUES ($1, $2) 
		RETURNING id;
		`

	var id int64

	err := s.db.QueryRow(query, phone, passHash).Scan(&id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return 0, storage.ErrUserExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) ProvideUser(phone string) (entity.User, error) {
	const op = "postgres.ProvideUser"

	query := `
		SELECT users.id,
		users.password_hashed
		FROM users
		WHERE phone = $1
		LIMIT 1;
		`
	var user entity.User

	err := s.db.QueryRow(query, phone).Scan(&user.ID, &user.PassHash)
	if err == sql.ErrNoRows {
		return user, storage.ErrUserNotFound
	} else if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) SaveRefreshToken(refreshToken string, uid int64) error {
	const op = "postgres.SaveRefreshToken"

	secondQuery := `
		INSERT INTO refresh_token (token, uid, is_active)
		VALUES ($1, $2, true);
		`

	query := `
		UPDATE refresh_token
		SET is_active = false
		WHERE uid = $1;
		`

	tx, err := s.db.Begin()

	_, err = tx.Exec(query, uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(secondQuery, refreshToken, uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) App(appID int32) (entity.App, error) {
	const op = "postgres.App"

	query := `
		SELECT apps.name,
		apps.secret,
		apps.ID
		FROM apps
		WHERE id = $1
		LIMIT 1;
		`
	var app entity.App

	err := s.db.QueryRow(query, appID).Scan(&app.Name, &app.Secret, &app.ID)
	if err == sql.ErrNoRows {
		return app, storage.ErrAppNotFound
	} else if err != nil {
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) CheckRefreshToken(refreshToken string) error {
	const op = "postgres.CheckRefreshToken"

	query := `
		SELECT refresh_token.is_active
		FROM refresh_token
		WHERE token = $1
		LIMIT 1;
		`

	var isValid bool

	err := s.db.QueryRow(query, refreshToken).Scan(&isValid)
	if err == sql.ErrNoRows {
		return storage.ErrInvalidRefreshToken
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if isValid != true {
		return fmt.Errorf("error %w", storage.ErrInvalidRefreshToken)
	}

	return nil
}
