package models

import (
	"time"

	"github.com/google/uuid"
)

type UserType int

const (
	FreeUser UserType = iota
	MonthlyMember
)

type Membership struct {
	Base
	UserID             uuid.UUID `gorm:"type:uuid;not null"`
	UserType           UserType  `gorm:"type:int;default:0"`
	DailyUsageLimit    int       `gorm:"type:int;default:10"`
	TotalUsageCount    int       `gorm:"type:int;default:0"`
	DailyUsageCount    int       `gorm:"type:int;default:0"`
	LastUsageDate      time.Time
	MembershipExpireAt time.Time
	BoosterExpireAt    time.Time
	BoosterUsageCount  int `gorm:"type:int;default:0"`
}

func (m *Membership) ResetDailyUsage() {
	m.DailyUsageCount = 0
	m.LastUsageDate = time.Now()
}

func (m *Membership) CheckAndUpdateExpiration() {
	now := time.Now()
	if m.UserType == MonthlyMember && now.After(m.MembershipExpireAt) {
		m.UserType = FreeUser
		// TODO DailyUsageLimit use config value
		m.DailyUsageLimit = 10 // Reset to free user limit
	}
}
