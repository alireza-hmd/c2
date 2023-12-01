package clients

import (
	"time"

	"github.com/alireza-hmd/c2/pkg/encrypt/aes"
)

const (
	Connected    = 1
	Disconnected = 2

	Timeout = 10
)

type Client struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Token      string    `json:"token" gorm:"unique"`
	Listener   string    `json:"listener"`
	RemoteIP   string    `json:"remote_ip"`
	ClientType string    `json:"client_type" gorm:"not null"`
	Timeout    int       `json:"timeout"`
	Connected  int       `json:"connected" gorm:"default:2"`
	Encrypted  bool      `json:"encrypted" gorm:"default:false"`
	SilentMode bool      `json:"silent_mode" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Register struct {
	Listener   string `json:"listener"`
	RemoteIP   string `json:"remote_ip"`
	SilentMode bool   `json:"silent_mode"`
	Encrypted  bool   `json:"encrypted"`
	ClientType string `json:"client_type"`
}

func (c *Client) Validate() error {
	return nil
}

func NewClient(listener, ip, clientType string, silentMode, encrypted bool) *Client {
	token := aes.GenerateToken(16)
	return &Client{
		Token:      token,
		Listener:   listener,
		RemoteIP:   ip,
		ClientType: clientType,
		SilentMode: silentMode,
		Encrypted:  encrypted,
		Timeout:    Timeout,
		Connected:  Connected,
	}
}

type ClientUpdate struct {
	Connected int       `json:"connected"`
	Timeout   int       `json:"timeout"`
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
	DeleteListenerClient(name string) error
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
	DeleteListenerClient(name string) error
}
