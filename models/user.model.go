package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleShopOwner  UserRole = "shop_owner"
	RoleUser       UserRole = "user"
)

type UserRoles []UserRole

func (r *UserRoles) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var roles []UserRole
	err := json.Unmarshal(bytes, &roles)
	*r = roles
	return err
}

func (r UserRoles) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	return json.Marshal(r)
}

func (r UserRoles) StringSlice() []string {
	slice := make([]string, len(r))
	for i, role := range r {
		slice[i] = string(role)
	}
	return slice
}

func (u *User) AddRole(role UserRole) {
	for _, r := range u.Roles {
		if r == role {
			return
		}
	}
	u.Roles = append(u.Roles, role)
}

func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) RemoveRole(role UserRole) {
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			return
		}
	}
}

type User struct {
	Base
	Name               string    `gorm:"type:varchar(255);not null"`
	Email              string    `gorm:"uniqueIndex;not null"`
	Password           string    `gorm:"not null"`
	Roles              UserRoles `gorm:"type:jsonb"`
	Provider           string    `gorm:"not null"`
	Photo              string    `gorm:"not null;default:'default.png'"`
	VerificationCode   string
	PasswordResetToken string
	PasswordResetAt    time.Time
	Verified           bool       `gorm:"not null"`
	Shops              []Shop     `gorm:"foreignKey:OwnerID"`
	Orders             []Order    `gorm:"foreignKey:UserID"`
	Membership         Membership `gorm:"foreignKey:UserID"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo           string `json:"photo"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Roles     UserRoles `json:"roles,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ur *UserResponse) MarshalJSON() ([]byte, error) {
	type Alias UserResponse
	return json.Marshal(&struct {
		*Alias
		Roles []string `json:"roles,omitempty"`
	}{
		Alias: (*Alias)(ur),
		Roles: ur.Roles.StringSlice(),
	})
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
