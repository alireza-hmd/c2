package client

import (
	"time"

	"github.com/alireza-hmd/c2/pkg/aes"
)

const (
	Connected    = 1
	Disconnected = 2

	Timeout = 5
)

type Client struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Token      string    `json:"token" gorm:"unique"`
	Listener   string    `json:"listener"`
	RemoteIP   string    `json:"remote_ip"`
	ClientType string    `json:"client_type" gorm:"not null"`
	Timeout    int       `json:"timeout" gorm:"default:5"`
	Connected  int       `json:"connected" gorm:"default:2"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Register struct {
	Listener   string `json:"listener"`
	RemoteIP   string `json:"remote_ip"`
	ClientType string `json:"client_type"`
}

func (c *Client) Validate() error {
	return nil
}

func NewClient(listener, ip, clientType string) *Client {
	token := aes.GenerateToken(32)
	return &Client{
		Token:      token,
		Listener:   listener,
		RemoteIP:   ip,
		ClientType: clientType,
		Timeout:    Timeout,
		Connected:  Connected,
	}
}

type ClientUpdate struct {
	Connected int       `json:"connected"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Reader interface {
	Get(token string) (*Client, error)
	List() ([]Client, error)
	ListConnected() ([]Client, error)
}

type Writer interface {
	Create(c *Client) (int, error)
	Update(token string, c *ClientUpdate) error
	Delete(token string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get(token string) (*Client, error)
	List() ([]Client, error)
	Create(c *Client) (int, error)
	Update(token string, c *ClientUpdate) error
	Delete(token string) error
}
