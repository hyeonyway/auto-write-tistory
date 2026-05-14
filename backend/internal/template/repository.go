package template

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindDefaultByType(ctx context.Context, postType string) (*Template, error) {
	const query = `
        SELECT id, type, name, title_format, body_format, is_default
        FROM templates
        WHERE type = $1 AND is_default = TRUE
        ORDER BY id ASC
        LIMIT 1
    `
	var tpl Template
	err := r.pool.QueryRow(ctx, query, postType).Scan(
		&tpl.ID,
		&tpl.Type,
		&tpl.Name,
		&tpl.TitleFormat,
		&tpl.BodyFormat,
		&tpl.IsDefault,
	)
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}
