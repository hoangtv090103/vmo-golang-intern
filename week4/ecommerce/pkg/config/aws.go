package config

type S3Config struct {
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
}

func NewS3Config(region, bucketName, accessKey, secretKey string) *S3Config {
	return &S3Config{
		Region:     region,
		BucketName: bucketName,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
	}
}
