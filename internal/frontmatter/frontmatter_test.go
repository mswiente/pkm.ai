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

// TestMarshalSimpleRoundTripsSpecialTitles guards against a real bug: titles
// containing YAML-significant characters (a colon-space, in particular) were
// written unquoted, which turned "key: value" into a nested YAML mapping and
// broke Parse for every Readwise article with a colon in its title.
func TestMarshalSimpleRoundTripsSpecialTitles(t *testing.T) {
	titles := []string{
		"GitHub - foo/bar: a tool for doing things",
		`"The decline of Europe is not inevitable, despite how much...`,
		"Ends with a colon:",
		"Has a # hash and a : colon and a , comma",
		"Plain simple title",
		"",
	}
	for _, title := range titles {
		n := Note{
			Title:   title,
			Type:    "resource",
			Status:  "inbox",
			Source:  "readwise",
			Created: "2026-04-03",
			Tags:    []string{"readwise"},
		}
		out := MarshalSimple(n)
		got, _, err := Parse(out)
		if err != nil {
			t.Fatalf("round-trip Parse failed for title %q: %v\nmarshaled:\n%s", title, err, out)
		}
		if got.Title != title {
			t.Errorf("round-trip Title = %q, want %q", got.Title, title)
		}
	}
}

func TestFormatTagsWithSpecialChars(t *testing.T) {
	got := FormatTags([]string{"normal", "has: colon", "has, comma"})
	_, _, err := Parse([]byte("---\ntitle: t\ntype: resource\nstatus: inbox\ncreated: 2026-01-01\ntags: " + got + "\n---\n"))
	if err != nil {
		t.Fatalf("tags with special chars broke parsing: %v\ntags: %s", err, got)
	}
}
