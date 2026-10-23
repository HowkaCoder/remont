package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const tokenTTL = 24 * time.Hour

type AuthHandler struct {
	repo   repository.AuthRepository
	secret []byte
}

func NewAuthHandler(repo repository.AuthRepository, secret []byte) *AuthHandler {
	return &AuthHandler{repo: repo, secret: secret}
}

type registerRequest struct {
	FirstName   string
	LastName    string
	MiddleName  string
	Email       string
	Password    string
	PhoneNumber string
	City        string
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, err)
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !strings.Contains(req.Email, "@") || len(req.Password) < 8 {
		return badRequest(c, errors.New("valid email and password of at least 8 characters are required"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return serverError(c, err)
	}

	user := entity.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		MiddleName:  req.MiddleName,
		Email:       req.Email,
		Password:    string(hash),
		PhoneNumber: req.PhoneNumber,
		City:        req.City,
	}
	if err := h.repo.CreateUserWithRole(&user, "client"); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user with this email already exists"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created successfully", "id": user.ID})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, err)
	}

	user, err := h.repo.FindUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	projects, err := h.repo.ProjectIDs(user.ID)
	if err != nil {
		return serverError(c, err)
	}

	role := ""
	if len(user.Roles) > 0 {
		role = user.Roles[0].Name
	}

	claims := entity.JwtClaims{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      role,
		Projects:  projects,
		City:      user.City,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.secret)
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"token": token, "claims": claims})
}

func (h *AuthHandler) Authenticate(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or malformed token"})
	}

	token, err := jwt.ParseWithClaims(strings.TrimPrefix(header, "Bearer "), &entity.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return h.secret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	claims, ok := token.Claims.(*entity.JwtClaims)
	if !ok || claims.UserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	c.Locals("userID", claims.UserID)
	c.Locals("claims", claims)
	return c.Next()
}

func (h *AuthHandler) Require(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		user, err := h.repo.FindUserWithRoles(userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}
		if err != nil {
			return serverError(c, err)
		}

		for _, role := range user.Roles {
			for _, p := range role.Permissions {
				if p.Name == permission {
					return c.Next()
				}
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*entity.JwtClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{
		"user_id":    claims.UserID,
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
		"role":       claims.Role,
		"city":       claims.City,
	})
}

func (h *AuthHandler) ProjectMembers(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid project id"))
	}

	workers, manager, err := h.repo.ProjectMembers(uint(id))
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(fiber.Map{"Workers": workers, "Manager": manager})
}

func (h *AuthHandler) CreateRole(c *fiber.Ctx) error {
	var role entity.Role
	if err := c.BodyParser(&role); err != nil {
		return badRequest(c, err)
	}
	if role.Name == "" {
		return badRequest(c, errors.New("name is required"))
	}
	if err := h.repo.CreateRole(&role); err != nil {
		return serverError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(role)
}

func (h *AuthHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.repo.GetRoles()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(roles)
}

func (h *AuthHandler) CreatePermission(c *fiber.Ctx) error {
	var permission entity.Permission
	if err := c.BodyParser(&permission); err != nil {
		return badRequest(c, err)
	}
	if permission.Name == "" {
		return badRequest(c, errors.New("name is required"))
	}
	if err := h.repo.CreatePermission(&permission); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(permission)
}

func (h *AuthHandler) GetPermissions(c *fiber.Ctx) error {
	permissions, err := h.repo.GetPermissions()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(permissions)
}

func (h *AuthHandler) AssignPermissionToRole(c *fiber.Ctx) error {
	var input struct {
		RoleID       uint `json:"RoleID"`
		PermissionID uint `json:"PermissionID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, err)
	}
	if err := h.repo.AssignPermissionToRole(input.RoleID, input.PermissionID); err != nil {
		return notFoundOr500(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(input)
}

func (h *AuthHandler) AssignRoleToUser(c *fiber.Ctx) error {
	var input struct {
		UserID uint `json:"UserID"`
		RoleID uint `json:"RoleID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, err)
	}
	if err := h.repo.AssignRoleToUser(input.UserID, input.RoleID); err != nil {
		return notFoundOr500(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(input)
}

func notFoundOr500(c *fiber.Ctx, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return serverError(c, err)
}
