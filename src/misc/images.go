package misc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

type ImageUploader struct {
	bucketName string
	prefix     string
}

func NewImageUploader(bucketName string, prefix string) *ImageUploader {
	return &ImageUploader{bucketName: bucketName, prefix: prefix}
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
func uploadImageS3(prefix string, key string, bucketName string, img *image.Image) (string, error) {
	sess, err := session.NewSession()
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
		Bucket: aws.String(bucketName),
		Key:    aws.String(hashString),
		Body:   bufferReader,
	})

	if err != nil {
		return "", err
	}

	return hashString, nil
}

// Upload a file
func (i *ImageUploader) ResizeAndUploadFile(reader io.Reader, width uint, height uint, key string) (string, error) {
	img, err := decodeImage(reader)
	if err != nil {
		return "", err
	}

	edited := resizeAndCropImage(img, width, height)

	return uploadImageS3(i.prefix, key, i.bucketName, &edited)
}
