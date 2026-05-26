package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/entities/entry/query"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresRepository) List(ctx context.Context) ([]*domain.Entry, error) {
	sql, args, err := query.ListEntries().ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.Entry
	for rows.Next() {
		var e domain.Entry
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&e.ID, &e.Title, &e.Content, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		e.CreatedAt = createdAt
		e.UpdatedAt = updatedAt
		entries = append(entries, &e)
	}
	return entries, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (*domain.Entry, error) {
	sql, args, err := query.GetEntry(id).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var e domain.Entry
	var createdAt, updatedAt time.Time
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&e.ID, &e.Title, &e.Content, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	e.CreatedAt = createdAt
	e.UpdatedAt = updatedAt
	return &e, nil
}

func (r *PostgresRepository) Create(ctx context.Context, title, content string) (*domain.Entry, error) {
	sql, args, err := query.CreateEntry(title, content).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var e domain.Entry
	var createdAt, updatedAt time.Time
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&e.ID, &e.Title, &e.Content, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}
	e.CreatedAt = createdAt
	e.UpdatedAt = updatedAt
	return &e, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id int64, title, content string) (*domain.Entry, error) {
	sql, args, err := query.UpdateEntry(id, title, content).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var e domain.Entry
	var createdAt, updatedAt time.Time
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&e.ID, &e.Title, &e.Content, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %w", err)
	}
	e.CreatedAt = createdAt
	e.UpdatedAt = updatedAt
	return &e, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	sql, args, err := query.DeleteEntry(id).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("entry not found")
	}
	return nil
}
