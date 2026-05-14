package post

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"devlog-studio/backend/internal/publisher/tistory"
	"devlog-studio/backend/internal/settings"
	apptemplate "devlog-studio/backend/internal/template"
)

type Service struct {
	repo     *Repository
	renderer *apptemplate.Renderer
	settings *settings.Service
	tistory  *tistory.Publisher
}

func NewService(repo *Repository, renderer *apptemplate.Renderer, settings *settings.Service, tistory *tistory.Publisher) *Service {
	return &Service{
		repo:     repo,
		renderer: renderer,
		settings: settings,
		tistory:  tistory,
	}
}

func (s *Service) Preview(ctx context.Context, req PreviewRequest) (*apptemplate.RenderResult, error) {
	if err := validateAlgorithm(req.Type, req.Input); err != nil {
		return nil, err
	}
	return s.renderer.RenderAlgorithm(ctx, req.Input)
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	if err := validateAlgorithm(req.Type, req.Input); err != nil {
		return nil, err
	}
	rendered, err := s.renderer.RenderAlgorithm(ctx, req.Input)
	if err != nil {
		return nil, err
	}
	created, err := s.repo.Create(ctx, req, rendered.Title, rendered.ContentMarkdown)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{
		ID:     created.ID,
		Title:  created.Title,
		Status: created.Status,
	}, nil
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (*UpdateResponse, error) {
	if err := validateAlgorithm(req.Type, req.Input); err != nil {
		return nil, err
	}
	rendered, err := s.renderer.RenderAlgorithm(ctx, req.Input)
	if err != nil {
		return nil, err
	}
	updated, err := s.repo.Update(ctx, id, req, rendered.Title, rendered.ContentMarkdown)
	if err != nil {
		return nil, err
	}
	return &UpdateResponse{
		ID:     updated.ID,
		Title:  updated.Title,
		Status: updated.Status,
	}, nil
}

func (s *Service) List(ctx context.Context, page int, size int) (*ListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	items, total, err := s.repo.List(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return &ListResponse{Items: items, Page: page, Size: size, Total: total}, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Post, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) PublishTistory(ctx context.Context, id int64, req PublishRequest) (*PublishResponse, error) {
	if strings.TrimSpace(req.Login.ID) == "" || req.Login.Password == "" {
		return nil, NewValidationError("TISTORY_LOGIN_REQUIRED", "Tistory 로그인 정보를 입력해주세요.")
	}

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cfg, err := s.settings.Get(ctx)
	if err != nil {
		return nil, err
	}

	categoryID := firstNonEmpty(req.CategoryID, cfg.Tistory.CategoryID)
	visibility := req.Visibility
	if visibility != 0 && visibility != 3 {
		visibility = cfg.Tistory.DefaultVisibility
	}

	result, err := s.tistory.Publish(ctx, tistory.PublishRequest{
		BlogURL:         cfg.Tistory.BlogURL,
		LoginID:         req.Login.ID,
		LoginPassword:   req.Login.Password,
		Title:           post.Title,
		ContentMarkdown: post.ContentMarkdown,
		Visibility:      visibility,
		CategoryID:      categoryID,
		Tags:            req.Tags,
	})
	if err != nil {
		_ = s.repo.MarkFailed(ctx, id, sanitizeError(err))
		return nil, classifyTistoryError(err)
	}

	if err := s.repo.MarkPublished(ctx, id, "TISTORY", result.ExternalPostID, result.ExternalURL); err != nil {
		return nil, err
	}
	externalURL := result.ExternalURL
	return &PublishResponse{Success: true, ExternalURL: &externalURL}, nil
}

func validateAlgorithm(postType string, input apptemplate.AlgorithmInput) error {
	if postType != TypeAlgorithm {
		return NewValidationError("VALIDATION_ERROR", "type must be ALGORITHM")
	}
	checks := map[string]string{
		"platform":     input.Platform,
		"problemTitle": input.ProblemTitle,
		"problemUrl":   input.ProblemURL,
		"language":     input.Language,
		"approach":     input.Approach,
		"code":         input.Code,
	}
	for field, value := range checks {
		if strings.TrimSpace(value) == "" {
			return NewValidationError("VALIDATION_ERROR", field+" is required")
		}
	}
	if _, err := url.ParseRequestURI(input.ProblemURL); err != nil {
		return NewValidationError("VALIDATION_ERROR", "problemUrl must be a valid URL")
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func sanitizeError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if len(message) > 500 {
		return message[:500]
	}
	return message
}

func classifyTistoryError(err error) error {
	message := err.Error()
	switch {
	case strings.Contains(message, "kakao authentication is still required"):
		return NewValidationError("TISTORY_AUTH_REQUIRED", "카카오 추가 인증이 필요합니다. Selenium 브라우저에서 인증을 완료한 뒤 다시 시도해주세요.")
	case strings.Contains(message, "tistory login failed"):
		return NewValidationError("TISTORY_LOGIN_FAILED", "Tistory 로그인에 실패했습니다. 카카오 로그인 정보를 확인해주세요.")
	case strings.Contains(message, "wait for tistory editor"):
		return NewValidationError("TISTORY_SESSION_NOT_READY", "Tistory 글쓰기 화면을 열지 못했습니다. 인증 상태를 확인해주세요.")
	default:
		return NewValidationError("TISTORY_PUBLISH_ERROR", "Tistory 발행 중 오류가 발생했습니다.")
	}
}

type ValidationError struct {
	Code    string
	Message string
}

func NewValidationError(code string, message string) *ValidationError {
	return &ValidationError{Code: code, Message: message}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func IsValidationError(err error) (*ValidationError, bool) {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return validationErr, true
	}
	return nil, false
}
