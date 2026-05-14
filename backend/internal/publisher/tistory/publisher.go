package tistory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Publisher struct {
	remoteURL string
	client    *http.Client
	mu        sync.Mutex
	session   *browserSession
}

type browserSession struct {
	ID        string
	BlogURL   string
	LoginID   string
	ExpiresAt time.Time
}

type PublishRequest struct {
	BlogURL         string
	LoginID         string
	LoginPassword   string
	Title           string
	ContentMarkdown string
	Visibility      int
	CategoryID      string
	Tags            []string
}

type CategoryOption struct {
	CategoryID string `json:"categoryId"`
	Label      string `json:"label"`
}

type PublishResult struct {
	ExternalPostID string
	ExternalURL    string
}

func NewPublisher(remoteURL string) *Publisher {
	return &Publisher{
		remoteURL: strings.TrimRight(remoteURL, "/"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	if strings.TrimSpace(req.BlogURL) == "" {
		return nil, errors.New("tistory blog url is empty")
	}
	if strings.TrimSpace(req.LoginID) == "" || req.LoginPassword == "" {
		return nil, errors.New("tistory login is required")
	}

	sessionID, err := p.acquireEditorSession(ctx, req.BlogURL, req.LoginID, req.LoginPassword)
	if err != nil {
		return nil, err
	}

	if err := p.fillEditor(ctx, sessionID, req); err != nil {
		return nil, fmt.Errorf("fill tistory editor: %w", err)
	}

	currentURL, _ := p.currentURL(ctx, sessionID)
	return &PublishResult{
		ExternalPostID: parsePostID(currentURL),
		ExternalURL:    currentURL,
	}, nil
}

func (p *Publisher) FetchCategories(ctx context.Context, blogURL, loginID, loginPassword string) ([]CategoryOption, error) {
	if strings.TrimSpace(blogURL) == "" {
		return nil, errors.New("tistory blog url is empty")
	}
	if strings.TrimSpace(loginID) == "" || loginPassword == "" {
		return nil, errors.New("tistory login is required")
	}

	sessionID, err := p.acquireEditorSession(ctx, blogURL, loginID, loginPassword)
	if err != nil {
		return nil, err
	}

	if err := p.clickFirst(ctx, sessionID, []string{"#category-btn"}); err != nil {
		return nil, fmt.Errorf("open category dropdown: %w", err)
	}
	if err := p.waitForElement(ctx, sessionID, "#category-list [role='option']", 15*time.Second); err != nil {
		return nil, fmt.Errorf("load category list: %w", err)
	}

	items, err := p.readCategoryOptions(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("read category list: %w", err)
	}
	return items, nil
}

func (p *Publisher) acquireEditorSession(ctx context.Context, blogURL string, loginID string, loginPassword string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.session != nil && p.session.BlogURL == normalizedBlogURL(blogURL) && p.session.LoginID == strings.TrimSpace(loginID) && time.Now().Before(p.session.ExpiresAt) {
		if err := p.openEditorWithExistingSession(ctx, p.session.ID, blogURL); err == nil {
			p.session.ExpiresAt = time.Now().Add(30 * time.Minute)
			return p.session.ID, nil
		}
		p.deleteSession(context.Background(), p.session.ID)
		p.session = nil
	}

	sessionID, err := p.createSession(ctx)
	if err != nil {
		return "", fmt.Errorf("create selenium session: %w", err)
	}

	if err := p.openEditorWithLogin(ctx, sessionID, blogURL, loginID, loginPassword); err != nil {
		p.deleteSession(context.Background(), sessionID)
		return "", err
	}

	p.session = &browserSession{
		ID:        sessionID,
		BlogURL:   normalizedBlogURL(blogURL),
		LoginID:   strings.TrimSpace(loginID),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	return sessionID, nil
}

func (p *Publisher) openEditorWithExistingSession(ctx context.Context, sessionID string, blogURL string) error {
	writeURL := normalizedBlogURL(blogURL) + "/manage/newpost"
	if err := p.navigate(ctx, sessionID, writeURL); err != nil {
		return fmt.Errorf("open tistory editor: %w", err)
	}
	return p.waitForEditorReady(ctx, sessionID, 20*time.Second)
}

func (p *Publisher) openEditorWithLogin(ctx context.Context, sessionID string, blogURL string, loginID string, loginPassword string) error {
	if err := p.navigate(ctx, sessionID, "https://www.tistory.com/auth/login"); err != nil {
		return fmt.Errorf("open tistory login: %w", err)
	}

	if err := p.tryKakaoLogin(ctx, sessionID, loginID, loginPassword); err != nil {
		return fmt.Errorf("tistory login failed: %w", err)
	}

	writeURL := normalizedBlogURL(blogURL) + "/manage/newpost"
	if err := p.navigate(ctx, sessionID, writeURL); err != nil {
		return fmt.Errorf("open tistory editor: %w", err)
	}
	if err := p.waitForEditorReady(ctx, sessionID, 2*time.Minute); err != nil {
		return fmt.Errorf("wait for tistory editor after kakao authentication: %w", err)
	}
	return nil
}

func (p *Publisher) createSession(ctx context.Context) (string, error) {
	payload := map[string]any{
		"capabilities": map[string]any{
			"alwaysMatch": map[string]any{
				"browserName": "chrome",
				"goog:chromeOptions": map[string]any{
					"args": []string{
						"--disable-gpu",
						"--no-sandbox",
						"--disable-dev-shm-usage",
						"--window-size=1440,1200",
					},
				},
			},
		},
	}
	var res struct {
		Value struct {
			SessionID string `json:"sessionId"`
		} `json:"value"`
		SessionID string `json:"sessionId"`
	}
	if err := p.webdriver(ctx, http.MethodPost, "/session", payload, &res); err != nil {
		return "", err
	}
	if res.Value.SessionID != "" {
		return res.Value.SessionID, nil
	}
	return res.SessionID, nil
}

func (p *Publisher) deleteSession(ctx context.Context, sessionID string) {
	_ = p.webdriver(ctx, http.MethodDelete, "/session/"+sessionID, nil, nil)
}

func (p *Publisher) navigate(ctx context.Context, sessionID string, target string) error {
	return p.webdriver(ctx, http.MethodPost, "/session/"+sessionID+"/url", map[string]string{"url": target}, nil)
}

func (p *Publisher) tryKakaoLogin(ctx context.Context, sessionID string, loginID string, password string) error {
	// Tistory/Kakao markup changes over time. Keep selectors centralized here.
	_ = p.clickFirst(ctx, sessionID, []string{
		"a.btn_login.link_kakao_id",
		"a[href*='kakao']",
		"button[class*='kakao']",
	})

	if err := p.typeFirst(ctx, sessionID, []string{
		"input[name='loginId']",
		"input[name='email']",
		"input[type='email']",
		"input#loginId",
	}, loginID); err != nil {
		return err
	}
	if err := p.typeFirst(ctx, sessionID, []string{
		"input[name='password']",
		"input[type='password']",
		"input#password",
	}, password); err != nil {
		return err
	}
	return p.clickFirst(ctx, sessionID, []string{
		"button[type='submit']",
		"button.submit",
		"input[type='submit']",
	})
}

func (p *Publisher) fillEditor(ctx context.Context, sessionID string, req PublishRequest) error {
	time.Sleep(2 * time.Second)
	_ = p.clickFirst(ctx, sessionID, []string{
		"button[data-mode='markdown']",
		"button[title*='Markdown']",
		"button[title*='마크다운']",
	})

	if strings.TrimSpace(req.CategoryID) != "" {
		if err := p.selectCategory(ctx, sessionID, req.CategoryID); err != nil {
			return fmt.Errorf("select category: %w", err)
		}
	}

	if err := p.typeFirst(ctx, sessionID, []string{
		"textarea[placeholder*='제목']",
		"input[placeholder*='제목']",
		"textarea.title",
		"input.title",
	}, req.Title); err != nil {
		return err
	}

	if err := p.typeFirst(ctx, sessionID, []string{
		"textarea",
		".CodeMirror textarea",
		"[contenteditable='true']",
	}, req.ContentMarkdown); err != nil {
		return err
	}

	if len(req.Tags) > 0 {
		_ = p.typeFirst(ctx, sessionID, []string{
			"input[placeholder*='태그']",
			"input[name*='tag']",
		}, strings.Join(req.Tags, ","))
	}

	_ = p.clickFirst(ctx, sessionID, []string{
		"button[data-visibility='private']",
		"button[title*='비공개']",
		"label:has(input[value='0'])",
	})

	return p.clickFirst(ctx, sessionID, []string{
		"button[type='submit']",
		"button[class*='publish']",
		"button[class*='save']",
		"button:has-text('완료')",
	})
}

func (p *Publisher) selectCategory(ctx context.Context, sessionID string, categoryID string) error {
	if strings.TrimSpace(categoryID) == "" {
		return nil
	}
	if err := p.clickFirst(ctx, sessionID, []string{"#category-btn"}); err != nil {
		return err
	}
	if err := p.waitForElement(ctx, sessionID, "#category-list [role='option']", 15*time.Second); err != nil {
		return err
	}
	selector := fmt.Sprintf(`#category-list [category-id="%s"]`, categoryID)
	return p.clickFirst(ctx, sessionID, []string{
		selector,
		`#category-list [role='option'][category-id="` + categoryID + `"]`,
	})
}

func (p *Publisher) typeFirst(ctx context.Context, sessionID string, selectors []string, text string) error {
	var lastErr error
	for _, selector := range selectors {
		elementID, err := p.findElement(ctx, sessionID, selector)
		if err != nil {
			lastErr = err
			continue
		}
		if err := p.webdriver(ctx, http.MethodPost, "/session/"+sessionID+"/element/"+elementID+"/clear", map[string]any{}, nil); err != nil {
			lastErr = err
		}
		if err := p.webdriver(ctx, http.MethodPost, "/session/"+sessionID+"/element/"+elementID+"/value", map[string]any{"text": text}, nil); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	if lastErr != nil {
		return lastErr
	}
	return errors.New("element not found")
}

func (p *Publisher) clickFirst(ctx context.Context, sessionID string, selectors []string) error {
	var lastErr error
	for _, selector := range selectors {
		elementID, err := p.findElement(ctx, sessionID, selector)
		if err != nil {
			lastErr = err
			continue
		}
		if err := p.webdriver(ctx, http.MethodPost, "/session/"+sessionID+"/element/"+elementID+"/click", map[string]any{}, nil); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	if lastErr != nil {
		return lastErr
	}
	return errors.New("element not found")
}

func (p *Publisher) findElement(ctx context.Context, sessionID string, selector string) (string, error) {
	var res struct {
		Value map[string]string `json:"value"`
	}
	err := p.webdriver(ctx, http.MethodPost, "/session/"+sessionID+"/element", map[string]string{
		"using": "css selector",
		"value": selector,
	}, &res)
	if err != nil {
		return "", err
	}
	for _, value := range res.Value {
		if value != "" {
			return value, nil
		}
	}
	return "", errors.New("element id missing")
}

func (p *Publisher) currentURL(ctx context.Context, sessionID string) (string, error) {
	var res struct {
		Value string `json:"value"`
	}
	if err := p.webdriver(ctx, http.MethodGet, "/session/"+sessionID+"/url", nil, &res); err != nil {
		return "", err
	}
	return res.Value, nil
}

func (p *Publisher) readCategoryOptions(ctx context.Context, sessionID string) ([]CategoryOption, error) {
	var res struct {
		Value []CategoryOption `json:"value"`
	}
	script := `return Array.from(document.querySelectorAll('#category-list [role="option"]')).map((el) => ({
		categoryId: el.getAttribute('category-id') || '',
		label: (el.getAttribute('aria-label') || el.textContent || '').trim(),
	}));`
	if err := p.executeScript(ctx, sessionID, script, nil, &res); err != nil {
		return nil, err
	}
	return res.Value, nil
}

func (p *Publisher) waitForElement(ctx context.Context, sessionID string, selector string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if _, err := p.findElement(ctx, sessionID, selector); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(300 * time.Millisecond)
	}
	if lastErr != nil {
		return lastErr
	}
	return errors.New("element not found")
}

func (p *Publisher) waitForEditorReady(ctx context.Context, sessionID string, timeout time.Duration) error {
	if err := p.waitForElement(ctx, sessionID, "#category-btn", timeout); err != nil {
		currentURL, _ := p.currentURL(ctx, sessionID)
		if strings.Contains(currentURL, "accounts.kakao.com") || strings.Contains(currentURL, "auth/login") {
			return fmt.Errorf("kakao authentication is still required; complete verification in the Selenium browser and retry")
		}
		return err
	}
	return nil
}

func (p *Publisher) executeScript(ctx context.Context, sessionID string, script string, args []any, dst any) error {
	if args == nil {
		args = []any{}
	}
	payload := map[string]any{
		"script": script,
		"args":   args,
	}
	return p.webdriver(ctx, http.MethodPost, "/session/"+sessionID+"/execute/sync", payload, dst)
}

func (p *Publisher) webdriver(ctx context.Context, method string, path string, payload any, dst any) error {
	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, p.remoteURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var webdriverErr struct {
			Value struct {
				Error   string `json:"error"`
				Message string `json:"message"`
			} `json:"value"`
		}
		_ = json.NewDecoder(res.Body).Decode(&webdriverErr)
		message := webdriverErr.Value.Message
		if message == "" {
			message = res.Status
		}
		return errors.New(message)
	}
	if dst == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(dst)
}

func parsePostID(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func normalizedBlogURL(blogURL string) string {
	return strings.TrimRight(strings.TrimSpace(blogURL), "/")
}
