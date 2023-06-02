package envs

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func Get(key, def string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	return def
}
