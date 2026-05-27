package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	"github.com/mkheyfets/ispro-app/internal/entities/link/query"
)

type PostgresLinkRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

func NewPostgresLinkRepository(pool *pgxpool.Pool) *PostgresLinkRepository {
	return &PostgresLinkRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresLinkRepository) List(ctx context.Context) ([]*domain.Link, error) {
	sql, args, err := query.ListLinks().ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list links: %w", err)
	}
	defer rows.Close()

	var links []*domain.Link
	for rows.Next() {
		var l domain.Link
		var createdAt time.Time
		if err := rows.Scan(&l.ID, &l.SourceID, &l.TargetID, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		l.CreatedAt = createdAt
		links = append(links, &l)
	}
	return links, nil
}

func (r *PostgresLinkRepository) Get(ctx context.Context, id int64) (*domain.Link, error) {
	sql, args, err := query.GetLink(id).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var l domain.Link
	var createdAt time.Time
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&l.ID, &l.SourceID, &l.TargetID, &createdAt)
	if err != nil {
		return nil, err
	}
	l.CreatedAt = createdAt
	return &l, nil
}

func (r *PostgresLinkRepository) Create(ctx context.Context, sourceID, targetID int64) (*domain.Link, error) {
	sql, args, err := query.CreateLink(sourceID, targetID).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var l domain.Link
	var createdAt time.Time
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&l.ID, &l.SourceID, &l.TargetID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}
	l.CreatedAt = createdAt
	return &l, nil
}

func (r *PostgresLinkRepository) Update(ctx context.Context, id int64, sourceID, targetID int64) (*domain.Link, error) {
	sql, args, err := query.UpdateLink(id, sourceID, targetID).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var l domain.Link
	var createdAt time.Time
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&l.ID, &l.SourceID, &l.TargetID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update link: %w", err)
	}
	l.CreatedAt = createdAt
	return &l, nil
}

func (r *PostgresLinkRepository) Delete(ctx context.Context, id int64) error {
	sql, args, err := query.DeleteLink(id).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("link not found")
	}
	return nil
}
