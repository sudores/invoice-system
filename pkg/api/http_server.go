package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
)

type HttpServer struct {
	server   *http.Server
	mux      *runtime.ServeMux
	log      zerolog.Logger
	registry []GrpcService
	config   Config
}

func NewHttpServer(registry []GrpcService, log *zerolog.Logger, config Config) *HttpServer {
	mux := runtime.NewServeMux()
	return &HttpServer{
		server: &http.Server{
			Addr:    config.HttpAddr,
			Handler: mux,
		}, // TODO: May be add timeouts some time
		mux:      mux,
		log:      log.With().Str("component", "http_server").Logger(),
		registry: registry,
		config:   config,
	}
}

func (hs *HttpServer) ListenAndServe(ctx context.Context) error {
	for _, k := range hs.registry {
		hs.log.Debug().Msg("Registered handler")
		fmt.Printf("%+v", hs)
		k.RegisterHttp(ctx, hs.mux)
	}
	if err := hs.server.ListenAndServe(); err != nil {
		return err
	}
	return hs.server.Shutdown(ctx)
}
