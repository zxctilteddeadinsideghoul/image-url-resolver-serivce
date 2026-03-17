package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	OurDomain       string
	GoldAppleDomain string
	ExistsTimeoutMs int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	timeout, _ := strconv.Atoi(getEnv("EXISTS_TIMEOUT_MS", "500"))

	return &Config{
		Port:            getEnv("PORT", "8080"),
		OurDomain:       getEnv("OUR_DOMAIN", ""),
		GoldAppleDomain: getEnv("GOLDAPPLE_DOMAIN", ""),
		ExistsTimeoutMs: timeout,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
