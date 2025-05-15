package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/database"
	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type PatientRepository struct{}

func (r *PatientRepository) Create(ctx context.Context, patient *models.Patient) error {
	sb := sqlbuilder.NewStruct(new(models.Patient)).For(sqlbuilder.PostgreSQL)

	patient.ID = uuid.New().String()
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()
	patient.IsActive = true

	query, args := sb.InsertInto("patients", patient).BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := database.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *PatientRepository) Get(ctx context.Context, filters models.GetPatientFilters) ([]models.Patient, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*").From("patients")

	if filters.ID != "" {
		sb.Where(sb.Equal("id", filters.ID))
	}
	if filters.Name != "" {
		sb.Where(sb.Equal("name", filters.Name))
	}
	if filters.Age > 0 {
		sb.Where(sb.Equal("age", filters.Age))
	}
	if filters.Gender != "" {
		sb.Where(sb.Equal("gender", filters.Gender))
	}
	if filters.Address != "" {
		sb.Where(sb.Equal("address", filters.Address))
	}
	if filters.Diagnosis != "" {
		sb.Where(sb.Equal("diagnosis", filters.Diagnosis))
	}
	if filters.RegisteredBy != "" {
		sb.Where(sb.Equal("registered_by", filters.RegisteredBy))
	}
	if filters.LastUpdatedBy != "" {
		sb.Where(sb.Equal("last_updated_by", filters.LastUpdatedBy))
	}
	if filters.IsActive != nil {
		sb.Where(sb.Equal("is_active", *filters.IsActive))
	}

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var patients []models.Patient
	err := database.DB.SelectContext(ctx, &patients, query, args...)
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *PatientRepository) Update(ctx context.Context, patient *models.Patient) error {
	patient.UpdatedAt = time.Now()

	structBuilder := sqlbuilder.NewStruct(new(models.Patient)).For(sqlbuilder.PostgreSQL)

	ub := structBuilder.Update("patients", patient)
	ub = ub.Where(ub.Equal("id", patient.ID))

	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := database.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *PatientRepository) Delete(ctx context.Context, id string) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("patients")
	sb.Set(
		sb.Assign("is_active", false),
		sb.Assign("updated_at", time.Now()),
	)
	sb.Where(sb.Equal("id", id))

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	result, err := database.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no patient found to deactivate")
	}
	return nil
} 