package tistory

import (
	"context"
	"errors"
	"strings"

	publishertistory "devlog-studio/backend/internal/publisher/tistory"
)

type Service struct {
	publisher *publishertistory.Publisher
}

func NewService(publisher *publishertistory.Publisher) *Service {
	return &Service{publisher: publisher}
}

func (s *Service) FetchCategories(ctx context.Context, req FetchCategoriesRequest) (*FetchCategoriesResponse, error) {
	if strings.TrimSpace(req.BlogURL) == "" {
		return nil, errors.New("blog url is required")
	}
	if strings.TrimSpace(req.Login.ID) == "" || req.Login.Password == "" {
		return nil, errors.New("login is required")
	}

	items, err := s.publisher.FetchCategories(ctx, req.BlogURL, req.Login.ID, req.Login.Password)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryOption, 0, len(items))
	for _, item := range items {
		result = append(result, CategoryOption{
			CategoryID: item.CategoryID,
			Label:      item.Label,
		})
	}

	return &FetchCategoriesResponse{Items: result}, nil
}
