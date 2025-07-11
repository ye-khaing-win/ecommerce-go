package repositories

import (
	"database/sql"
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repositories/sqlconnect"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository struct {
	DB *sql.DB
}

func ListCategories() ([]models.Category, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, name, description FROM categories"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	return cats, nil
}
func GetCategory(id int) (models.Category, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return models.Category{}, err
	}
	defer db.Close()

	var cat models.Category
	query := "SELECT id, name, description FROM categories WHERE id = ?"

	if err := db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Category{}, ErrCategoryNotFound
		}
		return models.Category{}, err
	}

	return cat, nil

}
func CreateCategory(cat models.Category) (models.Category, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return models.Category{}, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO categories (name, description) VALUES (?, ?)", cat.Name, cat.Description)
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
func UpdateCategory(id int, cat models.Category) (models.Category, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return models.Category{}, err
	}
	defer db.Close()

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

	res, err := db.Exec(query, args...)
	if err != nil {
		return models.Category{}, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return models.Category{}, err
	}

	row := db.QueryRow(fmt.Sprintf("SELECT %s FROM categories WHERE id = ?", strings.Join(cols, ", ")), id)
	if err = row.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Category{}, ErrCategoryNotFound
		}
		return models.Category{}, err
	}

	fmt.Println("Cols: ", cols)

	return cat, nil

}
func DeleteCategory(id int) error {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM categories WHERE id = ?", id)
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
