package db

import (
	"context"
	"database/sql"

	"mana/internal/models"

	"github.com/google/uuid"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

func (userStore *UserStore) InsertUser(ctx context.Context, user *models.User) error {
	insertUserSQL := `
		INSERT INTO users (id, username, email, password, activity_status, account_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := userStore.DB.ExecContext(
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

func (userStore *UserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	selectUserSQL := `
		SELECT id, username, email, password, activity_status, account_status, created_at
		FROM users
		WHERE email = $1
	`

	userRow := userStore.DB.QueryRowContext(ctx, selectUserSQL, email)

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

func (userStore *UserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	selectUserSQL := `
		SELECT id, username, email, password, activity_status, account_status, created_at
		FROM users
		WHERE username = $1
	`

	userRow := userStore.DB.QueryRowContext(ctx, selectUserSQL, username)

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

func (userStore *UserStore) GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	selectUserSQL := `
		SELECT id, username, email, password, activity_status, account_status, created_at
		FROM users
		WHERE ID = $1
	`

	userRow := userStore.DB.QueryRowContext(ctx, selectUserSQL, ID)

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

func (userStore *UserStore) CheckUserExistsByEmail(email string) (bool, error) {
	selectUserEmailSQL := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	row := userStore.DB.QueryRow(selectUserEmailSQL, email)

	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

func (userStore *UserStore) CheckUserExistsByUsername(username string) (bool, error) {
	selectUserUsernameSQL := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	row := userStore.DB.QueryRow(selectUserUsernameSQL, username)

	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

func (userStore *UserStore) UpdateActivityStatus(ctx context.Context, id uuid.UUID, status string) error {
	updateActivityStatusSQL := `UPDATE users SET activity_status = $1 WHERE id = $2`

	_, err := userStore.DB.ExecContext(ctx, updateActivityStatusSQL, status, id)
	return err
}

func (userStore *UserStore) UpdateAccountStatus(ctx context.Context, id uuid.UUID, status string) error {
	updateAccountStatusSQL := `UPDATE users SET account_status = $1 WHERE id = $2`

	_, err := userStore.DB.ExecContext(ctx, updateAccountStatusSQL, status, id)
	return err
}
