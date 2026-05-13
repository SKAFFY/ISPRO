package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkheyfets/ispro-app/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) ListEntries(ctx context.Context) ([]*models.Entry, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, content, created_at, updated_at FROM entries ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		var e models.Entry
		if err := rows.Scan(&e.ID, &e.Title, &e.Content, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entries = append(entries, &e)
	}
	return entries, nil
}

func (r *Repository) GetEntry(ctx context.Context, id int64) (*models.Entry, error) {
	var e models.Entry
	err := r.pool.QueryRow(ctx, "SELECT id, title, content, created_at, updated_at FROM entries WHERE id = $1", id).
		Scan(&e.ID, &e.Title, &e.Content, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository) CreateEntry(ctx context.Context, title, content string) (*models.Entry, error) {
	var e models.Entry
	err := r.pool.QueryRow(ctx,
		"INSERT INTO entries (title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at",
		title, content,
	).Scan(&e.ID, &e.Title, &e.Content, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}
	return &e, nil
}

func (r *Repository) UpdateEntry(ctx context.Context, id int64, title, content string) (*models.Entry, error) {
	var e models.Entry
	err := r.pool.QueryRow(ctx,
		"UPDATE entries SET title = $1, content = $2, updated_at = NOW() WHERE id = $3 RETURNING id, title, content, created_at, updated_at",
		title, content, id,
	).Scan(&e.ID, &e.Title, &e.Content, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %w", err)
	}
	return &e, nil
}

func (r *Repository) DeleteEntry(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM entries WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("entry not found")
	}
	return nil
}

func (r *Repository) ListLinks(ctx context.Context) ([]*models.Link, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, source_id, target_id FROM links ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to list links: %w", err)
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var l models.Link
		if err := rows.Scan(&l.ID, &l.SourceID, &l.TargetID); err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, &l)
	}
	return links, nil
}

func (r *Repository) GetLink(ctx context.Context, id int64) (*models.Link, error) {
	var l models.Link
	err := r.pool.QueryRow(ctx, "SELECT id, source_id, target_id FROM links WHERE id = $1", id).
		Scan(&l.ID, &l.SourceID, &l.TargetID)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repository) CreateLink(ctx context.Context, sourceID, targetID int64) (*models.Link, error) {
	var l models.Link
	err := r.pool.QueryRow(ctx,
		"INSERT INTO links (source_id, target_id) VALUES ($1, $2) RETURNING id, source_id, target_id",
		sourceID, targetID,
	).Scan(&l.ID, &l.SourceID, &l.TargetID)
	if err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}
	return &l, nil
}

func (r *Repository) UpdateLink(ctx context.Context, id int64, sourceID, targetID int64) (*models.Link, error) {
	var l models.Link
	err := r.pool.QueryRow(ctx,
		"UPDATE links SET source_id = $1, target_id = $2 WHERE id = $3 RETURNING id, source_id, target_id",
		sourceID, targetID, id,
	).Scan(&l.ID, &l.SourceID, &l.TargetID)
	if err != nil {
		return nil, fmt.Errorf("failed to update link: %w", err)
	}
	return &l, nil
}

func (r *Repository) DeleteLink(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM links WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("link not found")
	}
	return nil
}