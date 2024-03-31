package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MediaService) UploadFile(ctx context.Context, request *UploadFileRequest) (*UploadFileResponse, error) {
	if request.GetData() == nil {
		return nil, status.Error(codes.InvalidArgument, "empty file")
	}
	link, err := s.app.UploadFile(ctx, request.GetContentType(), request.GetData())
	if err != nil {
		return nil, status.Error(GetErrorCode(err), err.Error())
	}
	return UploadFileSuccessResponse(link), nil
}

func (s *MediaService) mustEmbedUnimplementedMediaServiceServer() {}
