package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID and set CreatedAt.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	base.CreatedAt = now
	base.UpdatedAt = now
	return nil
}

// BeforeUpdate will update the UpdatedAt field.
func (base *Base) BeforeUpdate(tx *gorm.DB) error {
	base.UpdatedAt = time.Now()
	return nil
}
