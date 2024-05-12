package app

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/soulmate-dating/media/internal/config"
	"log"
	"strings"

	"github.com/soulmate-dating/media/internal/adapters/s3"
)

var (
	ErrForbidden = fmt.Errorf("forbidden")
)

type App interface {
	UploadFile(ctx context.Context, name string, data []byte) (string, error)
}

type Application struct {
	publicHost string
	client     s3.Client
}

func (a Application) UploadFile(ctx context.Context, contentType string, data []byte) (string, error) {
	hash := sha256.Sum256(data)
	hashString := fmt.Sprintf("%x", hash)
	err := a.client.Upload(ctx, hashString, contentType, data)
	if err != nil {
		return "", err
	}
	url := a.client.GetURL(ctx, hashString)
	return strings.Replace(url.String(), url.Hostname(), a.publicHost, 1), nil
}

func New(ctx context.Context, cfg config.Config) App {
	client, err := s3.NewClient(ctx, s3.Config{
		Address:      cfg.S3.Address,
		AccessKey:    cfg.S3.AccessKey,
		SecretKey:    cfg.S3.SecretKey,
		SessionToken: cfg.S3.SessionToken,
		BucketName:   cfg.S3.Bucket,
		Policy:       cfg.S3.Policy,
		Secure:       cfg.S3.Secure,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Application{client: client, publicHost: cfg.API.PublicHost}
}
