package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type ProjectRepository interface {
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

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *projectRepository {
	return &projectRepository{db: db}
}

func (pr *projectRepository) GetAllProjects() ([]entity.Project, error) {
	var projects []entity.Project
	if err := pr.db.Preload("Members.User").Preload("Members.Role").Preload("States").Preload("DocumentFolders.Documents").Preload("PhotoFolders").Preload("Chars").Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (pr *projectRepository) GetProjectByID(id uint) (*entity.Project, error) {
	var project *entity.Project
	if err := pr.db.Preload("Members.User").Preload("Members.Role").Preload("States").Preload("DocumentFolders").Preload("PhotoFolders").Preload("Chars").First(&project, id).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (pr *projectRepository) CreateProject(project *entity.Project) error {
	return pr.db.Create(&project).Error
}

func (pr *projectRepository) UpdateProject(project *entity.Project, id uint) error {
	var eProject *entity.Project
	if err := pr.db.First(&eProject, id).Error; err != nil {
		return err
	}

	if project.Title != "" {
		eProject.Title = project.Title
	}

	return pr.db.Save(eProject).Error
}

func (pr *projectRepository) DeleteProject(id uint) error {
	var project *entity.Project
	if err := pr.db.First(&project, id).Error; err != nil {
		return err
	}
	return pr.db.Delete(&project).Error
}

func (pr *projectRepository) GetAllProjectRoles() ([]entity.ProjectRole, error) {
	var projectRoles []entity.ProjectRole
	if err := pr.db.Preload("Project").Preload("User").Preload("Role").Find(&projectRoles).Error; err != nil {
		return nil, err
	}
	return projectRoles, nil
}

func (pr *projectRepository) GetProjectRoleByID(id uint) (*entity.ProjectRole, error) {
	var projectRole *entity.ProjectRole
	if err := pr.db.Preload("Project").Preload("User").Preload("Role").First(&projectRole, id).Error; err != nil {
		return nil, err
	}
	return projectRole, nil
}

func (pr *projectRepository) CreateProjectRole(projectRole *entity.ProjectRole) error {
	return pr.db.Create(&projectRole).Error
}

func (pr *projectRepository) UpdateProjectRole(projectRole *entity.ProjectRole, id uint) error {
	var eProjectRole *entity.ProjectRole
	if err := pr.db.First(&eProjectRole, id).Error; err != nil {
		return err
	}

	if projectRole.ProjectID != 0 {
		eProjectRole.ProjectID = projectRole.ProjectID
	}
	if projectRole.UserID != 0 {
		eProjectRole.UserID = projectRole.UserID
	}
	if projectRole.RoleID != 0 {
		eProjectRole.RoleID = projectRole.RoleID
	}

	return pr.db.Save(eProjectRole).Error
}

func (pr *projectRepository) DeleteProjectRole(id uint) error {
	var projectRole *entity.ProjectRole
	if err := pr.db.First(&projectRole, id).Error; err != nil {
		return err
	}
	return pr.db.Delete(&projectRole).Error
}
