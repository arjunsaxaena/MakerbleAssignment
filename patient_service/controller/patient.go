package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arjunsaxaena/MakerbleAssignment/patient_service/repository"
	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/gin-gonic/gin"
)

type PatientController struct {
	repo repository.PatientRepository
}

func NewPatientController() *PatientController {
	return &PatientController{
		repo: repository.PatientRepository{},
	}
}

func (c *PatientController) Create(ctx *gin.Context) {
	var patient models.Patient
	if err := ctx.ShouldBindJSON(&patient); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "User ID not found in token",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Invalid user ID format in token",
		})
		return
	}

	patient.RegisteredBy = &userIDStr
	patient.LastUpdatedBy = &userIDStr

	if err := patient.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if patient.Meta == nil {
		patient.Meta = json.RawMessage("{}")
	}

	err := c.repo.Create(ctx, &patient)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to create patient: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    patient,
	})
}

func (c *PatientController) Get(ctx *gin.Context) {
	var filters models.GetPatientFilters

	filters.ID = ctx.Query("id")
	filters.Name = ctx.Query("name")
	
	if age := ctx.Query("age"); age != "" {
		var ageInt int
		if _, err := fmt.Sscanf(age, "%d", &ageInt); err == nil {
			filters.Age = ageInt
		}
	}
	
	filters.Gender = ctx.Query("gender")
	filters.Address = ctx.Query("address")
	filters.Diagnosis = ctx.Query("diagnosis")
	filters.RegisteredBy = ctx.Query("registered_by")
	filters.LastUpdatedBy = ctx.Query("last_updated_by")
	
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filters.IsActive = &isActive
	}

	patients, err := c.repo.Get(ctx, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch patients: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    patients,
	})
}

func (c *PatientController) Update(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Patient ID is required",
		})
		return
	}

	existingPatients, err := c.repo.Get(ctx, models.GetPatientFilters{ID: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch patient: " + err.Error(),
		})
		return
	}

	if len(existingPatients) == 0 {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Patient not found",
		})
		return
	}

	existingPatient := existingPatients[0]

	var updateRequest models.Patient
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "User ID not found in token",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Invalid user ID format in token",
		})
		return
	}

	updateRequest.LastUpdatedBy = &userIDStr

	updateRequest.ID = id

	if updateRequest.Name == "" {
		updateRequest.Name = existingPatient.Name
	}
	if updateRequest.Age == 0 {
		updateRequest.Age = existingPatient.Age
	}
	if updateRequest.Gender == "" {
		updateRequest.Gender = existingPatient.Gender
	}
	if updateRequest.Address == nil {
		updateRequest.Address = existingPatient.Address
	}
	if updateRequest.Diagnosis == nil {
		updateRequest.Diagnosis = existingPatient.Diagnosis
	}
	if updateRequest.RegisteredBy == nil {
		updateRequest.RegisteredBy = existingPatient.RegisteredBy
	}
	if updateRequest.Meta == nil {
		updateRequest.Meta = existingPatient.Meta
	}

	updateRequest.IsActive = existingPatient.IsActive

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
			Error:   "Failed to update patient: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    updateRequest,
	})
}

func (c *PatientController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Patient ID is required",
		})
		return
	}

	err := c.repo.Delete(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errors.New("no patient found to deactivate")) {
			status = http.StatusNotFound
		}
		ctx.JSON(status, models.Response{
			Success: false,
			Error:   "Failed to delete patient: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    map[string]string{"message": "Patient deleted successfully"},
	})
} 