package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type ProjectUsecase interface {
	GetAllProjects() ([]entity.Project, error)
	GetProjectByID(id uint) (*entity.Project, error)
	CreateProject(project *entity.Project) error
	UpdateProject(project *entity.Project, id uint) error
	DeleteProject(id uint) error

	GetAllProjectRoles() ([]entity.ProjectRole, error)
	GetProjectRoleByID(id uint) (*entity.ProjectRole, error)
	CreateProjectRole(projectRole *entity.ProjectRole) error
	UpdateProjectRole(projectRole *entity.ProjectRole, id uint) error
	DeleteProjectRole(id uint) error
}

type projectUsecase struct {
	projectRepository repository.ProjectRepository
}

func NewProjectUsecase(projectRepository repository.ProjectRepository) *projectUsecase {
	return &projectUsecase{projectRepository: projectRepository}
}

func (pu *projectUsecase) GetAllProjects() ([]entity.Project, error) {
	return pu.projectRepository.GetAllProjects()
}

func (pu *projectUsecase) GetProjectByID(id uint) (*entity.Project, error) {
	return pu.projectRepository.GetProjectByID(id)
}

func (pu *projectUsecase) CreateProject(project *entity.Project) error {
	return pu.projectRepository.CreateProject(project)
}

func (pu *projectUsecase) UpdateProject(project *entity.Project, id uint) error {
	return pu.projectRepository.UpdateProject(project, id)
}

func (pu *projectUsecase) DeleteProject(id uint) error {
	return pu.projectRepository.DeleteProject(id)
}

func (pu *projectUsecase) GetAllProjectRoles() ([]entity.ProjectRole, error) {
	return pu.projectRepository.GetAllProjectRoles()
}

func (pu *projectUsecase) GetProjectRoleByID(id uint) (*entity.ProjectRole, error) {
	return pu.projectRepository.GetProjectRoleByID(id)
}

func (pu *projectUsecase) CreateProjectRole(projectRole *entity.ProjectRole) error {
	return pu.projectRepository.CreateProjectRole(projectRole)
}

func (pu *projectUsecase) UpdateProjectRole(projectRole *entity.ProjectRole, id uint) error {
	return pu.projectRepository.UpdateProjectRole(projectRole, id)
}

func (pu *projectUsecase) DeleteProjectRole(id uint) error {
	return pu.projectRepository.DeleteProjectRole(id)
}
