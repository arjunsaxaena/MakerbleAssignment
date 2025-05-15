package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/arjunsaxaena/MakerbleAssignment/portal_service/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo repository.UserRepository
}

func NewUserController() *UserController {
	return &UserController{
		repo: repository.UserRepository{},
	}
}

func (c *UserController) Create(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := user.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to process password",
		})
		return
	}
	user.Password = string(hashedPassword)

	if user.Meta == nil {
		user.Meta = json.RawMessage("{}")
	}

	err = c.repo.Create(ctx, &user)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "Failed to create user: " + err.Error()
		
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			status = http.StatusConflict
			errMsg = "User with this username already exists"
		}
		
		ctx.JSON(status, models.Response{
			Success: false,
			Error:   errMsg,
		})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    user,
	})
}

func (c *UserController) Get(ctx *gin.Context) {
	var filters models.GetUserFilters

	filters.ID = ctx.Query("id")
	filters.Username = ctx.Query("username")
	filters.Role = ctx.Query("role")
	filters.IsActive = ctx.Query("is_active") == "true"

	// if filters.ID != "" {
	// 	users, err := c.repo.Get(ctx, filters)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, models.Response{
	// 			Success: false,
	// 			Error:   "Failed to fetch user: " + err.Error(),
	// 		})
	// 		return
	// 	}

	// 	if len(users) == 0 {
	// 		ctx.JSON(http.StatusNotFound, models.Response{
	// 			Success: false,
	// 			Error:   "User not found",
	// 		})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusOK, models.Response{
	// 		Success: true,
	// 		Data:    users[0],
	// 	})
	// 	return
	// }

	users, err := c.repo.Get(ctx, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch users: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    users,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "User ID is required",
		})
		return
	}

	existingUsers, err := c.repo.Get(ctx, models.GetUserFilters{ID: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch user: " + err.Error(),
		})
		return
	}

	if len(existingUsers) == 0 {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	existingUser := existingUsers[0]

	var updateRequest models.User
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	updateRequest.ID = id

	if updateRequest.Username == "" {
		updateRequest.Username = existingUser.Username
	}

	if updateRequest.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Error:   "Failed to process password",
			})
			return
		}
		updateRequest.Password = string(hashedPassword)
	} else {
		updateRequest.Password = existingUser.Password
	}

	if updateRequest.Role == "" {
		updateRequest.Role = existingUser.Role
	}

	if updateRequest.Meta == nil {
		updateRequest.Meta = existingUser.Meta
	}

	updateRequest.IsActive = existingUser.IsActive

	if err := updateRequest.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.repo.Update(ctx, &updateRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to update user: " + err.Error(),
		})
		return
	}

	updateRequest.Password = ""

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    updateRequest,
	})
}

func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "User ID is required",
		})
		return
	}

	err := c.repo.Delete(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errors.New("no user found to delete")) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, models.Response{
			Success: false,
			Error:   "Failed to delete user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    map[string]string{"message": "User deleted successfully"},
	})
}