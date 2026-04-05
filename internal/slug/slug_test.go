package slug

import (
	"strings"
	"testing"
)

func TestFromTitle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"RAG vs PKM", "rag-vs-pkm"},
		{"Fix ALB OIDC Issue", "fix-alb-oidc-issue"},
		{"Idea for Workshop", "idea-for-workshop"},
		{"", "untitled"},
		{"   ", "untitled"},
		{"!!!###", "untitled"},
		{"ä ö ü Ä Ö Ü ß", "ae-oe-ue-ae-oe-ue-ss"},
		{"hello--world", "hello-world"},
		{"-leading-and-trailing-", "leading-and-trailing"},
		{"under_score and.dot/slash", "under-score-and-dot-slash"},
		{"MixedCASE123", "mixedcase123"},
		{"already-a-slug", "already-a-slug"},
	}

	for _, tc := range tests {
		got := FromTitle(tc.input)
		if got != tc.want {
			t.Errorf("FromTitle(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFromTitleTruncation(t *testing.T) {
	// Build a title that would produce a slug > 64 chars
	long := strings.Repeat("word-", 20) // "word-word-word-..." = 100 chars
	got := FromTitle(long)
	if len(got) > maxSlugLen {
		t.Errorf("slug length %d exceeds maxSlugLen %d: %q", len(got), maxSlugLen, got)
	}
	if strings.HasSuffix(got, "-") {
		t.Errorf("slug ends with hyphen: %q", got)
	}
}
