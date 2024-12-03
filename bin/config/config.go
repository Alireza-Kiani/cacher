package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port int64
}

func GetConfig() Config {
	portStr := os.Getenv("port")
	if len(portStr) == 0 {
		log.Println("port not provided by env, going for default which is 3000")
		portStr = "3000"
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Println("error parsing port")
	}

	return Config{
		Port: port,
	}
}
