package models

type Item struct {
	ID          int     `json:"id,omitempty" db:"id"`
	Name        string  `json:"name,omitempty" db:"name" validate:"required"`
	Description string  `json:"description,omitempty" db:"description"`
	CategoryID  int     `json:"category_id,omitempty" db:"category_id" validate:"required"`
	Price       float64 `json:"price,omitempty" db:"price" validate:"required"`
	Available   bool    `json:"available,omitempty" db:"available"`
	Published   bool    `json:"published,omitempty" db:"published"`
	Category
}
