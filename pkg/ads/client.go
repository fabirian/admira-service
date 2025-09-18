package ads

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"admira-service/internal/models"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetPerformance(ctx context.Context) (*models.AdsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var adsResponse models.AdsResponse
	if err := json.Unmarshal(body, &adsResponse); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &adsResponse, nil
}