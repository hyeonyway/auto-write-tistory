package tistory

import (
	"context"
	"errors"
	"strings"

	publishertistory "devlog-studio/backend/internal/publisher/tistory"
	"devlog-studio/backend/internal/settings"
)

type Service struct {
	publisher *publishertistory.Publisher
	settings  *settings.Service
}

func NewService(publisher *publishertistory.Publisher, settings *settings.Service) *Service {
	return &Service{publisher: publisher, settings: settings}
}

func (s *Service) FetchCategories(ctx context.Context, req FetchCategoriesRequest) (*FetchCategoriesResponse, error) {
	if strings.TrimSpace(req.BlogURL) == "" {
		return nil, errors.New("blog url is required")
	}

	items, err := s.publisher.FetchCategories(ctx, req.BlogURL)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryOption, 0, len(items))
	settingsCategories := make([]settings.TistoryCategory, 0, len(items))
	for _, item := range items {
		result = append(result, CategoryOption{
			CategoryID: item.CategoryID,
			Label:      item.Label,
		})
		settingsCategories = append(settingsCategories, settings.TistoryCategory{
			CategoryID: item.CategoryID,
			Label:      item.Label,
		})
	}
	if err := s.settings.SaveCategories(ctx, settingsCategories); err != nil {
		return nil, err
	}

	return &FetchCategoriesResponse{Items: result}, nil
}

func (s *Service) StartSession(ctx context.Context) (*SessionStartResponse, error) {
	status, err := s.publisher.StartSession(ctx)
	if err != nil {
		return nil, err
	}
	return &SessionStartResponse{Connected: status.Connected, Message: status.Message}, nil
}

func (s *Service) ConfirmSession(ctx context.Context) (*SessionStatusResponse, error) {
	cfg, err := s.settings.Get(ctx)
	if err != nil {
		return nil, err
	}
	status, err := s.publisher.ConfirmSession(ctx, cfg.Tistory.BlogURL)
	if err != nil {
		return nil, err
	}
	return &SessionStatusResponse{Connected: status.Connected, Message: status.Message}, nil
}

func (s *Service) SessionStatus(ctx context.Context) (*SessionStatusResponse, error) {
	cfg, err := s.settings.Get(ctx)
	if err != nil {
		return nil, err
	}
	status, err := s.publisher.SessionStatus(ctx, cfg.Tistory.BlogURL)
	if err != nil {
		return nil, err
	}
	return &SessionStatusResponse{Connected: status.Connected, Message: status.Message}, nil
}

func (s *Service) DeleteSession() error {
	return s.publisher.DeleteSession()
}
