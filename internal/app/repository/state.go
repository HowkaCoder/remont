package repository

import (
	"fmt"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type StateRepository interface {
	CreateState(state *entity.State) error
	GetStateByID(id uint) (*entity.State, error)
	GetStatesByProjectID(projectID uint) ([]entity.State, error)
	UpdateState(state *entity.State) error
	DeleteState(id uint) error

	GetStatesByWorkerID(id uint) ([]entity.State, error)

	AssignWorkerToState(stateID, userID uint) error
	RemoveWorkerFromState(stateID, userID uint) error

	CreateRepairDetails(details *entity.RepairDetails) error
	GetRepairDetailsByProjectID(projectID uint) (*entity.RepairDetails, error)
	UpdateRepairDetails(details *entity.RepairDetails) error
	DeleteRepairDetails(id uint) error
}

type stateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *stateRepository {
	return &stateRepository{db: db}
}

func (sr *stateRepository) GetStatesByWorkerID(id uint) ([]entity.State, error) {
	var StateUsers []entity.StateUser
	if err := sr.db.Where("user_id = ?", id).Find(&StateUsers).Error; err != nil {
		return nil, err
	}

	var states []entity.State
	for _, stateUser := range StateUsers {
		var state entity.State
		if err := sr.db.First(&state, stateUser.StateID).Error; err != nil {
			return nil, err
		}
		states = append(states, state)
	}

	return states, nil
}

func (sr *stateRepository) CreateState(state *entity.State) error {

	sr.db.Create(state)

	var states []entity.State
	if err := sr.db.Where("project_id = ?", state.ProjectID).Find(&states).Error; err != nil {
		return err
	}

	var totalNecessary, totalPaid uint
	for _, state := range states {
		totalNecessary += state.NecessaryMoney
		totalPaid += state.PaidMoney
	}
	var details entity.RepairDetails
	if err := sr.db.Where("project_id = ?", state.ProjectID).First(&details).Error; err != nil {
		return err
	}

	details.NecessaryWork = totalNecessary
	details.PaidWork = totalPaid

	return sr.db.Save(&details).Error

}

func (sr *stateRepository) GetStateByID(id uint) (*entity.State, error) {
	var state entity.State
	if err := sr.db.Preload("Workers").First(&state, id).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (sr *stateRepository) GetStatesByProjectID(projectID uint) ([]entity.State, error) {
	var states []entity.State
	if err := sr.db.Preload("Workers").Where("project_id =?", projectID).Find(&states).Error; err != nil {
		return nil, err
	}
	//	sr.db.Model(&states).Association("Users").Clear() // Очистка всех связей с пользователями
	//	sr.db.Model(&states).Association("Users").Append(newUsers...) // Если нужно обновить список пользователей

	return states, nil
}

func (sr *stateRepository) UpdateState(state *entity.State) error {
	sr.db.Save(state)

	var states []entity.State
	if err := sr.db.Where("project_id = ?", state.ProjectID).Find(&states).Error; err != nil {
		return err
	}

	var totalNecessary, totalPaid uint
	for _, state := range states {
		totalNecessary += state.NecessaryMoney
		totalPaid += state.PaidMoney
	}
	var details entity.RepairDetails
	if err := sr.db.Where("project_id = ?", state.ProjectID).First(&details).Error; err != nil {
		return err
	}

	details.NecessaryWork = totalNecessary
	details.PaidWork = totalPaid

	return sr.db.Save(&details).Error

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

	//return sr.db.Where("state_id = ? AND user_id = ?", stateID, userID).Delete(&entity.StateUser{}).Error

	if err := sr.db.Where("user_id = ? AND state_id = ?", userID, stateID).Delete(&entity.StateUser{}).Error; err != nil {
		return err
	}

	var state entity.State
	if err := sr.db.First(&state, stateID).Error; err != nil {
		return err
	}

	sr.db.Model(&state).Association("Workers").Clear()

	return nil
}

func (rdr *stateRepository) CreateRepairDetails(details *entity.RepairDetails) error {
	var states []entity.State
	if err := rdr.db.Where("project_id = ?", details.ProjectID).Find(&states).Error; err != nil {
		return err
	}

	var totalNecessary, totalPaid uint
	for _, state := range states {
		totalNecessary += state.NecessaryMoney
		totalPaid += state.PaidMoney
	}

	details.NecessaryWork = totalNecessary
	details.PaidWork = totalPaid

	fmt.Println(details)
	return rdr.db.Create(details).Error
}

func (rdr *stateRepository) GetRepairDetailsByProjectID(projectID uint) (*entity.RepairDetails, error) {
	var details entity.RepairDetails
	if err := rdr.db.Where("project_id = ?", projectID).First(&details).Error; err != nil {
		return nil, err
	}
	return &details, nil
}

func (rdr *stateRepository) UpdateRepairDetails(details *entity.RepairDetails) error {
	return rdr.db.Save(details).Error
}

func (rdr *stateRepository) DeleteRepairDetails(id uint) error {
	return rdr.db.Delete(&entity.RepairDetails{}, id).Error
}
