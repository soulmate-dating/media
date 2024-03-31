package grpc

import (
	"errors"
	"github.com/TobbyMax/validator"
	"github.com/soulmate-dating/media/internal/app"
	"google.golang.org/grpc/codes"
)

var ErrMissingArgument = errors.New("required argument is missing")

func UploadFileSuccessResponse(link string) *UploadFileResponse {
	return &UploadFileResponse{
		Link: link,
	}
}

func GetErrorCode(err error) codes.Code {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		return codes.InvalidArgument
	case errors.Is(err, app.ErrForbidden):
		return codes.PermissionDenied
	}
	return codes.Internal
}
