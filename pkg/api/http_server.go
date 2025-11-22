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

func NewHttpServer(registry []GrpcService, log *zerolog.Logger, config Config, opts ...runtime.ServeMuxOption) *HttpServer {
	mux := runtime.NewServeMux(opts...)
	handler := corsMiddleware(mux)
	return &HttpServer{
		server: &http.Server{
			Addr:    config.HttpAddr,
			Handler: handler,
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

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			// Preflight request; respond immediately
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
