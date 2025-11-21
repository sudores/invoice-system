package api

type Config struct {
	GrpcAddr string `env:"GRPC_ADDR" envDefault:":50051"`
	HttpAddr string `env:"HTTP_ADDR" envDefault:":8080"`
}
