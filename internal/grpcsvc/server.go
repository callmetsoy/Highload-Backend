package grpcsvc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sis2/gen"
	"sis2/internal/domain"
)

type Server struct {
	gen.UnimplementedUrlServiceServer
	svc *domain.UrlService
}

func NewServer(svc *domain.UrlService) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateShortUrl(ctx context.Context, req *gen.CreateUrlRequest) (*gen.CreateUrlResponse, error) {
	if req.LongUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "long_url is required")
	}

	res, err := s.svc.Create(ctx, domain.CreateUrlInput{
		LongURL:        req.LongUrl,
		CustomAlias:    req.CustomAlias,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create: %v", err)
	}

	return &gen.CreateUrlResponse{
		ShortUrl:  res.ShortURL,
		ShortCode: res.ShortCode,
	}, nil
}
