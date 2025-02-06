package v1

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
)

func DefaultConfig() (*Config, error) {
	apiToken := os.Getenv("GOKAKASHI_API_SERVER_TOKEN")
	logAPITokenOnStartup := false
	if apiToken == "" {
		generatedToken, err := generateToken(32)
		if err != nil {
			return nil, fmt.Errorf("error generating an api token: %w", err)
		}
		apiToken = generatedToken
		logAPITokenOnStartup = true
	}

	apiHost := os.Getenv("GOKAKASHI_SERVER_HOST")
	apiPort, _ := strconv.Atoi(os.Getenv("GOKAKASHI_SERVER_PORT"))
	if apiPort == 0 {
		apiPort = 5556
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	if dbPort == 0 {
		dbPort = 5432
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	return &Config{
		Site: SiteConfig{
			APIToken:             apiToken,
			LogAPITokenOnStartup: logAPITokenOnStartup,
			Host:                 apiHost,
			Port:                 apiPort,
		},
		Database: DbConnection{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
		},
	}, nil
}

func generateToken(length int) (string, error) {
	token := make([]byte, length)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
