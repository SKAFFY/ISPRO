package domain

import "context"

type EntryRepository interface {
	List(ctx context.Context) ([]*Entry, error)
	Get(ctx context.Context, id int64) (*Entry, error)
	Create(ctx context.Context, title, content string) (*Entry, error)
	Update(ctx context.Context, id int64, title, content string) (*Entry, error)
	Delete(ctx context.Context, id int64) error
}
