package grpcsvc

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sis2/gen"
)

type Client struct {
	stub gen.UrlServiceClient
}

func NewClient(stub gen.UrlServiceClient) *Client {
	return &Client{stub: stub}
}

func (c *Client) CreateWithRetry(ctx context.Context, req *gen.CreateUrlRequest) (*gen.CreateUrlResponse, error) {
	maxAttempts := 3
	baseDelay := 100 * time.Millisecond
	timeout := 500 * time.Millisecond

	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, timeout)
		resp, err := c.stub.CreateShortUrl(attemptCtx, req)
		cancel()

		if err == nil {
			return resp, nil
		}

		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
			// ok
		default:
			return nil, fmt.Errorf("non-retryable: %w", err)
		}

		lastErr = err
		if attempt == maxAttempts-1 {
			break
		}

		window := baseDelay * time.Duration(1<<attempt)
		jitter := time.Duration(rand.Int63n(int64(window)))
		select {
		case <-time.After(jitter):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
