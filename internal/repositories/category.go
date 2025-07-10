package repositories

import (
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repositories/sqlconnect"
)

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
