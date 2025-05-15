package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/database"
	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type UserRepository struct{}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	sb := sqlbuilder.NewStruct(new(models.User)).For(sqlbuilder.PostgreSQL)

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	query, args := sb.InsertInto("users", user).BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := database.DB.ExecContext(ctx, query, args...)
	return err
}


func (r *UserRepository) Get(ctx context.Context, filters models.GetUserFilters) ([]models.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*").From("users")

	if filters.ID != "" {
		sb.Where(sb.Equal("id", filters.ID))
	}
	if filters.Username != "" {
		sb.Where(sb.Equal("username", filters.Username))
	}
	if filters.Role != "" {
		sb.Where(sb.Equal("role", filters.Role))
	}
	if ctx.Value("is_active_query_present") == true {
		sb.Where(sb.Equal("is_active", filters.IsActive))
	}

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var users []models.User
	err := database.DB.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	structBuilder := sqlbuilder.NewStruct(new(models.User)).For(sqlbuilder.PostgreSQL)

	ub := structBuilder.Update("users", user)

	ub = ub.Where(ub.Equal("id", user.ID))

	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)

	fmt.Println("Generated Update Query:", query)
	fmt.Println("Arguments:", args)

	_, err := database.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("users")
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
		return errors.New("no user found to deactivate")
	}
	return nil
}


