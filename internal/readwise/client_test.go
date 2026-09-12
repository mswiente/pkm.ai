package readwise

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestListAllFiltersByLocation guards against a real bug: ListAll used to fetch
// documents updated after a timestamp across the user's entire Reader library
// (any location), not just the inbox. A bulk edit that touches old archived
// items bumps their updated_at, which flooded incremental syncs with thousands
// of unrelated historical items. ListAll must pass the location filter through
// to the API so only documents in the requested location come back.
func TestListAllFiltersByLocation(t *testing.T) {
	var gotLocation, gotUpdatedAfter string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotLocation = r.URL.Query().Get("location")
		gotUpdatedAfter = r.URL.Query().Get("updatedAfter")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listResponse{
			Count:   1,
			Results: []Document{{ID: "doc1", Title: "Test"}},
		})
	}))
	defer server.Close()

	orig := baseURLV3
	baseURLV3 = server.URL
	defer func() { baseURLV3 = orig }()

	c := NewClient("test-token")
	since := time.Date(2026, 8, 1, 9, 45, 27, 0, time.UTC)
	docs, err := c.ListAll(context.Background(), &since, "new")
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if len(docs) != 1 || docs[0].ID != "doc1" {
		t.Fatalf("ListAll returned %v, want one doc with ID doc1", docs)
	}
	if gotLocation != "new" {
		t.Errorf("location query param = %q, want %q", gotLocation, "new")
	}
	if gotUpdatedAfter == "" {
		t.Error("updatedAfter query param was not set")
	}
}

// TestListAllNoLocationFilter confirms an empty location omits the query param,
// preserving the ability to fetch across all locations when that is wanted.
func TestListAllNoLocationFilter(t *testing.T) {
	var sawLocationParam bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, sawLocationParam = r.URL.Query()["location"]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listResponse{})
	}))
	defer server.Close()

	orig := baseURLV3
	baseURLV3 = server.URL
	defer func() { baseURLV3 = orig }()

	c := NewClient("test-token")
	if _, err := c.ListAll(context.Background(), nil, ""); err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if sawLocationParam {
		t.Error("expected no location query param when location is empty")
	}
}

// TestListAllPaginates confirms pagination via nextPageCursor continues to work
// alongside the new location parameter.
func TestListAllPaginates(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("pageCursor") == "" {
			cursor := "page2"
			json.NewEncoder(w).Encode(listResponse{
				Results:        []Document{{ID: "doc1"}},
				NextPageCursor: &cursor,
			})
			return
		}
		json.NewEncoder(w).Encode(listResponse{
			Results: []Document{{ID: "doc2"}},
		})
	}))
	defer server.Close()

	orig := baseURLV3
	baseURLV3 = server.URL
	defer func() { baseURLV3 = orig }()

	c := NewClient("test-token")
	docs, err := c.ListAll(context.Background(), nil, "new")
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 requests across pagination, got %d", calls)
	}
	if len(docs) != 2 || docs[0].ID != "doc1" || docs[1].ID != "doc2" {
		t.Errorf("ListAll returned %v, want [doc1 doc2]", docs)
	}
}
