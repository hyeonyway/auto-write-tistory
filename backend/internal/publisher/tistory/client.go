package tistory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	statePath string
	client    *http.Client
}

type storageState struct {
	Cookies []storageCookie `json:"cookies"`
}

type storageCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type CategoryOption struct {
	CategoryID string `json:"categoryId"`
	Label      string `json:"label"`
}

type PublishRequest struct {
	BlogURL         string
	Title           string
	ContentMarkdown string
	Visibility      int
	CategoryID      string
	Tags            []string
}

type PublishResult struct {
	ExternalPostID string
	ExternalURL    string
}

func NewClient(statePath string) *Client {
	return &Client{
		statePath: statePath,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Validate(ctx context.Context, blogURL string) error {
	if strings.TrimSpace(blogURL) != "" {
		body, err := c.get(ctx, normalizedBlogURL(blogURL)+"/manage/newpost")
		if err == nil && bytes.Contains(body, []byte("window.Config")) {
			return nil
		}
	}
	_, err := c.get(ctx, "https://www.tistory.com/legacy/member/blog/api/myBlogs")
	return err
}

func (c *Client) FetchCategories(ctx context.Context, blogURL string) ([]CategoryOption, error) {
	body, err := c.get(ctx, normalizedBlogURL(blogURL)+"/manage/newpost")
	if err != nil {
		return nil, err
	}
	return ParseCategories(body)
}

func (c *Client) Publish(ctx context.Context, req PublishRequest) (*PublishResult, error) {
	if strings.TrimSpace(req.BlogURL) == "" {
		return nil, errors.New("tistory blog url is empty")
	}
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.ContentMarkdown) == "" {
		return nil, errors.New("title and content are required")
	}

	visibility := req.Visibility
	if visibility == 3 {
		visibility = 20
	}
	if visibility != 0 && visibility != 15 && visibility != 20 {
		visibility = 0
	}

	payload := map[string]any{
		"id":                    "0",
		"title":                 req.Title,
		"content":               req.ContentMarkdown,
		"visibility":            visibility,
		"category":              req.CategoryID,
		"tag":                   strings.Join(req.Tags, ","),
		"published":             1,
		"type":                  "post",
		"uselessMarginForEntry": 1,
		"cclCommercial":         0,
		"cclDerive":             0,
		"attachments":           []any{},
		"recaptchaValue":        "",
		"draftSequence":         nil,
	}

	body, err := c.postJSON(ctx, normalizedBlogURL(req.BlogURL)+"/manage/post.json", payload)
	if err != nil {
		return nil, err
	}

	var responsePayload struct {
		URL   string `json:"url"`
		ID    any    `json:"id"`
		Post  any    `json:"post"`
		Entry any    `json:"entry"`
	}
	_ = json.Unmarshal(body, &responsePayload)

	externalURL := responsePayload.URL
	if externalURL == "" {
		externalURL = normalizedBlogURL(req.BlogURL)
	}
	externalPostID := ""
	if responsePayload.ID != nil {
		externalPostID = fmt.Sprint(responsePayload.ID)
	}
	return &PublishResult{
		ExternalPostID: externalPostID,
		ExternalURL:    externalURL,
	}, nil
}

func (c *Client) get(ctx context.Context, target string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	if err := c.attachCookies(req); err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) postJSON(ctx context.Context, target string, payload any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if err := c.attachCookies(req); err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("tistory returned %s", res.Status)
	}
	return body, nil
}

func (c *Client) attachCookies(req *http.Request) error {
	state, err := c.readState()
	if err != nil {
		return err
	}
	host := req.URL.Hostname()
	for _, cookie := range state.Cookies {
		domain := strings.TrimPrefix(cookie.Domain, ".")
		if domain == "" || host == domain || strings.HasSuffix(host, "."+domain) {
			req.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value, Path: cookie.Path})
		}
	}
	return nil
}

func (c *Client) readState() (*storageState, error) {
	data, err := readFile(c.statePath)
	if err != nil {
		return nil, err
	}
	var state storageState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	if len(state.Cookies) == 0 {
		return nil, ErrSessionRequired
	}
	return &state, nil
}

func ParseCategories(body []byte) ([]CategoryOption, error) {
	re := regexp.MustCompile(`(?s)categories\s*:\s*(\[[^\n;]*\])`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return parseCategoriesByPattern(body)
	}

	var raw []map[string]any
	if err := json.Unmarshal(matches[1], &raw); err != nil {
		return parseCategoriesByPattern(body)
	}

	items := make([]CategoryOption, 0, len(raw)+1)
	items = append(items, CategoryOption{CategoryID: "0", Label: "카테고리 없음"})
	for _, item := range raw {
		id := firstString(item, "id", "categoryId", "category_id")
		label := firstString(item, "name", "label", "title")
		if id == "" || label == "" {
			continue
		}
		items = append(items, CategoryOption{CategoryID: id, Label: label})
	}
	return items, nil
}

func parseCategoriesByPattern(body []byte) ([]CategoryOption, error) {
	re := regexp.MustCompile(`(?s)(?:id|categoryId)\s*:\s*['"]?(\d+)['"]?.{0,300}?(?:label|name|title)\s*:\s*['"]([^'"]+)['"]`)
	matches := re.FindAllSubmatch(body, -1)
	if len(matches) == 0 {
		return nil, errors.New("tistory categories config not found")
	}
	items := []CategoryOption{{CategoryID: "0", Label: "카테고리 없음"}}
	seen := map[string]bool{"0": true}
	for _, match := range matches {
		id := string(match[1])
		label := string(match[2])
		if id == "" || label == "" || seen[id] {
			continue
		}
		seen[id] = true
		items = append(items, CategoryOption{CategoryID: id, Label: label})
	}
	return items, nil
}

func firstString(item map[string]any, keys ...string) string {
	for _, key := range keys {
		switch value := item[key].(type) {
		case string:
			if value != "" {
				return value
			}
		case float64:
			return fmt.Sprintf("%.0f", value)
		}
	}
	return ""
}

func normalizedBlogURL(blogURL string) string {
	return strings.TrimRight(strings.TrimSpace(blogURL), "/")
}
