package handler

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var imageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".heic": true}

func formValue(form *multipart.Form, key string) string {
	if v := form.Value[key]; len(v) > 0 {
		return strings.TrimSpace(v[0])
	}
	return ""
}

func formUint(form *multipart.Form, key string) (uint, error) {
	n, err := strconv.ParseUint(formValue(form, key), 10, 64)
	if err != nil {
		return 0, errors.New(key + " is required")
	}
	return uint(n), nil
}

func formFile(form *multipart.Form) *multipart.FileHeader {
	if files := form.File["file"]; len(files) > 0 {
		return files[0]
	}
	return nil
}

func saveUpload(c *fiber.Ctx, file *multipart.FileHeader, dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := uuid.New().String() + strings.ToLower(filepath.Ext(file.Filename))
	if err := c.SaveFile(file, filepath.Join(dir, name)); err != nil {
		return "", err
	}
	return name, nil
}

func removeUpload(dir, stored string) error {
	if stored == "" {
		return nil
	}
	err := os.Remove(filepath.Join(dir, filepath.Base(stored)))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func badRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

func serverError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}
