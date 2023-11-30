package misc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

type ImageUploader struct {
	bucketName  string
	imageDomain string
}

func NewImageUploader(bucketName string, imageDomain string) *ImageUploader {
	return &ImageUploader{bucketName: bucketName, imageDomain: imageDomain}
}

// Decode an image
func decodeImage(reader io.Reader) (*image.Image, error) {
	img, _, err := image.Decode(reader)

	return &img, err
}

// Resize and crop the image
func resizeAndCropImage(img *image.Image, width uint, height uint) image.Image {
	resizedImg := resize.Resize(width, 0, *img, resize.Lanczos3)

	cropWidth := width
	cropHeight := height

	if resizedImg.Bounds().Dx() < int(width) {
		cropWidth = uint(resizedImg.Bounds().Dx())
	}

	if resizedImg.Bounds().Dy() < int(height) {
		cropHeight = uint(resizedImg.Bounds().Dy())
	}

	croppedImg := imaging.CropCenter(resizedImg, int(cropWidth), int(cropHeight))

	return croppedImg
}

// Upload an image to S3
func uploadImageS3(imageDomain string, prefix string, key string, bucketName string, img *image.Image) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws2.String(cfg.Region),
	})
	if err != nil {
		return "", nil
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s/%s", prefix, key)))
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%s.jpeg", hex.EncodeToString(hashBytes))

	s3Client := s3.New(sess)

	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, *img, nil); err != nil {
		return "", err
	}

	imageData := buffer.Bytes()
	bufferReader := bytes.NewReader(imageData)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws2.String(bucketName),
		Key:    aws2.String(hashString),
		Body:   bufferReader,
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s", imageDomain, hashString), nil
}

// Upload a file
func (i *ImageUploader) ResizeAndUploadFile(reader io.Reader, width uint, height uint, prefix string, key string) (string, error) {
	img, err := decodeImage(reader)
	if err != nil {
		return "", err
	}

	edited := resizeAndCropImage(img, width, height)

	return uploadImageS3(i.imageDomain, prefix, key, i.bucketName, &edited)
}
