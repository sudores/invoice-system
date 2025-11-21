package api

type Config struct {
	Addr string `env:"Addr" envDefault:":50051"`
}
