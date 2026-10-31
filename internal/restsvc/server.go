package restsvc

import (
	"encoding/json"
	"net/http"

	"sis2/internal/domain"
)

type Server struct {
	svc *domain.UrlService
}

func NewServer(svc *domain.UrlService) *Server {
	return &Server{svc: svc}
}

type createReq struct {
	LongURL        string `json:"long_url"`
	CustomAlias    string `json:"custom_alias"`
	IdempotencyKey string `json:"idempotency_key"`
}

type createResp struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}

func (s *Server) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	res, err := s.svc.Create(r.Context(), domain.CreateUrlInput{
		LongURL:        req.LongURL,
		CustomAlias:    req.CustomAlias,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createResp(res))
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/urls", s.CreateURL)
	return mux
}
