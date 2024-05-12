package s3

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"time"
)

const (
	NoSuchKey = "NoSuchKey"
)

type Client interface {
	IsObjectExist(ctx context.Context, objectName string) (bool, error)
	Upload(ctx context.Context, objectName, contentType string, content []byte) error
	Download(ctx context.Context, objectName, filePath string) error
	Delete(ctx context.Context, objectName string) error
	GetURL(ctx context.Context, objectName string) *url.URL
}

type Config struct {
	Address        string
	AccessKey      string
	SecretKey      string
	SessionToken   string
	BucketName     string
	ExpiryDuration time.Duration
	Secure         bool
	Policy         string
}

type client struct {
	client     *minio.Client
	bucketName string
	expiryTime time.Duration
}

func NewClient(ctx context.Context, cfg Config) (Client, error) {
	log.Printf("setting up s3 client")
	c, err := minio.New(cfg.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, cfg.SessionToken),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, err
	}
	exists, err := c.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = c.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{})
	}
	if cfg.Policy != "" {
		log.Printf("applying policy")
		err = c.SetBucketPolicy(ctx, cfg.BucketName, cfg.Policy)
	}

	if err != nil {
		return nil, err
	}
	return &client{client: c, bucketName: cfg.BucketName, expiryTime: cfg.ExpiryDuration}, err
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

func (c *client) GetURL(_ context.Context, objectName string) *url.URL {
	return c.client.EndpointURL().JoinPath(c.bucketName, objectName)
}
