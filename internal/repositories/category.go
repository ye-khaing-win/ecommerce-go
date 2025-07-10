package repositories

import (
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repositories/sqlconnect"
)

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
