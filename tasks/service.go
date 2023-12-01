package tasks

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(id int) (*Task, error) {
	return s.repo.Get(id)
}

func (s *Service) List() ([]*Task, error) {
	return s.repo.List()
}
func (s *Service) ListToDoTasks(client string) ([]*Task, error) {
	return s.repo.ListToDoTasks(client)
}
func (s *Service) ListDoneTasks(client string) ([]*Task, error) {
	return s.repo.ListDoneTasks(client)
}

func (s *Service) Create(client, listener, command, taskType string) (int, error) {
	t := NewTask(client, listener, command, taskType)
	if err := t.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(t)
}

func (s *Service) Update(id int, t *TaskUpdate) error {
	return s.repo.Update(id, t)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
