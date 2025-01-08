package service

import "goGin/internal/static/repository"

func GetUserByID(id string) (repository.User, error) {
	return repository.FindUserByID(id)
}
