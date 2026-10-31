package restsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	url  string
	http *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		url:  baseURL,
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Create(ctx context.Context, longURL, idemKey string) (*createResp, error) {
	reqBody := createReq{LongURL: longURL, IdempotencyKey: idemKey}
	b, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.url+"/api/v1/urls", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var out createResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
