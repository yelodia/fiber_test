package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	DB *sql.DB
}

// New returns a new Config struct
func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println("database connecting error: ", err)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		log.Println("cannot connect database: ", err)
		os.Exit(1)
	}
	return &Config{
		DB: db,
	}
}

func (c *Config) GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func (c *Config) GetEnvInt(name string, defaultVal int) int {
	valueStr := c.GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func (c *Config) GetEnvBool(name string, defaultVal bool) bool {
	valStr := c.GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func (c *Config) GetEnvSlice(name string, defaultVal []string, sep string) []string {
	valStr := c.GetEnv(name, "")
	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
