package handler

import (
	"strconv"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type StateHandler struct {
	usecase usecase.StateUsecase
}

func NewStateHandler(usecase usecase.StateUsecase) *StateHandler {
	return &StateHandler{usecase: usecase}
}

func (h *StateHandler) GetStatesByWorkerID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	states, err := h.usecase.GetStatesByWorkerID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.JSON(states)
}

func (h *StateHandler) GetStatesByProjectID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project ID",
			"error":   err.Error(),
		})
	}

	states, err := h.usecase.GetStatesByProjectID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting states by project ID",
			"error":   err.Error(),
		})
	}

	return c.JSON(states)
}

func (h *StateHandler) CreateState(c *fiber.Ctx) error {
	var state entity.State

	if err := c.BodyParser(&state); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
			"error":   err.Error(),
		})
	}

	err := h.usecase.CreateState(&state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating state",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "State created successfully",
	})
}

func (h *StateHandler) UpdateState(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid state ID",
			"error":   err.Error(),
		})
	}

	var state entity.State

	if err := c.BodyParser(&state); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
			"error":   err.Error(),
		})
	}

	state.ID = uint(id)

	err = h.usecase.UpdateState(&state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating state",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "State updated successfully",
	})
}

func (h *StateHandler) DeleteState(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid state ID",
			"error":   err.Error(),
		})
	}

	err = h.usecase.DeleteState(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting state",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "State deleted successfully",
	})
}

func (h *StateHandler) AssignWorkerToState(c *fiber.Ctx) error {
	/*	userId, err := strconv.Atoi(c.Query("userID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid user ID",
				"error":   err.Error(),
			})
		}

		stateID, err := strconv.Atoi(c.Query("stateID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid state ID",
				"error":   err.Error(),
			})
		}
	**/
	var request struct {
		UserID  uint `json:"userID"`
		StateID uint `json:"stateID"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.usecase.AssignWorkerToState(uint(request.StateID), uint(request.UserID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error assigning worker to state",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user successfully assigned",
	})

}

func (h *StateHandler) RemoveWorkerFromState(c *fiber.Ctx) error {
	/**
		userId, err := strconv.Atoi(c.Query("userID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid user ID",
				"error":   err.Error(),
			})
		}

		stateID, err := strconv.Atoi(c.Query("stateID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid state ID",
				"error":   err.Error(),
			})
		}
	**/

	var request struct {
		UserID  uint `json:"userID"`
		StateID uint `json:"stateID"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	if err := h.usecase.RemoveWorkerFromState(uint(request.StateID), uint(request.UserID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error assigning worker to state",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user successfully removed",
	})
}
