package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"thumbnail/internal/database/queries"

	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	DB   *sql.DB
	Path string
}

func NewConnection(path string) *Connection {
	connection := &Connection{
		Path: path,
	}

	if err := connection.connect(); err != nil {
		log.Fatal("Failed connect to sqlite:", err)
	}

	return connection
}

func (c *Connection) connect() error {
	var err error

	if err := os.Remove(c.Path); err != nil && !os.IsNotExist(err) {
		return err
	}

	c.DB, err = sql.Open("sqlite3", c.Path)
	if err != nil {
		return fmt.Errorf("error connecting to the database sqlite: %v", err)
	}

	_, err = c.DB.Exec(queries.CreateThumbnailTable)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func (c *Connection) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
