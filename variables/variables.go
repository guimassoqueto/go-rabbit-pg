package variables

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

var POSTGRES_PORT = getEnv("POSTGRES_PORT")
var POSTGRES_DB = getEnv("POSTGRES_DB")
var POSTGRES_USER = getEnv("POSTGRES_USER")
var POSTGRES_PASSWORD = getEnv("POSTGRES_PASSWORD")
var POSTGRES_HOST = getEnv("POSTGRES_HOST")
var POSTGRES_TABLE_NAME = getEnv("POSTGRES_TABLE_NAME")

var RABBITMQ_DEFAULT_USER = getEnv("RABBITMQ_DEFAULT_USER")
var RABBITMQ_DEFAULT_PASS = getEnv("RABBITMQ_DEFAULT_PASS")
var RABBITMQ_DEFAULT_HOST = getEnv("RABBITMQ_DEFAULT_HOST")
var RABBITMQ_DEFAULT_PORT = getEnv("RABBITMQ_DEFAULT_PORT")
var RABBITMQ_RECEIVE_QUEUE = getEnv("RABBITMQ_RECEIVE_QUEUE")
var RABBITMQ_SEND_QUEUE = getEnv("RABBITMQ_SEND_QUEUE")
var RABBITMQ_URL = fmt.Sprintf(
	"amqp://%s:%s@%s:%s",
	RABBITMQ_DEFAULT_USER,
	RABBITMQ_DEFAULT_PASS,
	RABBITMQ_DEFAULT_HOST,
	RABBITMQ_DEFAULT_PORT,
)

func MAX_CONCURRENCY() int {
	var max = getEnv("MAX_CONCURRENCY")
	maxConcurrency, err := strconv.Atoi(max)
	if err != nil {
		return 24
	}
	return maxConcurrency
}
