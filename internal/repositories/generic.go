package repositories

import "context"

type IRepository[T any] interface {
	Get(context.Context, map[string]interface{}) ([]T, error)
	GetOne(context.Context, string) (T, error)
	Create(context.Context, *T) error
	Update(context.Context, string, *T) error
	Delete(context.Context, string) error
}
