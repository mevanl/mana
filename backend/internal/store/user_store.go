package store

import (
	"context"
	"database/sql"
	"mana/internal/models"

	"github.com/google/uuid"
)

type UserStore interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CheckUserExistsByUsername(ctx context.Context, username string) (bool, error)
}

type sqlUserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return &sqlUserStore{db: db}
}

func (userStore *sqlUserStore) Create(ctx context.Context, user *models.User) error {
	insertUserSQL := `
		INSERT INTO users (id, username, email, password, activity_status, account_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := userStore.db.ExecContext(
		ctx,
		insertUserSQL,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.ActivityStatus,
		user.AccountStatus,
		user.CreatedAt,
	)

	return err
}

func (userStore *sqlUserStore) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	selectUserSQL := `
		SELECT id, username, email, password, activity_status, account_status, created_at
		FROM users
		WHERE id = $1
	`

	userRow := userStore.db.QueryRowContext(ctx, selectUserSQL, id)

	var user models.User
	err := userRow.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ActivityStatus,
		&user.AccountStatus,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (userStore *sqlUserStore) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	selectUserSQL := `
		SELECT id, username, email, password, activity_status, account_status, created_at
		FROM users
		WHERE email = $1
	`

	userRow := userStore.db.QueryRowContext(ctx, selectUserSQL, email)

	var user models.User
	err := userRow.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ActivityStatus,
		&user.AccountStatus,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (userStore *sqlUserStore) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	selectUserSQL := `
		SELECT id, username, email, password, activity_status, account_status, created_at
		FROM users
		WHERE username = $1
	`

	userRow := userStore.db.QueryRowContext(ctx, selectUserSQL, username)

	var user models.User
	err := userRow.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ActivityStatus,
		&user.AccountStatus,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (userStore *sqlUserStore) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	selectUserEmailSQL := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	row := userStore.db.QueryRowContext(ctx, selectUserEmailSQL, email)

	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

func (userStore *sqlUserStore) CheckUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	selectUserUsernameSQL := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	row := userStore.db.QueryRowContext(ctx, selectUserUsernameSQL, username)

	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

func (userStore *sqlUserStore) UpdateActivityStatus(ctx context.Context, id uuid.UUID, status string) error {
	updateActivityStatusSQL := `UPDATE users SET activity_status = $1 WHERE id = $2`

	_, err := userStore.db.ExecContext(ctx, updateActivityStatusSQL, status, id)
	return err
}

func (userStore *sqlUserStore) UpdateAccountStatus(ctx context.Context, id uuid.UUID, status string) error {
	updateAccountStatusSQL := `UPDATE users SET account_status = $1 WHERE id = $2`

	_, err := userStore.db.ExecContext(ctx, updateAccountStatusSQL, status, id)
	return err
}
