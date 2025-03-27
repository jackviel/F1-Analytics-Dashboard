package models

import (
	"time"

	"gorm.io/gorm"
)

type Driver struct {
	gorm.Model
	Name            string    `gorm:"not null"`
	Nationality     string    `gorm:"not null"`
	DateOfBirth     time.Time `gorm:"not null"`
	Number          int       `gorm:"unique"`
	TeamID          uint      `gorm:"not null"`
	Team            Team      `gorm:"foreignKey:TeamID"`
	Races           []Race    `gorm:"many2many:race_drivers;"`
	CareerPoints    int       `gorm:"default:0"`
	CareerWins      int       `gorm:"default:0"`
	CareerPoles     int       `gorm:"default:0"`
	CareerFastLaps  int       `gorm:"default:0"`
	CareerPodiums   int       `gorm:"default:0"`
	Active          bool      `gorm:"default:true"`
	ProfileImageURL string
	Biography       string
}

// RaceDriver represents the many-to-many relationship between drivers and races
type RaceDriver struct {
	DriverID    uint      `gorm:"primaryKey"`
	RaceID      uint      `gorm:"primaryKey"`
	Position    int
	Points      float64
	Grid        int
	FastestLap  time.Duration
	RaceTime    time.Duration
	Status      string    // Finished, DNF, DNS, etc.
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
} 