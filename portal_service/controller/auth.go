package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/arjunsaxaena/MakerbleAssignment/portal_service/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = os.Getenv("JWT_SECRET")

type AuthController struct {
	userRepo repository.UserRepository
}

func NewAuthController() *AuthController {
	return &AuthController{
		userRepo: repository.UserRepository{},
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginRequest models.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
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