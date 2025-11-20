package repo

type Config struct {
	DBURL string `env:"DB_URL,required,unset,notEmpty"`
}
