package misc

import (
	"context"
	"fmt"
	"time"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ImageUploader struct {
	bucketName  string
	imageDomain string
}

func NewImageUploader(bucketName string, imageDomain string) *ImageUploader {
	return &ImageUploader{bucketName: bucketName, imageDomain: imageDomain}
}

// Upload an image to S3
func (i *ImageUploader) GetUploadLink(key string) (string, string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", "", nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws2.String(cfg.Region),
	})
	if err != nil {
		return "", "", nil
	}

	s3Client := s3.New(sess)

	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    aws.String(key),
	})

	expiration := 15 * time.Minute
	url, err := req.Presign(expiration)
	if err != nil {
		return "", "", err
	}

	return url, fmt.Sprintf("https://%s/%s", i.imageDomain, key), nil
}
