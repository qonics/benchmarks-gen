package model

import (
	"time"
)

type ActivityLogs struct {
	Id        uint      `gorm:"primary_key;autoIncrement"`
	EventType string    `gorm:"index;default null;index"`
	Type      string    `gorm:"index;default null;index"`
	Activity  string    `gorm:"not null"`
	IpAddress string    `gorm:"not null"`
	Agent     string    `gorm:"not null"`
	Extra     string    `gorm:"type:text;null;default null"`
	Operator  uint      `gorm:"index;size:11;default null"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}
