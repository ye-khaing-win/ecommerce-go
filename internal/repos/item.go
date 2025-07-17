package repos

import (
	"context"
	"database/sql"
	mw "ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/models"
	"ecommerce-go/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ErrItemNotFound = errors.New("item not found")

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) List(ctx context.Context) ([]models.Item, error) {
	filters := mw.GetFilters(ctx)
	sorts := mw.GetSorts(ctx)

	query := `
			SELECT i.id, i.name, i.description, i.price, i.category_id,
				(
					SELECT JSON_OBJECT("id", c.id, "name", c.name, "description", c.description)
					FROM categories c
					WHERE c.id = i.category_id
				) AS category
			FROM items i
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

	var items []models.Item

	for rows.Next() {
		var item models.Item
		var catJSON []byte

		if err := rows.Scan(&item.ID, &item.Name,
			&item.Description, &item.Price,
			&item.CategoryID, &catJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(catJSON, &item.Category); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if items == nil {
		items = []models.Item{}
	}

	return items, nil

}
func (r *ItemRepository) Get(id int) (models.Item, error) {
	var item models.Item
	var catJSON []byte
	query := `
			SELECT 
				i.id, i.name, i.description, i.price, i.category_id,
				(
					SELECT JSON_OBJECT("id", c.id, "name", c.name, "description", c.description)
					FROM categories c
					WHERE c.id = i.category_id
				) AS category
			FROM items i
			WHERE i.id = ?
			`
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&item.ID,
		&item.Name, &item.Description,
		&item.Price, &item.CategoryID,
		&catJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, ErrItemNotFound
		}
		return models.Item{}, err
	}

	if err := json.Unmarshal(catJSON, &item.Category); err != nil {
		return models.Item{}, err
	}

	return item, nil
}
func (r *ItemRepository) Create(dto models.Item) (models.Item, error) {
	res, err := r.db.Exec(`INSERT INTO items 
    (name, description, category_id, price) 
	VALUES (?, ?, ?, ?)`, dto.Name, dto.Description, dto.CategoryID, dto.Price)

	if err != nil {
		return models.Item{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Item{}, err
	}

	var item models.Item
	query := `SELECT id, name, description, price FROM items WHERE id = ?`
	if err := r.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
		return models.Item{}, err
	}

	return item, nil

}
func (r *ItemRepository) Update(id int, item models.Item) (models.Item, error) {
	sets, args := utils.ApplyUpdates(item)
	query := fmt.Sprintf("UPDATE items SET %s WHERE id = ?", strings.Join(sets, ", "))
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return models.Item{}, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return models.Item{}, err
	}

	row := r.db.QueryRow(`SELECT id, name, description, price FROM items WHERE id = ?`, id)
	if err = row.Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, ErrItemNotFound
		}
		return models.Item{}, err
	}

	return item, nil
}
func (r *ItemRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrItemNotFound
	}

	return nil
}
func (r *ItemRepository) ListByCategory(ctx context.Context, catID int) ([]models.Item, error) {
	filters := mw.GetFilters(ctx)
	sorts := mw.GetSorts(ctx)

	query := `
			SELECT id, name, description, price, category_id
			FROM items
			WHERE category_id = ?
			`
	args := []any{catID}
	filterQuery, filterArgs := utils.ApplyFilters(filters)
	query += filterQuery
	args = append(args, filterArgs...)

	sortQuery := utils.ApplySorts(sorts)
	query += sortQuery

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item

		if err := rows.Scan(&item.ID, &item.Name,
			&item.Description, &item.Price,
			&item.CategoryID); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if items == nil {
		items = []models.Item{}
	}

	return items, nil
}
