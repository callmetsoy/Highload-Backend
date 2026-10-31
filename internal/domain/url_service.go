package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrAliasTaken = errors.New("alias already taken")

type CreateUrlInput struct {
	LongURL        string
	CustomAlias    string
	IdempotencyKey string
}

type CreateUrlOutput struct {
	ShortURL  string
	ShortCode string
}

type UrlService struct {
	mu      sync.Mutex
	counter int64
	idem    map[string]CreateUrlOutput
	aliases map[string]struct{}
}

func NewUrlService() *UrlService {
	return &UrlService{
		idem:    make(map[string]CreateUrlOutput),
		aliases: make(map[string]struct{}),
	}
}

func (s *UrlService) Create(ctx context.Context, in CreateUrlInput) (CreateUrlOutput, error) {
	if err := ctx.Err(); err != nil {
		return CreateUrlOutput{}, fmt.Errorf("context cancelled: %w", err)
	}

	if in.IdempotencyKey != "" {
		if cached, ok := s.idem[in.IdempotencyKey]; ok {
			return cached, nil
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var code string
	if in.CustomAlias != "" {
		if _, taken := s.aliases[in.CustomAlias]; taken {
			return CreateUrlOutput{}, fmt.Errorf("%w: %s", ErrAliasTaken, in.CustomAlias)
		}
		code = in.CustomAlias
		s.aliases[code] = struct{}{}
	} else {
		s.counter++
		code = fmt.Sprintf("c%06d", s.counter) // Упрощенный Base62
	}

	out := CreateUrlOutput{
		ShortURL:  "https://short.ly/" + code,
		ShortCode: code,
	}

	if in.IdempotencyKey != "" {
		s.idem[in.IdempotencyKey] = out
	}

	return out, nil
}
