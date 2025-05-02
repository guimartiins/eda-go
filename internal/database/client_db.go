package database

import (
	"database/sql"

	"github.com/guimartiins/eda-go/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{DB: db}
}

func (c *ClientDB) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}

	stmt, err := c.DB.Prepare("SELECT id, name, email, created_at FROM clients WHERE id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(id)

	if err := row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientDB) Save(client *entity.Client) error {
	stmt, err := c.DB.Prepare("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		println("1 - Saving client to database:", err)

		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt, client.UpdatedAt)
	if err != nil {
		println("2 - Saving client to database:", err)
		return err
	}

	return nil
}
