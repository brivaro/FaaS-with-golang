package repository

import (
	"faas/dataSource"
	"faas/models"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) InsertUser(user models.User) error {
	return dataSource.InsertUser(user)
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	return dataSource.GetUserByUsername(username)
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	return dataSource.GetAllUsers()
}
