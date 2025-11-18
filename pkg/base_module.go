package api

import "go.uber.org/fx"

// define

var module = fx.Options(
	fx.Provide(
		fx.Annotate(NewGrpcServiceRegistry,
			fx.ParamTags(`group:"grpc"`),
		),
	),
	fx.Provide(NewGrpcServer),

	fx.Provide(
		fx.Annotate(NewHttpHandlersRegistry,
			fx.ParamTags(`group:"http"`),
		),
	),
	fx.Provide(NewHttpServer),

	fx.Invoke(func(_ *GrpcServer) {}),
	fx.Invoke(func(_ *HttpServer) {}),
)

func WithGrpcTag(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(GrpcService)),
		fx.ResultTags(`group:"grpc"`),
	)
}

func WithHttpTag(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(HttpHandler)),
		fx.ResultTags(`group:"http"`),
	)
}

func NewBaseModule(grpcConfig *GrpcConfig, httpConfig *HttpConfig, metricsConfig *MetricsConfig) fx.Option {
	return fx.Options(
		fx.Supply(grpcConfig),
		fx.Supply(httpConfig),
		fx.Supply(metricsConfig),
		module,
	)
}
