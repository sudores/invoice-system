package main

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/sudores/invoice-system/pkg/api/invoice"
	"google.golang.org/grpc"
)

func main() {
	// 1. Listen on TCP port
	lis, err := net.Listen("tcp", ":50051") // choose any free port
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	// 2. Create gRPC server
	grpcServer := grpc.NewServer(
	// you can attach interceptors here, e.g., JWT
	// grpc.UnaryInterceptor(jwtInterceptor("my-secret")),
	)

	// 3. Register your service
	svc := invoice.NewInvoicesGrpcService(&log.Logger)
	invoice.RegisterInvoiceServiceServer(grpcServer, &svc)

	log.Debug().Msg("gRPC server listening on :50051")
	fmt.Println(grpcServer.GetServiceInfo())

	// 4. Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
