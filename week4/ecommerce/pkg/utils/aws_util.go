package utils

import (
	"bytes"
	"ecommerce/pkg/config"
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

func UploadFileToS3(s3Config config.S3Config, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
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
	fileExt := filepath.Ext(header.Filename)
	key := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), fileExt)

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
			ContentLength: aws.Int64(header.Size),
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
