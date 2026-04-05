package frontmatter

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	content := `---
title: Test Note
type: inbox
status: inbox
source: manual
created: 2026-04-03
tags: [aws, oidc]
---

## Content

Some body text here.
`
	n, body, err := Parse([]byte(content))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.Title != "Test Note" {
		t.Errorf("Title = %q, want %q", n.Title, "Test Note")
	}
	if n.Type != "inbox" {
		t.Errorf("Type = %q, want %q", n.Type, "inbox")
	}
	if len(n.Tags) != 2 || n.Tags[0] != "aws" || n.Tags[1] != "oidc" {
		t.Errorf("Tags = %v, want [aws oidc]", n.Tags)
	}
	if !strings.Contains(body, "Some body text") {
		t.Errorf("body missing expected text, got: %q", body)
	}
}

func TestParseNoFrontmatter(t *testing.T) {
	_, _, err := Parse([]byte("# Just a heading\n\nNo frontmatter here."))
	if err == nil {
		t.Error("expected error for missing frontmatter, got nil")
	}
}

func TestFormatTags(t *testing.T) {
	tests := []struct {
		tags []string
		want string
	}{
		{nil, "[]"},
		{[]string{}, "[]"},
		{[]string{"aws"}, "[aws]"},
		{[]string{"aws", "oidc", "alb"}, "[aws, oidc, alb]"},
	}
	for _, tc := range tests {
		got := FormatTags(tc.tags)
		if got != tc.want {
			t.Errorf("FormatTags(%v) = %q, want %q", tc.tags, got, tc.want)
		}
	}
}

func TestMarshalSimple(t *testing.T) {
	n := Note{
		Title:   "My Note",
		Type:    "inbox",
		Status:  "inbox",
		Source:  "manual",
		Created: "2026-04-03",
		Tags:    []string{"test"},
	}
	out := string(MarshalSimple(n))
	if !strings.HasPrefix(out, "---\n") {
		t.Error("output should start with ---")
	}
	if !strings.HasSuffix(out, "---\n") {
		t.Error("output should end with ---")
	}
	if !strings.Contains(out, "title: My Note") {
		t.Error("output missing title")
	}
	if !strings.Contains(out, "tags: [test]") {
		t.Error("output missing tags")
	}
	// Field order check: title must appear before type
	titleIdx := strings.Index(out, "title:")
	typeIdx := strings.Index(out, "type:")
	if titleIdx > typeIdx {
		t.Error("title should appear before type in frontmatter")
	}
}
