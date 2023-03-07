package user

import (
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
}

type service struct {
	respository Repository
}

func NewService(respository Repository) *service {
	return &service{respository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	if input.Role != "" {
		user.Role = input.Role
	} else {
		user.Role = "user"
	}

	newUser, err := s.respository.Save(user)
	if err != nil {
		return user, err
	}
	return newUser, nil
}
