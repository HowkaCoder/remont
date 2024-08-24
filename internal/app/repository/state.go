package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type StateRepository interface {
	CreateState(state *entity.State) error
	GetStateByID(id uint) (*entity.State, error)
	GetStatesByProjectID(projectID uint) ([]entity.State, error)
	UpdateState(state *entity.State) error
	DeleteState(id uint) error
	AssignWorkerToState(stateID, userID uint) error
	RemoveWorkerFromState(stateID, userID uint) error
}

type stateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *stateRepository {
	return &stateRepository{db: db}
}

func (sr *stateRepository) CreateState(state *entity.State) error {
	return sr.db.Create(state).Error
}

func (sr *stateRepository) GetStateByID(id uint) (*entity.State, error) {
	var state entity.State
	if err := sr.db.First(&state, id).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (sr *stateRepository) GetStatesByProjectID(projectID uint) ([]entity.State, error) {
	var states []entity.State
	if err := sr.db.Where("project_id =?", projectID).Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (sr *stateRepository) UpdateState(state *entity.State) error {
	return sr.db.Save(state).Error
}

func (sr *stateRepository) DeleteState(id uint) error {
	return sr.db.Delete(&entity.State{}, id).Error
}

func (sr *stateRepository) AssignWorkerToState(stateID, userID uint) error {
	stateUser := entity.StateUser{
		StateID: stateID,
		UserID:  userID,
	}
	return sr.db.Create(&stateUser).Error
}

func (sr *stateRepository) RemoveWorkerFromState(stateID, userID uint) error {
	return sr.db.Where("state_id = ? AND user_id = ?", stateID, userID).Delete(&entity.StateUser{}).Error

}
