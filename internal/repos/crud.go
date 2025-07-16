package repos

import "context"

type CRUDRepo[T any] interface {
	List(ctx context.Context) ([]T, error)
	Get(id int) (T, error)
	Create(cat T) (T, error)
	Update(id int, cat T) (T, error)
	Delete(id int) error
}
