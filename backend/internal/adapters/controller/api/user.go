package api

import (
	"backend/cmd/app"
	"backend/internal/adapters/controller/api/validator"
	"backend/internal/adapters/repository"
	"backend/internal/domain/dto"
	"backend/internal/domain/service"
	"backend/internal/domain/utils"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"time"
)

type UserHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
	validator    *validator.Validator
}

func NewUserHandler(app *app.App) *UserHandler {
	userService := service.NewUserService(app.DB)
	tokenService := service.NewTokenService(app.JwtSecret, time.Hour)

	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
		validator:    app.Validator,
	}
}

// @Accept       json
// @Produce      json
// @Param        body body  dto.UserLogin true  "User login body object"
func (h UserHandler) login(ctx context.Context, input *dto.UserLoginInput) (*dto.TokenOutput, error) {
	userDTO := input.Body

	if errValidate := h.validator.ValidateData(userDTO); errValidate != nil {
		return nil, huma.Error400BadRequest(errValidate.Error(), errValidate)
	}

	telegramData, tgErr := utils.ParseInitData(userDTO.InitDataRaw)
	if tgErr != nil {
		return nil, huma.Error400BadRequest(tgErr.Error(), tgErr)
	}

	_, errFetch := h.userService.GetByID(ctx, telegramData.ID)
	if errFetch != nil {
		var createErr error
		createErr = h.userService.Create(ctx, repository.CreateUserParams{
			ID:   telegramData.ID,
			Name: "default username", // TODO: Change in PROD to real user data
		})

		if createErr != nil {
			return nil, huma.Error400BadRequest(createErr.Error(), createErr)
		}
	}

	token, tokenErr := h.tokenService.GenerateToken(telegramData.ID)
	if tokenErr != nil || token == "" {
		return nil, huma.Error500InternalServerError("failed to generate auth token")
	}

	resp := &dto.TokenOutput{
		Body: dto.Token{Token: token},
	}

	return resp, nil
}

func (h UserHandler) Setup(router huma.API) {
	huma.Register(router, huma.Operation{
		OperationID: "login",
		Path:        "/api/user/login",
		Method:      http.MethodPost,
		Errors: []int{
			400,
			500,
		},
		Tags: []string{
			"user",
		},
		Summary:     "Login to existing user account",
		Description: "Login to existing user account using his email, username and password. Returns his ID, email, username, verifiedEmail boolean variable and role",
	}, h.login)
}
