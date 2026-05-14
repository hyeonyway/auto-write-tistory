package template

import "testing"

func TestMarkdownListFormatsItems(t *testing.T) {
	got := markdownList([]string{"DP", "BFS"})
	want := "- DP\n- BFS"
	if got != want {
		t.Fatalf("markdownList() = %q, want %q", got, want)
	}
}

func TestMarkdownListReturnsDashForEmptyInput(t *testing.T) {
	got := markdownList([]string{})
	if got != "-" {
		t.Fatalf("markdownList() = %q, want %q", got, "-")
	}
}

func TestReplaceVarsReplacesTemplateTokens(t *testing.T) {
	got := replaceVars("[{{platform}}] {{problemTitle}}", map[string]string{
		"platform":     "Programmers",
		"problemTitle": "등굣길",
	})
	want := "[Programmers] 등굣길"
	if got != want {
		t.Fatalf("replaceVars() = %q, want %q", got, want)
	}
}
