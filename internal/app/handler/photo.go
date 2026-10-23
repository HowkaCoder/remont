package handler

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type PhotoHandler struct {
	usecase   usecase.PhotoUsecase
	dir       string
	publicURL string
}

func NewPhotoHandler(usecase usecase.PhotoUsecase, dir, publicURL string) *PhotoHandler {
	return &PhotoHandler{usecase: usecase, dir: dir, publicURL: strings.TrimSuffix(publicURL, "/")}
}

func (h *PhotoHandler) GetAllPhotos(c *fiber.Ctx) error {
	photos, err := h.usecase.GetAllPhotos()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(photos)
}

func (h *PhotoHandler) GetPhotoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid photo id"))
	}

	photo, err := h.usecase.GetPhotoByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "photo not found"})
	}
	return c.JSON(photo)
}

func (h *PhotoHandler) CreatePhoto(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return badRequest(c, err)
	}

	projectID, err := formUint(form, "projectID")
	if err != nil {
		return badRequest(c, err)
	}
	folderID, err := formUint(form, "folderID")
	if err != nil {
		return badRequest(c, err)
	}

	file := formFile(form)
	if file == nil {
		return badRequest(c, errors.New("file is required"))
	}
	if !imageExts[strings.ToLower(filepath.Ext(file.Filename))] {
		return badRequest(c, errors.New("only jpg, png, webp and heic images are allowed"))
	}

	name, err := saveUpload(c, file, h.dir)
	if err != nil {
		return serverError(c, err)
	}

	photo := entity.Photo{
		Title:         formValue(form, "title"),
		Filepath:      h.publicURL + "/uploads/photos/" + name,
		ProjectID:     projectID,
		PhotoFolderID: folderID,
	}
	if err := h.usecase.CreatePhoto(&photo); err != nil {
		removeUpload(h.dir, name)
		return serverError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(photo)
}

func (h *PhotoHandler) UpdatePhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid photo id"))
	}

	form, err := c.MultipartForm()
	if err != nil {
		return badRequest(c, err)
	}

	existing, err := h.usecase.GetPhotoByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "photo not found"})
	}

	photo := entity.Photo{Title: formValue(form, "title")}
	if v, err := formUint(form, "projectID"); err == nil {
		photo.ProjectID = v
	}
	if v, err := formUint(form, "folderID"); err == nil {
		photo.PhotoFolderID = v
	}

	if file := formFile(form); file != nil {
		if !imageExts[strings.ToLower(filepath.Ext(file.Filename))] {
			return badRequest(c, errors.New("only jpg, png, webp and heic images are allowed"))
		}
		name, err := saveUpload(c, file, h.dir)
		if err != nil {
			return serverError(c, err)
		}
		photo.Filepath = h.publicURL + "/uploads/photos/" + name
	}

	if err := h.usecase.UpdatePhoto(&photo, uint(id)); err != nil {
		return serverError(c, err)
	}

	if photo.Filepath != "" {
		removeUpload(h.dir, existing.Filepath)
	}

	return c.JSON(fiber.Map{"message": "successfully updated"})
}

func (h *PhotoHandler) DeletePhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid photo id"))
	}

	photo, err := h.usecase.GetPhotoByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "photo not found"})
	}

	if err := h.usecase.DeletePhoto(uint(id)); err != nil {
		return serverError(c, err)
	}
	if err := removeUpload(h.dir, photo.Filepath); err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"message": "successfully deleted"})
}

func (h *PhotoHandler) GetPhotosByFolderID(c *fiber.Ctx) error {
	folderID, err := strconv.Atoi(c.Params("folderID"))
	if err != nil {
		return badRequest(c, errors.New("invalid folder id"))
	}

	photos, err := h.usecase.GetPhotosByFolderID(uint(folderID))
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(photos)
}
