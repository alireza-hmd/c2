package mysql

import (
	"errors"
	"log"
	"time"

	"github.com/alireza-hmd/c2/clients"
	"github.com/alireza-hmd/c2/pkg/response"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) Get(token string) (*clients.Client, error) {
	var c clients.Client
	res := r.db.Model(&clients.Client{}).Where("token = ?", token).Limit(1).Find(&c)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting client from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return &c, nil
}

func (r *ClientRepository) List() ([]clients.Client, error) {
	var cc []clients.Client
	res := r.db.Model(&clients.Client{}).Find(&cc)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting clients list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return cc, nil
}

func (r *ClientRepository) ListConnected() ([]clients.Client, error) {
	var cc []clients.Client
	res := r.db.Model(&clients.Client{}).Where("connected = ?", clients.Connected).Find(&cc)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting connected clients list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return cc, nil
}

func (r *ClientRepository) Create(c *clients.Client) (int, error) {
	err := r.db.Create(c).Error
	if err != nil {
		log.Println(err)
		return 0, errors.New("error creating client\n")
	}
	return c.ID, nil
}

func (r *ClientRepository) Update(token string, c *clients.ClientUpdate) error {
	c.UpdatedAt = time.Now()
	res := r.db.Model(&clients.Client{}).Where("token = ?", token).Updates(c)
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error updating client\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *ClientRepository) Delete(token string) error {
	res := r.db.Where("token = ?", token).Delete(&clients.Client{})
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error deleting client\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}
