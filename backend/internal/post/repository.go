package post

import (
	"context"
	"encoding/json"
	"errors"

	apptemplate "devlog-studio/backend/internal/template"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("post not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, input CreateRequest, title string, markdown string) (*Post, error) {
	inputJSON, err := json.Marshal(input.Input)
	if err != nil {
		return nil, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const insertPost = `
        INSERT INTO posts (title, type, source_platform, source_url, language, content_markdown, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, title, type, source_platform, source_url, language, content_markdown, content_html,
            status, external_provider, external_post_id, external_url, error_message, published_at,
            created_at, updated_at
    `
	row := tx.QueryRow(ctx, insertPost,
		title,
		TypeAlgorithm,
		input.Input.Platform,
		input.Input.ProblemURL,
		input.Input.Language,
		markdown,
		StatusDraft,
	)
	created, err := scanPost(row)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `INSERT INTO post_inputs (post_id, input_json) VALUES ($1, $2)`, created.ID, inputJSON); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *Repository) List(ctx context.Context, page int, size int) ([]Post, int64, error) {
	offset := (page - 1) * size

	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM posts`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, `
        SELECT id, title, type, source_platform, source_url, language, content_markdown, content_html,
            status, external_provider, external_post_id, external_url, error_message, published_at,
            created_at, updated_at
        FROM posts
        ORDER BY created_at DESC, id DESC
        LIMIT $1 OFFSET $2
    `, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, *p)
	}
	return posts, total, rows.Err()
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*Post, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, title, type, source_platform, source_url, language, content_markdown, content_html,
            status, external_provider, external_post_id, external_url, error_message, published_at,
            created_at, updated_at
        FROM posts
        WHERE id = $1
    `, id)
	p, err := scanPost(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	input, err := r.findLatestInput(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Input = input
	return p, err
}

func (r *Repository) Update(ctx context.Context, id int64, input UpdateRequest, title string, markdown string) (*Post, error) {
	inputJSON, err := json.Marshal(input.Input)
	if err != nil {
		return nil, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
        UPDATE posts
        SET title = $2,
            type = $3,
            source_platform = $4,
            source_url = $5,
            language = $6,
            content_markdown = $7,
            status = $8,
            external_provider = NULL,
            external_post_id = NULL,
            external_url = NULL,
            error_message = NULL,
            published_at = NULL,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING id, title, type, source_platform, source_url, language, content_markdown, content_html,
            status, external_provider, external_post_id, external_url, error_message, published_at,
            created_at, updated_at
    `,
		id,
		title,
		TypeAlgorithm,
		input.Input.Platform,
		input.Input.ProblemURL,
		input.Input.Language,
		markdown,
		StatusDraft,
	)
	updated, err := scanPost(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `INSERT INTO post_inputs (post_id, input_json) VALUES ($1, $2)`, id, inputJSON); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	updated.Input = &input.Input
	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) MarkPublished(ctx context.Context, id int64, provider string, externalPostID string, externalURL string) error {
	tag, err := r.pool.Exec(ctx, `
        UPDATE posts
        SET status = $2,
            external_provider = $3,
            external_post_id = NULLIF($4, ''),
            external_url = NULLIF($5, ''),
            error_message = NULL,
            published_at = CURRENT_TIMESTAMP,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `, id, StatusPublished, provider, externalPostID, externalURL)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id int64, message string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE posts
        SET status = $2,
            error_message = $3,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `, id, StatusFailed, message)
	return err
}

func (r *Repository) findLatestInput(ctx context.Context, postID int64) (*apptemplate.AlgorithmInput, error) {
	var input apptemplate.AlgorithmInput
	row := r.pool.QueryRow(ctx, `
        SELECT input_json
        FROM post_inputs
        WHERE post_id = $1
        ORDER BY created_at DESC, id DESC
        LIMIT 1
    `, postID)
	var raw []byte
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPost(row rowScanner) (*Post, error) {
	var p Post
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Type,
		&p.SourcePlatform,
		&p.SourceURL,
		&p.Language,
		&p.ContentMarkdown,
		&p.ContentHTML,
		&p.Status,
		&p.ExternalProvider,
		&p.ExternalPostID,
		&p.ExternalURL,
		&p.ErrorMessage,
		&p.PublishedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
