package repos

import (
	"context"
	"database/sql"
	mw "ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/models"
	"ecommerce-go/pkg/utils"
	"errors"
)

var ErrAdminNotFound = errors.New("admin not found")

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

func (r *AdminRepository) List(ctx context.Context) ([]models.Admin, error) {
	filters := mw.GetFilters(ctx)
	sorts := mw.GetSorts(ctx)

	query := `
			SELECT id, first_name, last_name, email, active, role, created_at
			FROM admins
			`
	// FILTERING
	filterQuery, args := utils.ApplyFilters(filters)
	query += filterQuery

	// SORTING
	sortQuery := utils.ApplySorts(sorts)
	query += sortQuery

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []models.Admin

	for rows.Next() {
		var admin models.Admin

		if err := rows.Scan(&admin.ID, &admin.FirstName,
			&admin.LastName, &admin.Email,
			&admin.Active, &admin.Role,
			&admin.CreatedAt); err != nil {
			return nil, err
		}

		admins = append(admins, admin)
	}

	if admins == nil {
		admins = []models.Admin{}
	}

	return admins, nil

}
func (r *AdminRepository) Create(dto models.Admin) (models.Admin, error) {
	// HASH PASSWORD
	passHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return models.Admin{}, err
	}
	dto.Password = passHash

	query := `
			INSERT INTO admins
			(first_name, last_name, email, password, role) VALUES
			(?, ?, ?, ?, ?)
			`

	res, err := r.db.Exec(
		query, dto.FirstName, dto.LastName,
		dto.Email, dto.Password, dto.Role,
	)
	if err != nil {
		return models.Admin{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Admin{}, err
	}

	var admin models.Admin
	query = `
			SELECT id, first_name, last_name, email, role, created_at 
			FROM admins
			WHERE id = ?
			`

	row := r.db.QueryRow(query, id)
	if err := row.Scan(&admin.ID, &admin.FirstName,
		&admin.LastName, &admin.Email,
		&admin.Role, &admin.CreatedAt); err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}
