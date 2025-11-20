package auth

import "time"

type Config struct {
	JwtSecret   string        `env:"JWT_SECRET,required,unset,notEmpty"`
	JwtTokenTTL time.Duration `env:"JWT_TOKEN_TTL,required,notEmpty"`
}
