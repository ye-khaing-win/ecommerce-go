package repositories

import (
	"context"
	"database/sql"
	"ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/models"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrCategoryNotFound = errors.New("category not found")

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) Repository[models.Category] {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List(ctx context.Context) ([]models.Category, error) {
	selected := middlewares.Selected(ctx)
	filtered := middlewares.Filtered(ctx)
	sorted := middlewares.Sorted(ctx)
	fmt.Println("Filtered: ", filtered)
	fmt.Println("Selected: ", selected)
	fmt.Println("Sorted: ", sorted)

	dbFields := GetDBFields(models.Category{}, selected)

	query := fmt.Sprintf("SELECT %s FROM categories WHERE 1 = 1", strings.Join(dbFields, ", "))

	// FILTERING
	var args []any
	for k, v := range filtered {
		query += fmt.Sprintf(" AND %s = ?", k)
		args = append(args, v)
	}

	if len(sorted) > 0 {

		var s []string
		for field, order := range sorted {
			s = append(s, fmt.Sprintf("%s %s", field, order))
		}
		fmt.Println("s: ", strings.Join(s, ","))
		query += " ORDER BY " + strings.Join(s, ",")
	}

	fmt.Println("Query: ", query)

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category

	for rows.Next() {
		var cat models.Category

		v := reflect.ValueOf(&cat).Elem()
		t := v.Type()

		var scanArgs []any
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i).Tag.Get("db")

			if _, ok := selected[field]; ok {
				scanArgs = append(scanArgs, v.Field(i).Addr().Interface())
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
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

	var sets []string
	var args []any
	var cols []string

	catVal := reflect.ValueOf(cat)
	catType := catVal.Type()
	for i := 0; i < catVal.NumField(); i++ {
		set := catType.Field(i).Tag.Get("db")
		arg := catVal.Field(i).Interface()

		cols = append(cols, set)
		if !reflect.ValueOf(arg).IsZero() {
			sets = append(sets, set+" = ?")
			args = append(args, arg)
		}

	}

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

	row := r.db.QueryRow(fmt.Sprintf(
		"SELECT %s FROM categories WHERE id = ?",
		strings.Join(cols, ", "),
	), id)
	if err = row.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Category{}, ErrCategoryNotFound
		}
		return models.Category{}, err
	}

	fmt.Println("Cols: ", cols)

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
