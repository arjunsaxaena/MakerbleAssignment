package models

import (
	"encoding/json"
	"errors"
	"time"
)

type UserRole string

const (
	RoleReceptionist UserRole = "receptionist"
	RoleDoctor       UserRole = "doctor"
)

type User struct {
	ID        string          `json:"id" db:"id"`
	Username  string          `json:"username" db:"username"`
	Password  string          `json:"password" db:"password"`
	Role      string          `json:"role" db:"role"`
	Meta      json.RawMessage `json:"meta" db:"meta"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type GetUserFilters struct {
	ID        string `json:"id" form:"id"`
	Username  string `json:"username" form:"username"`
	Role      string `json:"role" form:"role"`
	IsActive  bool   `json:"is_active" form:"is_active"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	
	if u.Password == "" {
		return errors.New("password is required")
	}
	
	if u.Role != string(RoleReceptionist) && u.Role != string(RoleDoctor) {
		return errors.New("role must be 'receptionist' or 'doctor'")
	}
	
	return nil
}
