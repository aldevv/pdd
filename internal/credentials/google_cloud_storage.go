package credentials

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

const (
	projectID  = "plantdiseasedetection-365617"
	bucketName = "pdd_api"
	uploadPath = "photos/"
)

type UploaderClient struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var GClient *UploaderClient

// UploadFile uploads an object
func (c *UploaderClient) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (c *UploaderClient) RemoveFile(object string) error {
	ctx := context.Background()

	err := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete object: %v", err)
	}

	return nil
}

func init() {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./google_cloud_credentials.json")
	}
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}

	GClient = &UploaderClient{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: uploadPath,
	}

}
