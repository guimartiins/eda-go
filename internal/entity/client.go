package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidName = errors.New("invalid name")
var ErrInvalidEmail = errors.New("invalid email")

type Client struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewClient(name string, email string) (*Client, error) {
	client := &Client{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := client.Validate(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Validate() error {
	if c.Name == "" {
		return ErrInvalidName
	}
	if c.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}
