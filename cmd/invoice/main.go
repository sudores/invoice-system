package main

import (
	"log"
	"net"

	pb "github.com/sudores/invoice-system"
	"google.golang.org/grpc"
)

func main() {
	// 1. Listen on TCP port
	lis, err := net.Listen("tcp", ":50051") // choose any free port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Create gRPC server
	grpcServer := grpc.NewServer(
	// you can attach interceptors here, e.g., JWT
	// grpc.UnaryInterceptor(jwtInterceptor("my-secret")),
	)

	// 3. Register your service
	svc := &invoices.{
		svc: &InvoiceService{
			repo: NewInvoiceRepository(db), // assuming you have a DB connection
		},
	}
	pb.RegisterInvoiceServiceServer(grpcServer, svc)

	log.Println("gRPC server listening on :50051")

	// 4. Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
