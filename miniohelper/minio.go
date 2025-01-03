package miniohelper

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioHelper struct {
	Client *minio.Client
}

func NewMinioHelper(endpoint, accessKey, secretKey string, useSSL bool) (*MinioHelper, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioHelper{Client: client}, nil
}

func (m *MinioHelper) ListBuckets() ([]minio.BucketInfo, error) {
	buckets, err := m.Client.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}
	return buckets, nil
}

func (m *MinioHelper) UploadFile(bucketName, fileName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = m.Client.PutObject(
		context.Background(),
		bucketName,
		fileName,
		file,
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	return err
}

func (m *MinioHelper) RemoveFile(bucketName, objectName string) error {
	return m.Client.RemoveObject(
		context.Background(),
		bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
}

func (m *MinioHelper) GetFileStream(bucketName, objectName string) (io.Reader, error) {
	object, err := m.Client.GetObject(
		context.Background(),
		bucketName,
		objectName,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get object with name: %v %w", object, err)
	}
	return object, nil
}
