package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type StateUsecase interface {
	CreateState(state *entity.State) error
	GetStateByID(id uint) (*entity.State, error)
	GetStatesByProjectID(projectID uint) ([]entity.State, error)
	UpdateState(state *entity.State) error
	DeleteState(id uint) error
	AssignWorkerToState(stateID, userID uint) error
	RemoveWorkerFromState(stateID, userID uint) error
}

type stateUsecase struct {
	repo repository.StateRepository
}

func NewStateUsecase(repo repository.StateRepository) *stateUsecase {
	return &stateUsecase{repo: repo}
}

func (su *stateUsecase) CreateState(state *entity.State) error {
	return su.repo.CreateState(state)
}

func (su *stateUsecase) GetStateByID(id uint) (*entity.State, error) {
	return su.repo.GetStateByID(id)
}

func (su *stateUsecase) GetStatesByProjectID(projectID uint) ([]entity.State, error) {
	return su.repo.GetStatesByProjectID(projectID)
}

func (su *stateUsecase) UpdateState(state *entity.State) error {
	return su.repo.UpdateState(state)
}

func (su *stateUsecase) DeleteState(id uint) error {
	return su.repo.DeleteState(id)

}

func (su *stateUsecase) AssignWorkerToState(stateID, userID uint) error {
	return su.repo.AssignWorkerToState(stateID, userID)
}

func (su *stateUsecase) RemoveWorkerFromState(stateID, userID uint) error {
	return su.repo.RemoveWorkerFromState(stateID, userID)
}
