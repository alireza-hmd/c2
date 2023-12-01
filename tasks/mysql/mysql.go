package mysql

import (
	"errors"
	"log"
	"time"

	"github.com/alireza-hmd/c2/pkg/response"
	"github.com/alireza-hmd/c2/tasks"

	"gorm.io/gorm"
)

type TasksRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) tasks.Repository {
	return &TasksRepository{
		db: db,
	}
}

func (r *TasksRepository) Get(id int) (*tasks.Task, error) {
	var t tasks.Task
	res := r.db.Model(&tasks.Task{}).Where("id = ?", id).Limit(1).Find(&t)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting task from db\n")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return &t, nil
}

func (r *TasksRepository) List() ([]*tasks.Task, error) {
	var tt []*tasks.Task
	res := r.db.Model(&tasks.Task{}).Find(&tt)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting tasks from db")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return tt, nil
}

func (r *TasksRepository) ListClientTasks(client string) ([]*tasks.Task, error) {
	var tt []*tasks.Task
	res := r.db.Model(&tasks.Task{}).Where("client = ?", client).Find(&tt)
	if err := res.Error; err != nil {
		log.Println(err)
		return nil, errors.New("error getting tasks from db")
	}
	if res.RowsAffected == 0 {
		return nil, response.ErrNotFound
	}
	return tt, nil
}

func (r *TasksRepository) Create(t *tasks.Task) (int, error) {
	err := r.db.Create(t).Error
	if err != nil {
		log.Println(err)
		return 0, errors.New("error creating task\n")
	}
	return t.ID, nil
}

func (r *TasksRepository) Update(id int, t *tasks.TaskUpdate) error {
	t.UpdatedAt = time.Now()
	res := r.db.Model(&tasks.Task{}).Where("id = ?", id).Updates(t)
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error updating task\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *TasksRepository) Delete(id int) error {
	res := r.db.Where("id = ?", id).Delete(&tasks.Task{})
	if err := res.Error; err != nil {
		log.Println(err)
		return errors.New("error deleting task\n")
	}
	if res.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}
