package config

import (
	"encoding/json"
	"log"
	"os"

	"avito-backend-trainee-assignment/internal/models"
)

func LoadConfig(filename string) models.Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	return config
}
