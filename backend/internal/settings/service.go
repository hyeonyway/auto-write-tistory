package settings

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context) (*Response, error) {
	values, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &Response{Tistory: decodeTistory(values)}, nil
}

func (s *Service) Update(ctx context.Context, input TistorySettings) (*Response, error) {
	tags, err := json.Marshal(input.DefaultTags)
	if err != nil {
		return nil, err
	}

	pairs := map[string]string{
		"tistory.blog_url":           input.BlogURL,
		"tistory.category_id":        input.CategoryID,
		"tistory.default_visibility": strconv.Itoa(input.DefaultVisibility),
		"tistory.default_tags":       string(tags),
	}
	for key, value := range pairs {
		if err := s.repo.Upsert(ctx, key, value); err != nil {
			return nil, err
		}
	}
	return s.Get(ctx)
}

func (s *Service) SaveCategories(ctx context.Context, categories []TistoryCategory) error {
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}
	if err := s.repo.Upsert(ctx, "tistory.categories", string(data)); err != nil {
		return err
	}
	return s.repo.Upsert(ctx, "tistory.categories_synced_at", time.Now().Format(time.RFC3339))
}

func decodeTistory(values map[string]string) TistorySettings {
	visibility, err := strconv.Atoi(values["tistory.default_visibility"])
	if err != nil {
		visibility = 0
	}

	tags := []string{}
	_ = json.Unmarshal([]byte(values["tistory.default_tags"]), &tags)

	categories := []TistoryCategory{}
	_ = json.Unmarshal([]byte(values["tistory.categories"]), &categories)

	blogURL := values["tistory.blog_url"]
	if blogURL == "" {
		blogURL = "https://hyeonyway.tistory.com"
	}

	categoryID := values["tistory.category_id"]
	categoryLabel := categoryLabel(categories, categoryID)

	return TistorySettings{
		BlogURL:            blogURL,
		CategoryID:         categoryID,
		CategoryLabel:      categoryLabel,
		DefaultVisibility:  visibility,
		DefaultTags:        tags,
		Categories:         categories,
		CategoriesSyncedAt: values["tistory.categories_synced_at"],
		Session: TistorySessionInfo{
			Connected: false,
			Message:   "세션 상태를 확인하지 않았습니다.",
		},
	}
}

func categoryLabel(categories []TistoryCategory, categoryID string) string {
	for _, category := range categories {
		if category.CategoryID == categoryID {
			return category.Label
		}
	}
	return ""
}
