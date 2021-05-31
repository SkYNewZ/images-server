package minio

import (
	"os"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	*minio.Client

	// BucketName is the default bucket used by this app
	BucketName string
}

func validateMinioSettings() {
	var variables = []string{"MINIO_ENDPOINT", "MINIO_USER", "MINIO_PASSWORD"}
	for _, k := range variables {
		if v := os.Getenv(k); v == "" {
			log.Fatalf("you must specify $%s", k)
		}
	}
}

// New creates a new minio.Client
func New(bucketName string) *Client {
	validateMinioSettings()
	useSSL := !(os.Getenv("MINIO_DISABLE_SSL") == "true")
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASSWORD"), useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	c := &Client{
		Client:     minioClient,
		BucketName: bucketName,
	}

	// ensure given bucket exists
	c.makeBucket()
	return c
}

func (c *Client) makeBucket() {
	ok, err := c.BucketExists(c.BucketName)
	if err != nil {
		log.Fatalln(err)
	}

	switch ok {
	case true: // bucket exist
		return
	case false: // bucket does not exist
		if err := c.MakeBucket(c.BucketName, ""); err != nil {
			log.Fatalln(err)
		}
	}
}
