package services

import (
	"github.com/pantyukhov/imageresizeserver/pkg/setting"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Service struct {
	MinioClient *minio.Client
}

func NewS3Service() S3Service {
	endpoint := setting.Settings.S3Config.Endpoint
	accessKeyID := setting.Settings.S3Config.AccessKeyID
	secretAccessKey := setting.Settings.S3Config.SecretAccessKey
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return S3Service{
		MinioClient: minioClient,
	}
}

func (s *S3Service) GetOrCreteFile(filepath string) (*minio.Object, error) {
	n, err := s.MinioClient.GetObject(setting.Settings.Context.Context, setting.Settings.S3Config.Bucket, filepath, minio.GetObjectOptions{})

	return n, err
}
