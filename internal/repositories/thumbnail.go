package repositories

import (
	"database/sql"
	"fmt"
	"thumbnail/internal/database/queries"
)

type ThumbnailRepository struct {
	db *sql.DB
}

func NewThumbnailRepository(db *sql.DB) *ThumbnailRepository {
	return &ThumbnailRepository{db: db}
}

func (tr *ThumbnailRepository) Save(URL string, thumbnail []byte) error {
	_, err := tr.db.Exec(queries.SaveThumbnail, URL, thumbnail)
	if err != nil {
		return fmt.Errorf("failed to save thumbnail: %s", err.Error())
	}
	return nil
}

func (tr *ThumbnailRepository) Get(URL string) ([]byte, error) {
	var thumbnail []byte
	err := tr.db.QueryRow(queries.GetThumbnail, URL).Scan(&thumbnail)
	if err != nil {
		return nil, fmt.Errorf("failed to get thumbnail: %s", err.Error())
	}
	return thumbnail, nil
}

func (tr *ThumbnailRepository) Exists(URL string) (bool, error) {
	var exists bool
	err := tr.db.QueryRow(queries.ExistsThumnail, URL).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if thumnail exists: %s", err.Error())
	}
	return exists, nil
}
