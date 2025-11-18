package api

import (
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server   *grpc.Server
	log      *zerolog.Logger
	registry []GrpcServer
}

type GrpcService interface {
	Descriptor() *grpc.ServiceDesc
}

func NewGrpcServer() {}
