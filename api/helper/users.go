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

// GetUserByUsernameOrEmail retrieves a user by username or email
func (s *UserStore) GetUserByUsernameOrEmail(ctx context.Context, tx *sql.Tx, user *models.User) (error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE username = $1 OR email = $2
	`

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		slog.Warn("User not found", "username", user.Username, "email", user.Email)
		return nil
	}

	if err != nil {
		slog.Error("Failed to get user by username or email", "error", err)
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

// GetUserByProviderID retrieves a user by OAuth provider and provider ID
func (s *UserStore) GetUserByProviderID(ctx context.Context, provider, providerID string) (*models.User, error) {
	query := `
		SELECT id, username, email, provider, provider_id, avatar_url, email_verified, created_at, updated_at
		FROM users
		WHERE provider = $1 AND provider_id = $2
	`

	user := &models.User{}
	err := s.Db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Provider,
		&user.ProviderID,
		&user.AvatarURL,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		slog.Error("Failed to get user by provider ID", "error", err, "provider", provider)
		return nil, err
	}

	return user, nil
}

// CreateOAuthUser creates a new OAuth user
func (s *UserStore) CreateOAuthUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users(username, email, provider, provider_id, avatar_url, email_verified, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, NULL)
		RETURNING id, created_at, updated_at
	`

	err := s.Db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Provider,
		user.ProviderID,
		user.AvatarURL,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		slog.Error("Failed to create OAuth user", "error", err, "email", user.Email)
		return err
	}

	slog.Info("OAuth user created successfully", "id", user.ID, "email", user.Email, "provider", user.Provider)
	return nil
}

// FindOrCreateOAuthUser finds an existing OAuth user or creates a new one
func (s *UserStore) FindOrCreateOAuthUser(ctx context.Context, googleUser *models.GoogleUserInfo) (*models.User, error) {
	// First, try to find existing user by provider ID
	user, err := s.GetUserByProviderID(ctx, "google", googleUser.ID)
	if err != nil {
		return nil, err
	}

	// If user exists, return it
	if user != nil {
		slog.Info("Found existing OAuth user", "id", user.ID, "email", user.Email)
		return user, nil
	}

	// User doesn't exist, create new one
	// Generate username from email (part before @)
	username := googleUser.Email
	if idx := len(googleUser.Email); idx > 0 {
		for i, c := range googleUser.Email {
			if c == '@' {
				username = googleUser.Email[:i]
				break
			}
		}
	}

	newUser := &models.User{
		Username:      username,
		Email:         googleUser.Email,
		Provider:      "google",
		ProviderID:    &googleUser.ID,
		AvatarURL:     &googleUser.Picture,
		EmailVerified: googleUser.VerifiedEmail,
	}

	if err := s.CreateOAuthUser(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
