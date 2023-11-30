package mysql

import (
	"errors"
	"log"
	"time"

	"github.com/alireza-hmd/c2/listener"
	"github.com/alireza-hmd/c2/pkg/response"
	"gorm.io/gorm"
)

type ListenerRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) listener.Repository {
	return &ListenerRepository{
		db: db,
	}
}

func (r *ListenerRepository) Get(name string) (*listener.Listener, error) {
	var l listener.Listener
	res := r.db.Model(&listener.Listener{}).Where("name = ?", name).Limit(1).Find(&l)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting listener from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return &l, nil
}

func (r *ListenerRepository) List() ([]*listener.Listener, error) {
	var ll []*listener.Listener
	res := r.db.Model(&listener.Listener{}).Find(&ll)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting listeners list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return ll, nil
}

func (r *ListenerRepository) ListActive() ([]*listener.Listener, error) {
	var ll []*listener.Listener
	res := r.db.Model(&listener.Listener{}).Where("active = ?", listener.ActiveStatus).Find(&ll)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting active listeners list from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return ll, nil
}

func (r *ListenerRepository) Create(l *listener.Listener) (int, error) {
	err := r.db.Create(l).Error
	if err != nil {
		log.Println(err)
		return 0, errors.New("error creating listener\n")
	}
	return l.ID, nil
}
func (r *ListenerRepository) Update(name string, l *listener.ListenerUptade) error {
	l.UpdatedAt = time.Now()
	res := r.db.Model(&listener.Listener{}).Where("name = ?", name).Updates(l)
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error updating listener\n")
	}
	return nil
}
func (r *ListenerRepository) Delete(name string, stop chan listener.Cancel) error {
	res := r.db.Where("name = ?", name).Delete(&listener.Listener{})
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error deleting listener\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	stop <- listener.Cancel{}
	return nil
}
