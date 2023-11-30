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

func (s *Service) Create(t *Task) (int, error) {
	return s.repo.Create(t)
}

func (s *Service) Update(id int, t *TaskUpdate) error {
	return s.repo.Update(id, t)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
