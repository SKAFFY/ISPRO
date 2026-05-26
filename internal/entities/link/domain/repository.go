package domain

import "context"

type LinkRepository interface {
	List(ctx context.Context) ([]*Link, error)
	Get(ctx context.Context, id int64) (*Link, error)
	Create(ctx context.Context, sourceID, targetID int64) (*Link, error)
	Update(ctx context.Context, id int64, sourceID, targetID int64) (*Link, error)
	Delete(ctx context.Context, id int64) error
}
