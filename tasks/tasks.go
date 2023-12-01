package tasks

import "time"

const (
	//task status
	ToDo      = 1
	Delivered = 2
	Done      = 3
	Failed    = 4

	//task type
	Timeout = "timeout"
	PWD     = "pwd"
)

type Task struct {
	ID        int       `json:"id"`
	Client    string    `json:"client"`
	Listener  string    `json:"listener"`
	Command   string    `json:"command"`
	Result    string    `json:"result"`
	TaskType  string    `json:"task_type"`
	Status    int       `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Task) Validate() error {
	return nil
}

func NewTask(client, listener, command string, taskType string) *Task {
	return &Task{
		Client:   client,
		Listener: listener,
		Command:  command,
		Status:   ToDo,
		TaskType: taskType,
	}
}

type TaskUpdate struct {
	Status    int       `json:"status"`
	Result    string    `json:"result"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Reader interface {
	Get(id int) (*Task, error)
	List() ([]*Task, error)
	ListToDoTasks(client string) ([]*Task, error)
	ListDoneTasks(client string) ([]*Task, error)
}

type Writer interface {
	Create(t *Task) (int, error)
	Update(id int, t *TaskUpdate) error
	Delete(id int) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get(id int) (*Task, error)
	List() ([]*Task, error)
	ListToDoTasks(client string) ([]*Task, error)
	ListDoneTasks(client string) ([]*Task, error)
	Create(client, listener, command, taskType string) (int, error)
	Update(id int, t *TaskUpdate) error
	Delete(id int) error
}
