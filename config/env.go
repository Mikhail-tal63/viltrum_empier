package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MongoURI      string
	host          string
	DBName        string
	JWTSecret     string
	JWTexpiration int64
}

var ENVs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://127.0.0.1:27017"),
		DBName:        getEnv("DB_NAME", "viltrum_empier"),
		Port:          getEnv("PORT", "8080"),
		host:          getEnv("HOST", "localhost"),
		JWTexpiration: getEnvInt("JWT_EXPIRATION", 3600*24*7),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
	}
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value

}
func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 14)
		if err != nil {
			return fallback
		}
		return i

	}
	return fallback

}
