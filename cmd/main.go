package main

import (
	"log"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	db = internal.Init()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	docRepo := repository.NewDocumentRepository(db)
	docUsecase := usecase.NewDocumentUsecase(docRepo)
	docHandler := handler.NewDocumentHandler(docUsecase)
	app := fiber.New()

	app.Get("/api/docs", docHandler.GetAllDocuments)
	app.Post("/api/docs", docHandler.CreateDocument)
	app.Patch("/api/docs/:id", docHandler.UpdateDocument)

	app.Post("/roles", createRole)
	app.Post("/permissions", createPermission)
	app.Post("/roles/assign", assignPermissionToRole)
	app.Post("/users/assign-role", assignRoleToUser)
	app.Post("/login", login)
	app.Post("/users", userHandler.CreateUser)

	app.Use(authMiddleware)
	app.Get("/users", PermissionMiddleware("get_all_users"), userHandler.GetAllUsers)

	log.Fatal(app.Listen(":3000"))
}

func PermissionMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint) // Assuming userID is set in context after authentication
		var user entity.User
		db.Preload("Roles.Permissions").Preload("Permissions").First(&user, userID)

		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				if permission.Name == requiredPermission {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}

func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userID"].(float64))
	c.Locals("userID", userID)

	return c.Next()
}

var jwtSecret = []byte("your_secret_key")

func login(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user entity.User
	db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func createRole(c *fiber.Ctx) error {
	var role entity.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	db.Create(&role)
	return c.Status(fiber.StatusCreated).JSON(role)
}

func createPermission(c *fiber.Ctx) error {
	permission := new(entity.Permission)
	if err := c.BodyParser(permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	db.Create(&permission)
	return c.Status(fiber.StatusCreated).JSON(permission)
}

func assignPermissionToRole(c *fiber.Ctx) error {
	var input struct {
		RoleID       uint `json:"RoleID"`
		PermissionID uint `json:"PermissionID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	rolePermission := entity.RolePermission{
		RoleID:       input.RoleID,
		PermissionID: input.PermissionID,
	}

	db.Create(&rolePermission)
	return c.Status(fiber.StatusCreated).JSON(rolePermission)
}

func assignRoleToUser(c *fiber.Ctx) error {
	var input struct {
		UserID uint `json:"UserID"`
		RoleID uint `json:"RoleID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userRole := entity.UserRole{
		UserID: input.UserID,
		RoleID: input.RoleID,
	}

	db.Create(&userRole)
	return c.Status(fiber.StatusCreated).JSON(userRole)
}
