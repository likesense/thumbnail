package grpc

import (
	"context"
	pb "thumbnail/internal/proto"
	"thumbnail/internal/services"
)

type ThumbnailHandler struct {
	services *services.Services
	pb.UnimplementedThumbnailServiceServer
}

func NewThumbnailHandler(services *services.Services) *ThumbnailHandler {
	return &ThumbnailHandler{
		services: services,
	}
}

func (th *ThumbnailHandler) GetThumbnail(ctx context.Context, req *pb.GetThumbnailRequest) (*pb.ThumbnailResponse, error) {
	thumbnail, err := th.services.Thumbnail.GetThumbnail(req.Url)
	if err != nil {
		return &pb.ThumbnailResponse{
			Url:   req.Url,
			Error: err.Error(),
		}, nil
	}

	return &pb.ThumbnailResponse{
		Url:       req.Url,
		Thumbnail: thumbnail,
	}, nil
}

func (th *ThumbnailHandler) GetThumbnailAsync(req *pb.GetThumbnailsRequestAsync, stream pb.ThumbnailService_GetThumbnailAsyncServer) error {
	results, err := th.services.Thumbnail.GetThumbnailAsync(req.Urls)
	if err != nil {
		return err
	}

	for result := range results {
		response := &pb.ThumbnailResponse{
			Url:       result.URL,
			Thumbnail: result.Thumbnail,
			Error:     result.Error,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return nil
}
