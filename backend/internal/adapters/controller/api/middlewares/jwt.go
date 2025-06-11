package middlewares

import (
	"backend/cmd/app"
	"backend/internal/adapters/database/postgres"
	"backend/internal/domain/entity"
	"backend/internal/domain/service"
	"backend/internal/domain/utils/auth"
	"context"
	"github.com/gofiber/fiber/v3"
)

type UserService interface {
	GetByID(ctx context.Context, uuid string) (*entity.User, error)
}

type MiddlewareHandler struct {
	userService UserService
}

// NewMiddlewareHandler is a function that returns a new instance of MiddlewareHandler.
func NewMiddlewareHandler(app *app.App) *MiddlewareHandler {
	userStorage := postgres.NewUserStorage(app.DB)
	userService := service.NewUserService(userStorage)

	return &MiddlewareHandler{
		userService: userService,
	}
}

// IsAuthenticated is a function that checks whether the user has sufficient rights to access the endpoint
/*
 * tokenType string - the type of token that is required to access the endpoint
 * requiredRights ...string - the rights that the user must have
 */
func (h MiddlewareHandler) IsAuthenticated(tokenType string, requiredRights ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		user, fetchErr := auth.GetUserFromJWT(authHeader, tokenType, c.Context(), h.userService.GetByID)
		if fetchErr != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": fetchErr.Error(),
			})
		}

		/*
			if !config.RoleHasRights(user.Role, requiredRights) {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"success": false,
						"message": errorz.Forbidden,
					})
				}
		*/

		c.Locals("user", user)

		return c.Next()
	}
}
