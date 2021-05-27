package internal

import (
	"os"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

func validateMinioSettings() {
	var variables = []string{"MINIO_ENDPOINT", "MINIO_USER", "MINIO_PASSWORD"}
	for _, k := range variables {
		if v := os.Getenv(k); v == "" {
			log.Fatalf("you must specify $%s", k)
		}
	}
}

func newMinio(bucketName string) *minio.Client {
	validateMinioSettings()
	useSSL := !(os.Getenv("MINIO_DISABLE_SSL") == "true")
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASSWORD"), useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// ensure given bucket exists
	makeBucket(minioClient, bucketName)
	return minioClient
}

func makeBucket(client *minio.Client, name string) {
	ok, err := client.BucketExists(name)
	if err != nil {
		log.Fatalln(err)
	}

	switch ok {
	case true: // bucket exist
		return
	case false: // bucket does not exist
		if err := client.MakeBucket(name, ""); err != nil {
			log.Fatalln(err)
		}
	}
}
