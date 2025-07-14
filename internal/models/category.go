package models

type Category struct {
	ID          int    `json:"id,omitempty" db:"id"`
	Name        string `json:"name,omitempty" db:"name" validate:"required"`
	Description string `json:"description,omitempty" db:"description"  validate:"required"`
}
