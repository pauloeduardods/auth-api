package admin_repo

import (
	"auth-api/src/internal/domain/admin"
	"auth-api/src/pkg/logger"
	"database/sql"
)

type AdminRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewAdminRepository(db *sql.DB, logger logger.Logger) admin.AdminRepository {
	return &AdminRepository{
		db:     db,
		logger: logger,
	}
}

func (r *AdminRepository) GetByID(input *admin.GetAdminInput) (*admin.Admin, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	var adm admin.Admin
	query := `SELECT id, name, email FROM admins WHERE id = $1`
	if err := r.db.QueryRow(query, input.ID).Scan(&adm.ID, &adm.Name, &adm.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, admin.ErrAdminNotFound
		}
		r.logger.Error("Error getting admin by ID: %v", err)
		return nil, err
	}
	return &adm, nil
}

func (r *AdminRepository) GetByEmail(input *admin.GetAdminByEmailInput) (*admin.Admin, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	var adm admin.Admin
	query := `SELECT id, name, email FROM admins WHERE email = $1`
	if err := r.db.QueryRow(query, input.Email).Scan(&adm.ID, &adm.Name, &adm.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, admin.ErrAdminNotFound
		}
		r.logger.Error("Error getting admin by email: %v", err)
		return nil, err
	}
	return &adm, nil
}

func (r *AdminRepository) Create(input *admin.CreateAdminInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO admins (id, name, email) VALUES ($1, $2, $3)`
	if _, err := r.db.Exec(query, input.ID.String(), input.Name, input.Email); err != nil {
		r.logger.Error("Error creating admin: %v", err)
		return err
	}
	return nil
}

func (r *AdminRepository) Update(input *admin.UpdateAdminInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := `UPDATE admins SET name = COALESCE($1, name), email = COALESCE($2, email) WHERE id = $4`
	_, err := r.db.Exec(query, input.Name, input.Email, input.ID.String())
	if err != nil {
		r.logger.Error("Error updating admin: %v", err)
	}
	return err
}

func (r *AdminRepository) Delete(input *admin.DeleteAdminInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	query := `DELETE FROM admins WHERE id = $1`
	if _, err := r.db.Exec(query, input.ID.String()); err != nil {
		r.logger.Error("Error deleting admin: %v", err)
		return err
	}
	return nil
}
