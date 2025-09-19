package product_s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
)

type Minio struct {
	mc         *minio.Client
	bucketName string
}

func New(mc *minio.Client, bucketName string) *Minio {
	return &Minio{
		mc:         mc,
		bucketName: bucketName,
	}
}

func (m *Minio) SaveImage(ctx context.Context, id string, image []byte) (string, error) {
	const op = "repository.minio.product.SaveImage"

	reader := bytes.NewReader(image)

	_, err := m.mc.PutObject(
		ctx,
		m.bucketName,
		id,
		reader,
		int64(len(image)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	url, err := m.GetImageUrl(ctx, id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (m *Minio) GetImageUrl(ctx context.Context, imageId string) (string, error) {
	const op = "repository.minio.product.GetImageUrl"

	url, err := m.mc.PresignedGetObject(ctx, m.bucketName, imageId, 5*24*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url.String(), nil
}

func (m *Minio) DeleteImage(ctx context.Context, imageId string) error {
	const op = "repository.minio.product.DeleteImage"

	err := m.mc.RemoveObject(ctx, m.bucketName, imageId, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
