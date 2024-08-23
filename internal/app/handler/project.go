package handler

import (
	"strconv"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	usecase usecase.ProjectUsecase
}

func NewProjectHandler(usecase usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{usecase}
}

func (h *ProjectHandler) GetAllProjects(c *fiber.Ctx) error {
	projects, err := h.usecase.GetAllProjects()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(projects)
}

func (h *ProjectHandler) GetProjectByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid project ID")
	}

	project, err := h.usecase.GetProjectByID(uint(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(project)
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var project entity.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).SendString("Invalid request payload")
	}

	if err := h.usecase.CreateProject(&project); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(project)
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid project ID")
	}

	var project *entity.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(400).SendString("Invalid request payload")
	}

	project.ID = uint(id)

	if err := h.usecase.UpdateProject(project, uint(id)); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(project)
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid project ID")
	}

	if err := h.usecase.DeleteProject(uint(id)); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendString("Project deleted successfully")
}

func (ph *ProjectHandler) GetAllProjectRole(c *fiber.Ctx) error {

	projectRole, err := ph.usecase.GetAllProjectRoles()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(projectRole)
}

func (ph *ProjectHandler) GetProjectRoleByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid project role ID")
	}

	projectRole, err := ph.usecase.GetProjectRoleByID(uint(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(projectRole)
}

func (ph *ProjectHandler) CreateProjectRole(c *fiber.Ctx) error {
	var projectRole entity.ProjectRole
	if err := c.BodyParser(&projectRole); err != nil {
		return c.Status(400).SendString("Invalid request payload")
	}

	if err := ph.usecase.CreateProjectRole(&projectRole); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(projectRole)
}

func (ph *ProjectHandler) UpdateProjectRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid project role ID")
	}

	var projectRole *entity.ProjectRole
	if err := c.BodyParser(&projectRole); err != nil {
		return c.Status(400).SendString("Invalid request payload")
	}

	projectRole.ID = uint(id)

	if err := ph.usecase.UpdateProjectRole(projectRole, uint(id)); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(projectRole)
}

func (ph *ProjectHandler) DeleteProjectRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid project role ID")
	}

	if err := ph.usecase.DeleteProjectRole(uint(id)); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendString("Project role deleted successfully")
}
