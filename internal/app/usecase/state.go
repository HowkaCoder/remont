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

	GetStatesByWorkerID(id uint) ([]entity.State, error)
	DeleteState(id uint) error
	AssignWorkerToState(stateID, userID uint) error
	RemoveWorkerFromState(stateID, userID uint) error

	CreateRepairDetails(details *entity.RepairDetails) error
	GetRepairDetailsByProjectID(projectID uint) (*entity.RepairDetails, error)
	UpdateRepairDetails(details *entity.RepairDetails) error
	DeleteRepairDetails(id uint) error
}

type stateUsecase struct {
	repo repository.StateRepository
}

func NewStateUsecase(repo repository.StateRepository) *stateUsecase {
	return &stateUsecase{repo: repo}
}

func (su *stateUsecase) GetStatesByWorkerID(id uint) ([]entity.State, error) {
	return su.repo.GetStatesByWorkerID(id)
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

func (su *stateUsecase) CreateRepairDetails(details *entity.RepairDetails) error {
	return su.repo.CreateRepairDetails(details)
}

func (su *stateUsecase) GetRepairDetailsByProjectID(projectID uint) (*entity.RepairDetails, error) {
	return su.repo.GetRepairDetailsByProjectID(projectID)
}

func (su *stateUsecase) UpdateRepairDetails(details *entity.RepairDetails) error {
	return su.repo.UpdateRepairDetails(details)
}

func (su *stateUsecase) DeleteRepairDetails(id uint) error {
	return su.repo.DeleteRepairDetails(id)
}
