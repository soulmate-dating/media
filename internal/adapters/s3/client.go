package s3

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	NoSuchKey = "NoSuchKey"
)

type Client interface {
	IsObjectExist(ctx context.Context, objectName string) (bool, error)
	Upload(ctx context.Context, objectName, contentType string, content []byte) error
	Download(ctx context.Context, objectName, filePath string) error
	Delete(ctx context.Context, objectName string) error
	GetLinkByName(objectName string) string
}

type client struct {
	client     *minio.Client
	bucketName string
}

func NewClient(endpoint, accessKey, secretKey, bucketName string, secure bool) (Client, error) {
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}
	exists, err := c.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = c.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	}
	if err != nil {
		return nil, err
	}
	return &client{client: c, bucketName: bucketName}, err
}

func (c *client) IsObjectExist(ctx context.Context, objectName string) (bool, error) {
	_, err := c.client.StatObject(ctx, c.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == NoSuchKey {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (c *client) Upload(ctx context.Context, objectName, contentType string, content []byte) error {
	_, err := c.client.PutObject(ctx, c.bucketName, objectName, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (c *client) Download(ctx context.Context, objectName, filePath string) error {
	err := c.client.FGetObject(ctx, c.bucketName, objectName, filePath, minio.GetObjectOptions{})
	return err
}

func (c *client) Delete(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
	return err
}

func (c *client) GetLinkByName(objectName string) string {
	return c.client.EndpointURL().String() + "/" + c.bucketName + "/" + objectName
}
