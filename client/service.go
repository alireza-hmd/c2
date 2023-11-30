package client

type Service struct {
	repo Repository
}

func NewService(r Repository) UseCase {
	return &Service{
		repo: r,
	}
}

func (s *Service) Get(token string) (*Client, error) {
	return s.repo.Get(token)
}

func (s *Service) List() ([]Client, error) {
	return s.repo.List()
}

func (s *Service) Create(c *Client) (int, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(c)
}

func (s *Service) Update(token string, c *ClientUpdate) error {
	return s.repo.Update(token, c)
}

func (s *Service) Delete(token string) error {
	return s.repo.Delete(token)
}
