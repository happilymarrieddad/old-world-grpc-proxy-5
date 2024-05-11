package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	authpb "github.com/happilymarrieddad/old-world/grpc-proxy-5/pb/proto/auth"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Auth
	if err := authpb.RegisterAuthHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	// V1

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	port := int(8081)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      cors.AllowAll().Handler(mux),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	log.Printf("Server running on port %d\n", port)
	return s.ListenAndServe()
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
