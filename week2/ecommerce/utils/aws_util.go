package utils

import (
	"bytes"
	"ecommerce/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

func UploadFileToS3(s3Config config.S3Config, file multipart.File, handler *multipart.FileHeader) (string, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s3Config.Region),
			Credentials: credentials.NewStaticCredentials(
				s3Config.AccessKey,
				s3Config.SecretKey,
				"",
			),
		})

	if err != nil {
		return "", err
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Create a unique file name using timestamp
	fileExt := filepath.Ext(handler.Filename)
	key := fmt.Sprintf("product_images/%d%s", time.Now().UnixNano(), fileExt)

	// Read file content into buffer
	buffer := new(bytes.Buffer)

	_, err = buffer.ReadFrom(file)

	if err != nil {
		return "", err
	}

	// Prepare S3 upload input

	_, err = svc.PutObject(
		&s3.PutObjectInput{
			Bucket:        aws.String(s3Config.BucketName),
			Key:           aws.String(key),
			Body:          bytes.NewReader(buffer.Bytes()),
			ContentLength: aws.Int64(handler.Size),
			ContentType:   aws.String(http.DetectContentType(buffer.Bytes())),
			//ACL: aws.String("public-read"), // Make file publicly accessible

		},
	)

	if err != nil {
		return "", fmt.Errorf("Failed to upload file to S3: %w", err)
	}

	// Return the S3 URL for the uploaded image
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Config.BucketName, s3Config.Region, key)
	return fileURL, nil
}
