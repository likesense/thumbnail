package grpc

import (
	"thumbnail/internal/services"

	pb "thumbnail/internal/proto"

	"google.golang.org/grpc"
)

type Handler struct {
	ThumbnailHandler *ThumbnailHandler
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		ThumbnailHandler: NewThumbnailHandler(services),
	}
}

func (h *Handler) RegisterYoutubeHandler(srv *grpc.Server) {
	pb.RegisterThumbnailServiceServer(srv, h.ThumbnailHandler)
}
