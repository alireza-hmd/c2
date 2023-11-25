package listener

import (
	"errors"
	"time"
)

const (
	ActiveStatus   = 1
	DeactiveStatus = 2

	HTTPConnection   = "http"
	SocketConnection = "socket"
)

type Cancel struct{}

type Listener struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null;unique"`
	Port       uint      `json:"port" gorm:"not null"`
	IpAddress  string    `json:"ip_address"`
	Connection string    `json:"connection" gorm:"not null"`
	Active     int       `json:"active" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (l *Listener) Validate() error {
	if l.Port < 1 && l.Port > 65536 {
		return errors.New("invalid port number\n")
	}
	if l.Connection != HTTPConnection && l.Connection != SocketConnection {
		return errors.New("undefined connection\n")
	}
	return nil
}

func NewListener(name, connection string, port uint) *Listener {
	return &Listener{
		Name:       name,
		Port:       port,
		Connection: connection,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

type ListenerUptade struct {
	Active    int       `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Listener) TableName() string {
	return "listeners"
}

type Reader interface {
	Get(name string) (*Listener, error)
	List() ([]Listener, error)
	ListActive() ([]Listener, error)
}

type Writer interface {
	Create(l *Listener) (int, error)
	Update(name string, l *ListenerUptade) error
	Delete(name string, stop chan Cancel) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get(name string) (*Listener, error)
	List() ([]Listener, error)
	Create(l *Listener, stop map[int](chan Cancel)) (int, error)
	Update(name string, l *ListenerUptade) error
	Delete(name string, stop chan Cancel) error
	Activation(l *Listener, stop chan Cancel, status int) error
	Run(l *Listener, stop chan Cancel)
	RunActiveListeners(stop map[int](chan Cancel)) error
}
