package readwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURLV3 = "https://readwise.io/api/v3"
	baseURLV2 = "https://readwise.io/api/v2"
)

// Client is an HTTP client for the Readwise Reader API.
type Client struct {
	token      string
	httpClient *http.Client
}

// Document represents a Readwise Reader document or highlight.
// Highlights have ParentID set to the parent document's ID.
type Document struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Author      string         `json:"author"`
	URL         string         `json:"url"`
	SourceURL   string         `json:"source_url"`
	Summary     string         `json:"summary"`
	Notes       string         `json:"notes"`
	Content     string         `json:"content"`
	Category    string         `json:"category"`
	Location    string         `json:"location"`
	Tags        map[string]any `json:"tags"`
	ParentID    *string        `json:"parent_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	SavedAt     time.Time      `json:"saved_at"`
	WordCount   int            `json:"word_count"`
	ReadingTime int            `json:"reading_time"`
}

type listResponse struct {
	Count          int        `json:"count"`
	Results        []Document `json:"results"`
	NextPageCursor *string    `json:"nextPageCursor"`
}

// NewClient creates a Client with the given API token.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ValidateToken checks the token against the Readwise auth endpoint.
// Returns nil if valid, error otherwise.
func (c *Client) ValidateToken(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURLV2+"/auth/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("validate token: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid token (HTTP %d)", resp.StatusCode)
	}
	return nil
}

// ListAll fetches all documents updated after the given time (nil = all).
// Paginates automatically and retries on 429 using the Retry-After header.
func (c *Client) ListAll(ctx context.Context, updatedAfter *time.Time) ([]Document, error) {
	var all []Document
	var cursor *string

	for {
		params := url.Values{}
		if updatedAfter != nil {
			params.Set("updatedAfter", updatedAfter.UTC().Format(time.RFC3339))
		}
		if cursor != nil {
			params.Set("pageCursor", *cursor)
		}

		reqURL := baseURLV3 + "/list/?" + params.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}
		req.Header.Set("Authorization", "Token "+c.token)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("list documents: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			delay := 3 * time.Second
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, err := strconv.Atoi(ra); err == nil {
					delay = time.Duration(secs) * time.Second
				}
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("list documents: HTTP %d: %s", resp.StatusCode, string(body))
		}

		var result listResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decode response: %w", err)
		}
		resp.Body.Close()

		all = append(all, result.Results...)

		if result.NextPageCursor == nil {
			break
		}
		cursor = result.NextPageCursor
	}

	return all, nil
}
