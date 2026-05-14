package tistory

import (
	"context"
	"errors"
)

type Publisher struct {
	session *SessionManager
}

func NewPublisher(runtimeDir string) *Publisher {
	return &Publisher{session: NewSessionManager(runtimeDir)}
}

func (p *Publisher) StartSession(ctx context.Context) (*SessionStatus, error) {
	return p.session.Start(ctx)
}

func (p *Publisher) ConfirmSession(ctx context.Context, blogURL string) (*SessionStatus, error) {
	client, err := p.client()
	if err != nil {
		return &SessionStatus{Connected: false, Message: "저장된 Tistory 세션이 없습니다."}, nil
	}
	return p.session.Confirm(ctx, client, blogURL)
}

func (p *Publisher) SessionStatus(ctx context.Context, blogURL string) (*SessionStatus, error) {
	client, err := p.client()
	if err != nil {
		if errors.Is(err, ErrSessionRequired) {
			return &SessionStatus{Connected: false, Message: "저장된 Tistory 세션이 없습니다."}, nil
		}
		return nil, err
	}
	return p.session.Status(ctx, client, blogURL)
}

func (p *Publisher) DeleteSession() error {
	return p.session.Delete()
}

func (p *Publisher) FetchCategories(ctx context.Context, blogURL string) ([]CategoryOption, error) {
	client, err := p.client()
	if err != nil {
		return nil, err
	}
	return client.FetchCategories(ctx, blogURL)
}

func (p *Publisher) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	client, err := p.client()
	if err != nil {
		return nil, err
	}
	return client.Publish(ctx, req)
}

func (p *Publisher) client() (*Client, error) {
	statePath, err := p.session.RequireState()
	if err != nil {
		return nil, err
	}
	return NewClient(statePath), nil
}
