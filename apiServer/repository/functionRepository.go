package repository

import (
	"faas/dataSource"
	"faas/models"
)

type FunctionRepository struct {
}

func NewFunctionRepository() *FunctionRepository {
	return &FunctionRepository{}
}

func (r *FunctionRepository) InsertFunction(function models.Function) (string, error) {
	return dataSource.InsertFunction(function)
}

func (r *FunctionRepository) DeleteFunction(id string) error {
	return dataSource.DeleteFunction(id)
}

func (r *FunctionRepository) GetFunctionByID(id string) (models.Function, error) {
	return dataSource.GetFunctionByID(id)
}

func (r *FunctionRepository) GetFunctionsByUsername(username string) ([]models.Function, error) {
	return dataSource.GetFunctionsByUsername(username)
}
