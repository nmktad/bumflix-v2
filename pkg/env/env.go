package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Endpoint        []byte
	AccessKeyID     []byte
	SecretAccessKey []byte
	UseSSL          bool
}

var EnvInstance *Env

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	endpoint := []byte(os.Getenv("MINIO_ENDPOINT"))
	accessKeyID := []byte(os.Getenv("MINIO_ACCESS_KEY"))
	secretAccessKey := []byte(os.Getenv("MINIO_SECRET_KEY"))
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	EnvInstance = &Env{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		UseSSL:          useSSL,
	}

	return nil
}
