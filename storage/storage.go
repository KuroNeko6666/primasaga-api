package storage

import (
	"log"

	"github.com/minio/minio-go"
)

func ConnectStorage() (*minio.Client, error) {
	endpoint := "localhost:9090"
	accessKeyID := "LcEP4PU3l0kNEUz2"
	secretAccessKey := "9k9R56kSZuNWIms0yisBX7D6uK6tTCEU"
	ssl := false

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, ssl)
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient, err
}
