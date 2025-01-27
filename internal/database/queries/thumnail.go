package queries

const (
	CreateThumbnailTable = `
	CREATE TABLE IF NOT EXISTS thumbnail (
		url TEXT PRIMARY KEY,
		thumbnail BLOB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	SaveThumbnail = `
	INSERT INTO thumbnail (url, thumbnail)
	VALUES(?, ?)
	`
	GetThumbnail = `
	SELECT thumbnail FROM thumbnail WHERE url = ?
	`
	ExistsThumnail = `
	SELECT EXISTS(SELECT 1 FROM thumbnail WHERE url = ?)
	`
)
