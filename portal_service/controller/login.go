package controller

import (
	"net/http"

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/middleware"
	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/arjunsaxaena/MakerbleAssignment/portal_service/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	userRepo repository.UserRepository
}

func NewLoginController() *LoginController {
	return &LoginController{
		userRepo: repository.UserRepository{},
	}
}

func (c *LoginController) Login(ctx *gin.Context) {
	var loginRequest models.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	jwtSecret := middleware.GetJWTSecret()
	if jwtSecret == "" {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "JWT_SECRET not configured, authentication cannot proceed",
		})
		return
	}

	users, err := c.userRepo.Get(ctx, models.GetUserFilters{
		Username: loginRequest.Username,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to authenticate user: " + err.Error(),
		})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}

	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}

	tokenString, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to generate authentication token",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.TokenResponse{
			Token: tokenString,
		},
	})
} 