package env

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once    sync.Once
	loadErr error
)

func load() error {
	once.Do(func() {
		loadErr = godotenv.Load(".env")
	})
	return loadErr
}

func GetString(key, fallback string) string {
	if err := load(); err != nil {
		log.Print("Warning: .env not found; Using fallback variables")
	}
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) (int, bool) {
	if err := load(); err != nil {
		log.Print("Warning: .env not found; Using fallback variables")
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback, false
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback, false
	}
	return valAsInt, true
}
