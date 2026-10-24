package user

import (
	"errors"
	"fmt"
	"strings"
)

type Notifier interface {
	Send(to, msg string) error
}

type Service struct {
	notifier Notifier
	repo     Repository
}

func NewService(repo Repository, n Notifier) *Service {
	return &Service{repo: repo, notifier: n}
}

func (s *Service) Register(name, email string) (User, error) {
	if strings.TrimSpace(name) == "" {
		return User{}, errors.New("empty name")
	}
	user, err := s.repo.Save(User{Name: name, Email: email})
	if err != nil {
		return User{}, fmt.Errorf("register %s: %w", name, err)
	}
	err = s.notifier.Send(email, "Welcome, "+name+"!")
	if err != nil {
		return user, fmt.Errorf("notify %s: %w", name, err)
	}
	return user, nil
}

func (s *Service) Find(email string) (User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Users() ([]User, error) {
	return s.repo.List()
}
