package post

import (
	"time"

	apptemplate "devlog-studio/backend/internal/template"
)

const (
	TypeAlgorithm = "ALGORITHM"

	StatusDraft     = "DRAFT"
	StatusPublished = "PUBLISHED"
	StatusFailed    = "FAILED"
)

type Post struct {
	ID               int64                       `json:"id"`
	Title            string                      `json:"title"`
	Type             string                      `json:"type"`
	SourcePlatform   *string                     `json:"sourcePlatform"`
	SourceURL        *string                     `json:"sourceUrl"`
	Language         *string                     `json:"language"`
	ContentMarkdown  string                      `json:"contentMarkdown"`
	ContentHTML      *string                     `json:"contentHtml,omitempty"`
	Status           string                      `json:"status"`
	ExternalProvider *string                     `json:"externalProvider"`
	ExternalPostID   *string                     `json:"externalPostId,omitempty"`
	ExternalURL      *string                     `json:"externalUrl"`
	ErrorMessage     *string                     `json:"errorMessage"`
	PublishedAt      *time.Time                  `json:"publishedAt"`
	CreatedAt        time.Time                   `json:"createdAt"`
	UpdatedAt        time.Time                   `json:"updatedAt"`
	Input            *apptemplate.AlgorithmInput `json:"input,omitempty"`
}

type PreviewRequest struct {
	Type  string                     `json:"type"`
	Input apptemplate.AlgorithmInput `json:"input"`
}

type CreateRequest struct {
	Type  string                     `json:"type"`
	Input apptemplate.AlgorithmInput `json:"input"`
}

type CreateResponse struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type UpdateRequest struct {
	Type  string                     `json:"type"`
	Input apptemplate.AlgorithmInput `json:"input"`
}

type UpdateResponse struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type ListResponse struct {
	Items []Post `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int64  `json:"total"`
}

type PublishRequest struct {
	Visibility int      `json:"visibility"`
	CategoryID string   `json:"categoryId"`
	Tags       []string `json:"tags"`
}

type PublishResponse struct {
	Success     bool    `json:"success"`
	ExternalURL *string `json:"externalUrl"`
}
