package repository


import (
  "github.com/HowkaCoder/remont/internal/app/entity"
  "gorm.io/gorm"
)


type ActRepository interface {
  GetAllActs() ([]entity.Act , error)
  GetAllActsByProjectID(id uint) ([]entity.Act , error) 
  GetActByID(id uint) (*entity.Act , error)
  CreateAct(act *entity.Act) error 
  UpdateAct(act *entity.Act , id uint) error 
  DeleteAct
}
