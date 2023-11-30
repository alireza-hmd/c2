package mysql

import (
	"errors"
	"log"
	"time"

	"github.com/alireza-hmd/c2/client"
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

func (r *ClientRepository) Get(token string) (*client.Client, error) {
	var c client.Client
	res := r.db.Model(&client.Client{}).Where("token = ?", token).Limit(1).Find(&c)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting client from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return &c, nil
}

func (r *ClientRepository) List() ([]client.Client, error) {
	var cc []client.Client
	res := r.db.Model(&client.Client{}).Find(&cc)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting clients list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return cc, nil
}

func (r *ClientRepository) ListConnected() ([]client.Client, error) {
	var cc []client.Client
	res := r.db.Model(&client.Client{}).Where("connected = ?", client.Connected).Find(&cc)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting connected clients list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return cc, nil
}

func (r *ClientRepository) Create(c *client.Client) (int, error) {
	err := r.db.Create(c).Error
	if err != nil {
		log.Println(err)
		return 0, errors.New("error creating client\n")
	}
	return c.ID, nil
}

func (r *ClientRepository) Update(token string, c *client.ClientUpdate) error {
	c.UpdatedAt = time.Now()
	res := r.db.Model(&client.Client{}).Where("token = ?", token).Updates(c)
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
	res := r.db.Where("token = ?", token).Delete(&client.Client{})
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error deleting client\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}
