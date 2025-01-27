package repositories

import "database/sql"

type Thumbnail interface {
	Save(URL string, thumbnail []byte) error
	Get(URL string) ([]byte, error)
	Exists(URL string) (bool, error)
}

type Repositories struct {
	Thumbnail Thumbnail
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Thumbnail: NewThumbnailRepository(db),
	}
}
