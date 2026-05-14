package template

import (
	"context"
	"fmt"
	"strings"
)

type AlgorithmInput struct {
	Platform        string   `json:"platform"`
	ProblemTitle    string   `json:"problemTitle"`
	ProblemURL      string   `json:"problemUrl"`
	Language        string   `json:"language"`
	Categories      []string `json:"categories"`
	Approach        string   `json:"approach"`
	Code            string   `json:"code"`
	TimeComplexity  string   `json:"timeComplexity"`
	SpaceComplexity string   `json:"spaceComplexity"`
	Review          string   `json:"review"`
	Tags            []string `json:"tags"`
}

type RenderResult struct {
	Title           string `json:"title"`
	ContentMarkdown string `json:"contentMarkdown"`
}

type Renderer struct {
	repo *Repository
}

func NewRenderer(repo *Repository) *Renderer {
	return &Renderer{repo: repo}
}

func (r *Renderer) RenderAlgorithm(ctx context.Context, input AlgorithmInput) (*RenderResult, error) {
	tpl, err := r.repo.FindDefaultByType(ctx, "ALGORITHM")
	if err != nil {
		return nil, fmt.Errorf("load template: %w", err)
	}

	values := map[string]string{
		"platform":        input.Platform,
		"problemTitle":    input.ProblemTitle,
		"problemUrl":      input.ProblemURL,
		"language":        input.Language,
		"categories":      markdownList(input.Categories),
		"approach":        input.Approach,
		"code":            input.Code,
		"timeComplexity":  optional(input.TimeComplexity),
		"spaceComplexity": optional(input.SpaceComplexity),
		"review":          optional(input.Review),
		"tags":            strings.Join(input.Tags, ", "),
	}

	return &RenderResult{
		Title:           replaceVars(tpl.TitleFormat, values),
		ContentMarkdown: replaceVars(tpl.BodyFormat, values),
	}, nil
}

func replaceVars(format string, values map[string]string) string {
	out := format
	for key, value := range values {
		out = strings.ReplaceAll(out, "{{"+key+"}}", value)
	}
	return out
}

func markdownList(items []string) string {
	if len(items) == 0 {
		return "-"
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			lines = append(lines, "- "+trimmed)
		}
	}
	if len(lines) == 0 {
		return "-"
	}
	return strings.Join(lines, "\n")
}

func optional(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "-"
	}
	return trimmed
}
