package services

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nfnt/resize"
	"github.com/pantyukhov/imageresizeserver/pkg/setting"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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

func (s *S3Service) GetResizeSettings(filepath string) (uint, uint, string) {
	items := strings.Split(filepath, "/")
	sizes := strings.Split(strings.ToLower(items[len(items)-2]), "x")

	height, err := strconv.Atoi(sizes[0])
	if err != nil {
		height = 0
	}

	width, err := strconv.Atoi(sizes[1])
	if err != nil {
		width = 0
	}
	path := strings.Join(items[:len(items)-2], "/") + "/" + items[len(items)-1]
	return uint(height), uint(width), path
}

func (s *S3Service) ResizeImage(localPath string, height uint, width uint) (image.Image, error) {
	file, err := os.Open(localPath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	file.Close()
	if err != nil {
		return nil, err
	}
	m := resize.Resize(width, height, img, resize.Lanczos3)

	return m, err

}

func (s *S3Service) ResizeFilePath(filepath string) error {
	height, width, path := s.GetResizeSettings(filepath)
	file, err := s.MinioClient.GetObject(setting.Settings.Context.Context, setting.Settings.S3Config.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "source-")

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

	m, err := s.ResizeImage(tmpFile.Name(), height, width)
	os.Remove(tmpFile.Name())

	tmpFile, err = ioutil.TempFile(os.TempDir(), "resize-")
	if err != nil {
		return err
	}

	out, err := os.Create(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	_, err = s.MinioClient.FPutObject(setting.Settings.Context.Context, setting.Settings.S3Config.Bucket, filepath, tmpFile.Name(), minio.PutObjectOptions{
		ContentType: info.ContentType,
	})
	os.Remove(tmpFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) GetOrCreteFile(filepath string, allow_resize bool) (*minio.Object, *minio.ObjectInfo, error) {
	if strings.HasPrefix(filepath, "/") {
		filepath = filepath[1:]
	}
	file, err := s.MinioClient.GetObject(setting.Settings.Context.Context, setting.Settings.S3Config.Bucket, filepath, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	info, err := file.Stat()
	if info.Key == "" && allow_resize {
		err := s.ResizeFilePath(filepath)
		if err != nil {
			return nil, nil, err
		}
		return s.GetOrCreteFile(filepath, false)
	}
	return file, &info, err
}
