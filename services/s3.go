package services

import (
	"github.com/pantyukhov/imageresizeserver/pkg/setting"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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

func (s *S3Service) GetResizeSettings(filepath string) (int, int, string) {
	items := strings.Split(filepath, "/")
	sizes := strings.Split(strings.ToLower(items[len(items)-2]), "x")

	height, err := strconv.Atoi(sizes[0])
	if err != nil {
		height = -1
	}

	width, err := strconv.Atoi(sizes[1])
	if err != nil {
		width = -1
	}
	return height, width
}

func (s *S3Service) ResizeFilePath(filepath string) error {
	height, width, path := s.GetResizeSettings(filepath)

	file, err := s.MinioClient.GetObject(setting.Settings.Context.Context, setting.Settings.S3Config.Bucket, filepath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "resize-")

	if err != nil {
		return err
	}
	func() {
		// make a buffer to keep chunks that are read
		buf := make([]byte, 1024)
		for {
			// read a chunk
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
			}
			if n == 0 {
				break
			}

			// write a chunk
			if _, err := tmpFile.Write(buf[:n]); err != nil {
			}
		}

		// Close the file
		if err := tmpFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	os.Remove(tmpFile.Name())
}

func (s *S3Service) GetOrCreteFile(filepath string) (*minio.Object, *minio.ObjectInfo, error) {
	if strings.HasPrefix(filepath, "/") {
		filepath = filepath[1:]
	}
	file, err := s.MinioClient.GetObject(setting.Settings.Context.Context, setting.Settings.S3Config.Bucket, filepath, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	info, err := file.Stat()
	if info.Key == "" {
	}
	return file, &info, err
}
