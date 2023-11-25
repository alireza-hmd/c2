package listener

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) UseCase {
	return &Service{
		repo: r,
	}
}

func (s *Service) Get(name string) (*Listener, error) {
	return s.repo.Get(name)
}

func (s *Service) List() ([]Listener, error) {
	return s.repo.ListActive()
}

func (s *Service) Create(l *Listener, stop map[int](chan Cancel)) (int, error) {
	if err := l.Validate(); err != nil {
		return 0, err
	}
	if l.Active == ActiveStatus {
		stop[l.ID] = make(chan Cancel, 1)
		go s.Run(l, stop[l.ID])
	}
	return s.repo.Create(l)
}

func (s *Service) Update(name string, l *ListenerUptade) error {
	return s.repo.Update(name, l)
}

func (s *Service) Delete(name string, stop chan Cancel) error {
	return s.repo.Delete(name, stop)
}

func (s *Service) Activation(l *Listener, stop chan Cancel, status int) error {
	if status == ActiveStatus {
		if l.Active == ActiveStatus {
			return errors.New("listener is already active\n")
		}
		if err := changeActivationStatus(s.repo, l, ActiveStatus); err != nil {
			return err
		}
		go s.Run(l, stop)
		return nil
	}

	if l.Active == DeactiveStatus {
		return errors.New("listener is already deactivated\n")
	}
	if err := changeActivationStatus(s.repo, l, DeactiveStatus); err != nil {
		return err
	}
	stop <- Cancel{}
	return nil
}

func (s *Service) Run(l *Listener, stop chan Cancel) {
	for {
		select {
		case <-stop:
			fmt.Printf("%s listener deactivated\n", l.Name)
			close(stop)
			return
		default:
			fmt.Printf("%s listener is active\n", l.Name)
			time.Sleep(5 * time.Second)
		}
	}
}

func (s *Service) RunActiveListeners(stop map[int](chan Cancel)) error {
	ll, err := s.repo.ListActive()
	if err != nil {
		return fmt.Errorf("error running active listeners | %v", err)
	}
	for _, l := range ll {
		stop[l.ID] = make(chan Cancel, 1)
		go s.Run(&l, stop[l.ID])
	}
	if len(ll) > 0 {
		log.Println("listeners Activated")
	}
	return nil
}

func changeActivationStatus(r Repository, l *Listener, status int) error {
	listener := &ListenerUptade{
		Active: status,
	}
	if err := r.Update(l.Name, listener); err != nil {
		return fmt.Errorf("error changing listener active status | %v", err)
	}
	return nil
}
