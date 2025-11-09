package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BenjaminAGH/nocturnescope/backend/internal/infrastructure/security"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/interface/http/handlers"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/interface/http/middleware"
	"github.com/BenjaminAGH/nocturnescope/backend/internal/usecase/service"
)

func Register(app *fiber.App, userService *service.UserService, authService *service.AuthService, jwtService *security.JWTService) {
	api := app.Group("/api")

	// públicas
	RegisterAuthRoutes(api, authService, userService)

	// protegidas
	protected := api.Group("")
	protected.Use(middleware.JWTProtected(jwtService, authService))

	authHandler := handlers.NewAuthHandler(authService, userService)
	protected.Post("/auth/logout", authHandler.Logout)

	RegisterUserRoutes(protected, userService)

}
