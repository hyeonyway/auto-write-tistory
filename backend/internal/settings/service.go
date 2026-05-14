package settings

import (
	"context"
	"encoding/json"
	"strconv"
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

func decodeTistory(values map[string]string) TistorySettings {
	visibility, err := strconv.Atoi(values["tistory.default_visibility"])
	if err != nil {
		visibility = 0
	}

	tags := []string{}
	_ = json.Unmarshal([]byte(values["tistory.default_tags"]), &tags)

	blogURL := values["tistory.blog_url"]
	if blogURL == "" {
		blogURL = "https://hyeonyway.tistory.com"
	}

	return TistorySettings{
		BlogURL:           blogURL,
		CategoryID:        values["tistory.category_id"],
		DefaultVisibility: visibility,
		DefaultTags:       tags,
	}
}
