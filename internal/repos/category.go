package repos

import (
	"context"
	"database/sql"
	"ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/models"
	"ecommerce-go/pkg/utils"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrCategoryNotFound = errors.New("category not found")

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CRUDRepo[models.Category] {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List(ctx context.Context) ([]models.Category, error) {
	filters := middlewares.GetFilters(ctx)
	sorts := middlewares.GetSorts(ctx)
	fmt.Println("Filters: ", filters)
	fmt.Println("Sorts: ", sorts)

	query := `
			SELECT id, name, description, created_at
			FROM categories 
			WHERE 1 = 1
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

	var cats []models.Category

	for rows.Next() {
		var cat models.Category

		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	if cats == nil {
		cats = []models.Category{}
	}

	return cats, nil
}
func (r *categoryRepository) Get(id int) (models.Category, error) {

	var cat models.Category
	query := "SELECT id, name, description FROM categories WHERE id = ?"

	if err := r.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Category{}, ErrCategoryNotFound
		}
		return models.Category{}, err
	}

	return cat, nil

}
func (r *categoryRepository) Create(cat models.Category) (models.Category, error) {

	res, err := r.db.Exec("INSERT INTO categories (name, description) VALUES (?, ?)", cat.Name, cat.Description)
	if err != nil {
		return models.Category{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Category{}, err
	}
	cat.ID = int(id)

	return cat, nil

}

func (r *categoryRepository) Update(id int, cat models.Category) (models.Category, error) {

	sets, args := utils.ApplyUpdates(cat)

	query := fmt.Sprintf("UPDATE categories SET %s WHERE id = ?", strings.Join(sets, ", "))
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return models.Category{}, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return models.Category{}, err
	}

	row := r.db.QueryRow(
		`SELECT id, name, description, created_at FROM categories WHERE id = ?`, id)
	if err = row.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Category{}, ErrCategoryNotFound
		}
		return models.Category{}, err
	}

	return cat, nil

}
func (r *categoryRepository) Delete(id int) error {

	res, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}

	return nil

}

func GetDBFields(model any, whitelist map[string]struct{}) []string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		dbTag := t.Field(i).Tag.Get("db")
		if _, ok := whitelist[dbTag]; ok {
			fields = append(fields, dbTag)
		}
	}

	return fields
}
