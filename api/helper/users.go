package helper

import (
	"context"
	"database/sql"
	"log/slog"

	"shawty-ur/api/models"
)

// UserStore handles all database operations for users
type UserStore struct {
	Db *sql.DB
}

// NewUserStore creates a new user store
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{Db: db}
}

// Create inserts a new user into the database
// func (s *UserStore) CreateUser(ctx context.Context, user *models.User) error {
// 	query := `
// 		INSERT INTO users(email, password, username)
// 		VALUES ($1, $2, $3)
// 		ON CONFLICT(email) DO NOTHING
// 		RETURNING id, created_at, updated_at
// 	`

// 	err := s.db.QueryRowContext(
// 		ctx,
// 		query,
// 		user.Email,
// 		user.Password.Hash,
// 		user.Username,
// 	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

// 	if err == sql.ErrNoRows {
// 		slog.Warn("User already exists", "email", user.Email)
// 		return nil
// 	}

// 	if err != nil {
// 		slog.Error("Failed to create user", "error", err)
// 		return err
// 	}

// 	slog.Info("User created successfully", "id", user.ID, "email", user.Email)
// 	return nil
// }

func (s *UserStore) CreateUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	slog.Info("Creating User Query: ", "email: ", user.Email)
	query := `INSERT INTO users(email, password, username)
			VALUES ($1, $2, $3) ON CONFLICT(email) DO NOTHING
			RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx, query,
		user.Email,
		user.Password.Hash,
		user.Username).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		//user already exists
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// GetByEmail retrieves a user by email
func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := s.Db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		slog.Warn("User not found", "email", email)
		return nil, nil
	}

	if err != nil {
		slog.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	return user, nil
}

// GetUserByUsernameOrEmail retrieves a user by username or email
func (s *UserStore) GetUserByUsernameOrEmail(ctx context.Context, email string, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $2
	`

	user := &models.User{}
	err := s.Db.QueryRowContext(ctx, query, username, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		slog.Warn("User not found", "username", username, "email", email)
		return nil, nil
	}

	if err != nil {
		slog.Error("Failed to get user by username or email", "error", err)
		return nil, err
	}

	return user, nil
}

// List retrieves all users
func (s *UserStore) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan user row", "error", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// Delete removes a user by ID
func (s *UserStore) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.Db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete user", "error", err, "id", id)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("User not found for deletion", "id", id)
	}

	slog.Info("User deleted", "id", id)
	return nil
}

// Update updates user information
func (s *UserStore) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	err := s.Db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		slog.Error("Failed to update user", "error", err, "id", user.ID)
		return err
	}

	slog.Info("User updated", "id", user.ID)
	return nil
}
