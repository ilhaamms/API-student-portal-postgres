package service

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/repository"
)

type UserService interface {
	Login(user model.User) error
	Register(user model.User) error

	CheckPassLength(pass string) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Login(user model.User) error {
	err := s.userRepository.CheckAvail(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Register(user model.User) error {
	err := s.userRepository.Add(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) CheckPassLength(pass string) bool {
	if len(pass) <= 5 {
		return true
	}

	return false
}
