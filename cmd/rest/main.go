package main

import (
	"log"
	"net/http"

	"sis2/internal/domain"
	"sis2/internal/restsvc"
)

func main() {
	svc := domain.NewUrlService()
	srv := restsvc.NewServer(svc)
	log.Println("REST server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv.Router()))
}
