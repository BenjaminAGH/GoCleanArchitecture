package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/BenjaminAGH/nocturnescope/backend/internal/infrastructure/database"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/infrastructure/repository"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/infrastructure/security"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/infrastructure/session"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/interface/http/routes"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/usecase/service"
)

func main() {
	_ = godotenv.Load("config/.env")

	db := database.Connect()

	userRepo := repository.NewUserGormRepository(db)
	userService := service.NewUserService(userRepo)

	jwtService := security.NewJWTServiceFromEnv()
	sessionStore := session.NewMemoryStore()
	authService := service.NewAuthService(userRepo, jwtService, sessionStore)

	app := fiber.New()

	routes.Register(app, userService, authService, jwtService)

	log.Fatal(app.Listen(":3000"))
}
