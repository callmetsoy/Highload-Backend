package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"sis2/gen"
	"sis2/internal/domain"
	"sis2/internal/grpcsvc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := domain.NewUrlService()
	s := grpc.NewServer()
	gen.RegisterUrlServiceServer(s, grpcsvc.NewServer(svc))

	log.Println("gRPC server listening on :9090")
	log.Fatal(s.Serve(lis))
}
