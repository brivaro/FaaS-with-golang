package functions

import (
	"context"
	"faas/models"
	"faas/repository"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerClient interface {
	ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
	Close() error
}

type FunctionService struct {
	repo   repository.FunctionRepository
	docker DockerClient
}

func NewService(repo repository.FunctionRepository) (*FunctionService, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &FunctionService{
		repo:   repo,
		docker: docker,
	}, nil
}

func (s *FunctionService) RegisterFunction(ctx context.Context, user models.User, req RegisterRequest) (string, error) {
	function := models.Function{
		ID:        "NotID",
		Name:      req.Name,
		Data:      req.Data,
		CreatedAt: "",
		UserID:    user.Username,
	}

	reader, err := s.docker.ImagePull(ctx, req.Data, image.PullOptions{})
	if err != nil {
		return "", ErrFunctionNotPulled 
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	return s.repo.InsertFunction(function)
}

func (s *FunctionService) DeleteFunction(ctx context.Context, user models.User, functionID string) error {
	function, err := s.repo.GetFunctionByID(functionID)
	if err != nil {
		return ErrFunctionNotFound
	}

	if function.UserID != user.Username {
		return ErrUnauthorized
	}

	return s.repo.DeleteFunction(functionID)
}

func (s *FunctionService) GetUserFunctions(ctx context.Context, username string) ([]models.Function, error) {
	return s.repo.GetFunctionsByUsername(username)
}
