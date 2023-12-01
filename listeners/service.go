package listeners

import (
	"errors"
	"fmt"
	"log"

	"github.com/alireza-hmd/c2/clients"
	"github.com/alireza-hmd/c2/tasks"
)

type Service struct {
	lRepo Repository
	cRepo clients.Repository
	tRepo tasks.Repository
}

func NewService(lr Repository, cr clients.Repository, tr tasks.Repository) UseCase {
	return &Service{
		lRepo: lr,
		cRepo: cr,
		tRepo: tr,
	}
}

func (s *Service) Get(name string) (*Listener, error) {
	return s.lRepo.Get(name)
}

func (s *Service) List() ([]*Listener, error) {
	return s.lRepo.ListActive()
}

func (s *Service) Create(l *Listener, stop map[int](chan Cancel)) (int, error) {
	if err := l.Validate(); err != nil {
		return 0, err
	}
	if l.Active == ActiveStatus {
		stop[l.ID] = make(chan Cancel, 1)
		go s.Run(l, stop[l.ID])
	}
	return s.lRepo.Create(l)
}

func (s *Service) Update(name string, l *ListenerUptade) error {
	return s.lRepo.Update(name, l)
}

func (s *Service) Delete(name string, stop chan Cancel) error {
	return s.lRepo.Delete(name, stop)
}

func (s *Service) Activation(l *Listener, stop chan Cancel, status int) error {
	if status == ActiveStatus {
		if l.Active == ActiveStatus {
			return errors.New("listener is already active\n")
		}
		if err := changeActivationStatus(s.lRepo, l, ActiveStatus); err != nil {
			return err
		}
		go s.Run(l, stop)
		return nil
	}

	if l.Active == DeactiveStatus {
		return errors.New("listener is already deactivated\n")
	}
	if err := changeActivationStatus(s.lRepo, l, DeactiveStatus); err != nil {
		return err
	}
	stop <- Cancel{}
	return nil
}

func (s *Service) Run(l *Listener, stop chan Cancel) {
	cService := clients.NewService(s.cRepo)
	tService := tasks.NewService(s.tRepo)
	InitHandler(s, cService, tService, l.Name, l.Port, stop)
}

func (s *Service) RunActiveListeners(stop map[int](chan Cancel)) error {
	ll, err := s.lRepo.ListActive()
	if err != nil {
		return fmt.Errorf("error running active listeners | %v", err)
	}
	for _, l := range ll {
		stop[l.ID] = make(chan Cancel, 1)
		go s.Run(l, stop[l.ID])
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
	if status == DeactiveStatus {
		listener.Connected = Disconnected
	}
	if err := r.Update(l.Name, listener); err != nil {
		return fmt.Errorf("error changing listener active status | %v", err)
	}
	return nil
}
