package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser       string
	DBPassword   string
	DBAddress    string
	DBPort       int64
	DBName       string
	DBConnection string

	Port string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbAddress := getEnv("DB_ADDRESS", "localhost")
	dbPort := getEnvAsInt("DB_PORT", 5432)
	dbName := getEnv("DB_NAME", "postgres")

	dbConnection := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPassword, dbAddress, dbPort, dbName)

	Port := getEnv("PORT", "3000")

	return Config{
		DBUser:       dbUser,
		DBPassword:   dbPassword,
		DBAddress:    dbAddress,
		DBPort:       dbPort,
		DBName:       dbName,
		DBConnection: dbConnection,
		Port:         Port,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
