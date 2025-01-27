package services

import (
	"thumbnail/internal/repositories"
	"thumbnail/internal/utils"
)

type Thumbnail interface {
	GetThumbnail(URL string) ([]byte, error)
	GetThumbnailAsync(URLs []string) (<-chan utils.ThumbnailResult, error)
}

type Services struct {
	Thumbnail Thumbnail
}

func NewServices(repos *repositories.Repositories) *Services {
	return &Services{
		Thumbnail: NewThumbnailService(repos),
	}
}
