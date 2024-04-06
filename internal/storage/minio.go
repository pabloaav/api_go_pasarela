package storage

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage struct {
	S3 *s3.S3
}

func NewS3Session() Storage {
	minioSession := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.AWS_KEY, config.AWS_SECRET, ""),
		Region:      aws.String(config.AWS_REGION),
		// Endpoint:         aws.String(config.MINIO_ENDPOINT),
		// DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	}))
	return Storage{
		S3: s3.New(minioSession),
	}
}
