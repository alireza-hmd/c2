package mysql

import (
	"errors"
	"log"
	"time"

	"github.com/alireza-hmd/c2/listeners"
	"github.com/alireza-hmd/c2/pkg/response"
	"gorm.io/gorm"
)

type ListenerRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) listeners.Repository {
	return &ListenerRepository{
		db: db,
	}
}

func (r *ListenerRepository) Get(name string) (*listeners.Listener, error) {
	var l listeners.Listener
	res := r.db.Model(&listeners.Listener{}).Where("name = ?", name).Limit(1).Find(&l)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting listener from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return &l, nil
}

func (r *ListenerRepository) List() ([]*listeners.Listener, error) {
	var ll []*listeners.Listener
	res := r.db.Model(&listeners.Listener{}).Find(&ll)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting listeners list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return ll, nil
}

func (r *ListenerRepository) ListActive() ([]*listeners.Listener, error) {
	var ll []*listeners.Listener
	res := r.db.Model(&listeners.Listener{}).Where("active = ?", listeners.ActiveStatus).Find(&ll)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting active listeners list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return ll, nil
}

func (r *ListenerRepository) Create(l *listeners.Listener) (int, error) {
	err := r.db.Create(l).Error
	if err != nil {
		log.Println(err)
		return 0, errors.New("error creating listener\n")
	}
	return l.ID, nil
}
func (r *ListenerRepository) Update(name string, l *listeners.ListenerUptade) error {
	l.UpdatedAt = time.Now()
	res := r.db.Model(&listeners.Listener{}).Where("name = ?", name).Updates(l)
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error updating listener\n")
	}
	return nil
}
func (r *ListenerRepository) Delete(name string, stop chan listeners.Cancel) error {
	res := r.db.Where("name = ?", name).Delete(&listeners.Listener{})
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error deleting listener\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	stop <- listeners.Cancel{}
	return nil
}
