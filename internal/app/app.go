package app

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/soulmate-dating/media/internal/adapters/s3"
)

var (
	ErrForbidden = fmt.Errorf("forbidden")
)

type App interface {
	UploadFile(ctx context.Context, name string, data []byte) (string, error)
}

type Application struct {
	client s3.Client
}

func (a Application) UploadFile(ctx context.Context, contentType string, data []byte) (string, error) {
	hash := sha256.Sum256(data)
	hashString := fmt.Sprintf("%x.png", hash)
	err := a.client.Upload(ctx, hashString, contentType, data)
	if err != nil {
		return "", err
	}
	return a.client.GetLinkByName(hashString), nil
}

func NewApp(client s3.Client) App {
	return &Application{client: client}
}
