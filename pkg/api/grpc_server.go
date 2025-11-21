package api

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	config Config
	server *grpc.Server
	log    zerolog.Logger
}

type GrpcService interface {
	Descriptor() *grpc.ServiceDesc
}

func NewGrpcServer(registry []GrpcService, log *zerolog.Logger, config Config, opts grpc.ServerOption) *GrpcServer {
	grpcServer := grpc.NewServer(opts)

	for _, svc := range registry {
		grpcServer.RegisterService(svc.Descriptor(), svc)
	}
	return &GrpcServer{
		server: grpcServer,
		log:    log.With().Str("component", "grpc_server").Logger(),
		config: config,
	}
}

func (gr GrpcServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", gr.config.Addr)
	if err != nil {
		gr.log.Error().Err(err).Msg("failed to listen")
	}
	gr.log.Info().Str("addr", gr.config.Addr).Msg("Listening on port!")
	gr.log.Debug().Str("services", fmt.Sprint(gr.server.GetServiceInfo())).Msg("Service info")

	return gr.server.Serve(lis)
}
