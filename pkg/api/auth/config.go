package auth

type Config struct {
	JWTSecret string `env:"JWT_SECRET,required,unset,notEmpty"`
}
