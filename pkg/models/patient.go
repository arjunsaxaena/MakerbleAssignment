package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Patient struct {
	ID            string          `json:"id" db:"id"`
	Name          string          `json:"name" db:"name"`
	Age           int             `json:"age" db:"age"`
	Gender        string          `json:"gender" db:"gender"`
	Address       *string         `json:"address,omitempty" db:"address"`
	Diagnosis     *string         `json:"diagnosis,omitempty" db:"diagnosis"`
	RegisteredBy  *string         `json:"registered_by,omitempty" db:"registered_by"`
	LastUpdatedBy *string         `json:"last_updated_by,omitempty" db:"last_updated_by"`
	Meta          json.RawMessage `json:"meta" db:"meta"`
	IsActive      bool            `json:"is_active" db:"is_active"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

type GetPatientFilters struct {
	ID            string  `json:"id" form:"id"`
	Name          string  `json:"name" form:"name"`
	Age           int     `json:"age" form:"age"`
	Gender        string  `json:"gender" form:"gender"`
	Address       string  `json:"address" form:"address"`
	Diagnosis     string  `json:"diagnosis" form:"diagnosis"`
	RegisteredBy  string  `json:"registered_by" form:"registered_by"`
	LastUpdatedBy string  `json:"last_updated_by" form:"last_updated_by"`
	IsActive      *bool   `json:"is_active" form:"is_active"`
}

func (p *Patient) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	
	if p.Age <= 0 {
		return errors.New("age must be greater than 0")
	}
	
	if p.Gender != string(GenderMale) && p.Gender != string(GenderFemale) && p.Gender != string(GenderOther) {
		return errors.New("gender must be 'male', 'female', or 'other'")
	}
	if p.RegisteredBy == nil {
		return errors.New("registered_by is required")
	}
	
	return nil
} 