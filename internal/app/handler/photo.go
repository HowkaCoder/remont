package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type PhotoHandler struct {
	usecase usecase.PhotoUsecase
}

func NewPhotoHandler(usecase usecase.PhotoUsecase) *PhotoHandler {
	return &PhotoHandler{usecase}
}

func (h *PhotoHandler) GetAllPhotos(c *fiber.Ctx) error {
	photos, err := h.usecase.GetAllPhotos()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(photos)
}

func (h *PhotoHandler) GetPhotoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid photo ID")
	}

	photo, err := h.usecase.GetPhotoByID(uint(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(photo)
}

func (h *PhotoHandler) CreatePhoto(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	projectID, err := strconv.Atoi(form.Value["projectID"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project ID is required"})
	}

	folderID, err := strconv.Atoi(form.Value["folderID"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project ID is required"})
	}

	title := form.Value["title"][0]

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	file := files[0]

	uploadDir := "./uploads/photos"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	backURL := "https://remont-production.up.railway.app/"
	filePath := filepath.Join(uploadDir, fileName)
	Fpath := filepath.Join(backURL, filePath)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	photo := entity.Photo{
		Title:         title,
		Filepath:      Fpath,
		ProjectID:     uint(projectID),
		PhotoFolderID: uint(folderID),
	}

	if err := h.usecase.CreatePhoto(&photo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "successfully created"})

}

func (h *PhotoHandler) UpdatePhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	projectID, err := strconv.Atoi(form.Value["projectID"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project ID is required"})
	}

	folderID, err := strconv.Atoi(form.Value["folderID"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project ID is required"})
	}

	title := form.Value["title"][0]

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	ePhoto, err := h.usecase.GetPhotoByID(uint(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if len(files) > 0 {
		if err := os.Remove(ePhoto.Filepath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		file := files[0]

		uploadDir := "./uploads/photos"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := filepath.Join(uploadDir, fileName)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		photo := entity.Photo{
			Title:         title,
			Filepath:      filePath,
			ProjectID:     uint(projectID),
			PhotoFolderID: uint(folderID),
		}

		if err := h.usecase.UpdatePhoto(&photo, uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		photo := entity.Photo{
			Title:         title,
			ProjectID:     uint(projectID),
			PhotoFolderID: uint(folderID),
		}

		if err := h.usecase.UpdatePhoto(&photo, uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"message": "Successfully updated"})

}

func (h *PhotoHandler) DeletePhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	photo, err := h.usecase.GetPhotoByID(uint(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err := os.Remove(photo.Filepath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.usecase.DeletePhoto(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Successfully deleted"})
}

func (h *PhotoHandler) GetPhotosByFolderID(c *fiber.Ctx) error {
	folderID, err := strconv.Atoi(c.Params("folderID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	photos, err := h.usecase.GetPhotosByFolderID(uint(folderID))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(photos)
}
