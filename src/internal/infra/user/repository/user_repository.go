package user_repo

import (
	"auth-api/src/internal/domain/user"
	"auth-api/src/pkg/logger"
	"database/sql"
)

type UserRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewUserRepository(db *sql.DB, logger logger.Logger) user.UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) GetByID(input *user.GetUserInput) (*user.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	var usr user.User
	query := `SELECT id, name, email, phone FROM users WHERE id = $1`
	if err := r.db.QueryRow(query, input.ID).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		r.logger.Error("Error getting user by ID: %v", err)
		return nil, err
	}
	return &usr, nil
}

func (r *UserRepository) GetByEmail(input *user.GetUserByEmailInput) (*user.User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	var usr user.User
	query := `SELECT id, name, email, phone FROM users WHERE email = $1`
	if err := r.db.QueryRow(query, input.Email).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		r.logger.Error("Error getting user by email: %v", err)
		return nil, err
	}
	return &usr, nil
}

func (r *UserRepository) Create(input *user.CreateUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO users (id, name, email, phone) VALUES ($1, $2, $3, $4)`
	if _, err := r.db.Exec(query, input.ID.String(), input.Name, input.Email, input.Phone); err != nil {
		r.logger.Error("Error creating user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) Update(input *user.UpdateUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := `UPDATE users SET name = COALESCE($1, name), email = COALESCE($2, email), phone = COALESCE($3, phone) WHERE id = $4`
	_, err := r.db.Exec(query, input.Name, input.Email, input.Phone, input.ID.String())
	if err != nil {
		r.logger.Error("Error updating user: %v", err)
	}
	return err
}

func (r *UserRepository) Delete(input *user.DeleteUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := `DELETE FROM users WHERE id = $1`
	if _, err := r.db.Exec(query, input.ID.String()); err != nil {
		r.logger.Error("Error deleting user: %v", err)
		return err
	}
	return nil
}
