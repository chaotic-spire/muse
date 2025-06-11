package api

import (
	"backend/cmd/app"
	"backend/internal/adapters/controller/api/validator"
	"backend/internal/adapters/database/postgres"
	"backend/internal/domain/dto"
	"backend/internal/domain/entity"
	"backend/internal/domain/service"
	"backend/internal/domain/utils/auth"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

type UserService interface {
	Create(ctx context.Context, registerReq dto.UserReturn) (*entity.User, error)
	GetByID(ctx context.Context, uuid string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type TokenService interface {
	GenerateAuthTokens(c context.Context, userID string) (*dto.AuthTokens, error)
	GenerateToken(ctx context.Context, userID string, expires time.Time, tokenType string) (string, error)
}

type UserHandler struct {
	userService  UserService
	tokenService TokenService
	validator    *validator.Validator
}

func NewUserHandler(app *app.App) *UserHandler {
	userStorage := postgres.NewUserStorage(app.DB)

	return &UserHandler{
		userService:  service.NewUserService(userStorage),
		tokenService: service.NewTokenService(),
		validator:    app.Validator,
	}
}

// @Accept       json
// @Produce      json
// @Param        body body  dto.UserLogin true  "User login body object"
func (h UserHandler) login(ctx context.Context, input *dto.UserLoginInput) (*dto.UserRegisterOutput, error) {
	userDTO := input.Body

	if errValidate := h.validator.ValidateData(userDTO); errValidate != nil {
		return nil, huma.Error400BadRequest(errValidate.Error(), errValidate)
	}

	telegramData, tgErr := auth.ParseInitData(userDTO.InitDataRaw)
	if tgErr != nil {
		return nil, huma.Error400BadRequest(tgErr.Error(), tgErr)
	}

	user, errFetch := h.userService.GetByID(ctx, strconv.FormatInt(telegramData.ID, 10))
	if errFetch != nil {
		var createErr error
		user, createErr = h.userService.Create(ctx, *telegramData)

		if createErr != nil {
			return nil, huma.Error400BadRequest(createErr.Error(), createErr)
		}
	}

	tokens, tokensErr := h.tokenService.GenerateAuthTokens(ctx, strconv.FormatInt(user.ID, 10))
	if tokensErr != nil || tokens == nil {
		return nil, huma.Error500InternalServerError("failed to generate auth tokens")
	}

	resp := &dto.UserRegisterOutput{
		Body: dto.UserRegisterResponse{
			User: dto.UserReturn{
				ID:        user.ID,
				Username:  user.Username,
				Firstname: user.Firstname,
				Lastname:  user.Lastname,
				PhotoUrl:  user.PhotoUrl,
			},
			Tokens: *tokens,
		},
	}

	return resp, nil
}

// @Accept       json
// @Produce      json
// @Param        body body  dto.Token true  "Access token object"
// @Success      200  {object}  dto.Token
func (h UserHandler) refreshToken(ctx context.Context, input *dto.TokenInput) (*dto.TokenOutput, error) {
	accessTokenDTO := input.Body

	if errValidate := h.validator.ValidateData(accessTokenDTO); errValidate != nil {
		return nil, huma.Error400BadRequest(errValidate.Error(), errValidate)
	}

	userID, errToken := auth.VerifyToken(accessTokenDTO.Token, viper.GetString("service.backend.jwt.secret"), auth.TokenTypeAccess)

	if errToken != nil {
		return nil, huma.Error401Unauthorized(errToken.Error(), errToken)
	}

	expTime := time.Now().UTC().Add(time.Minute * time.Duration(viper.GetInt("service.backend.jwt.access-token-expiration")))

	newAccess, errNewAccess := h.tokenService.GenerateToken(ctx,
		userID,
		expTime,
		auth.TokenTypeAccess)

	if errNewAccess != nil {
		return nil, huma.Error500InternalServerError(errNewAccess.Error(), errNewAccess)
	}

	resp := &dto.TokenOutput{
		Body: dto.Token{
			Token:   newAccess,
			Expires: expTime,
		},
	}

	return resp, nil
}

func (h UserHandler) Setup(router huma.API) {
	huma.Register(router, huma.Operation{
		OperationID: "refreshToken",
		Path:        "/api/user/refresh",
		Method:      http.MethodPost,
		Errors: []int{
			400,
			401,
			500,
		},
		Tags: []string{
			"user",
		},
		Summary:     "Refresh the access token",
		Description: "Get a new access token using a valid refresh token",
	}, h.refreshToken)

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
