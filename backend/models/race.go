package models

import (
	"time"

	"gorm.io/gorm"
)

type Circuit struct {
	gorm.Model
	Name            string    `gorm:"not null;unique"`
	Location        string    `gorm:"not null"`
	Country         string    `gorm:"not null"`
	Length          float64   // Circuit length in kilometers
	FirstGrandPrix  time.Time
	LapRecord       time.Duration
	LapRecordHolder string
	LapRecordYear   int
	ImageURL        string
	Description     string
	Races           []Race    `gorm:"foreignKey:CircuitID"`
}

type Race struct {
	gorm.Model
	Name            string    `gorm:"not null"`
	Season          int       `gorm:"not null"`
	Round           int       `gorm:"not null"`
	CircuitID       uint      `gorm:"not null"`
	Circuit         Circuit   `gorm:"foreignKey:CircuitID"`
	Date            time.Time `gorm:"not null"`
	RaceTime        time.Time
	QualifyingTime  time.Time
	Practice1Time   time.Time
	Practice2Time   time.Time
	Practice3Time   time.Time
	SprintTime      time.Time // Optional, for sprint races
	Status          string    // Scheduled, Completed, Cancelled
	Laps            int
	RaceDistance    float64   // Total race distance in kilometers
	Weather         string
	Temperature     float64
	TrackCondition  string
	Drivers         []Driver  `gorm:"many2many:race_drivers;"`
	Teams           []Team    `gorm:"many2many:race_teams;"`
	Results         []RaceDriver
	TeamResults     []RaceTeam
}

// Lap represents individual lap times for a driver in a race
type Lap struct {
	gorm.Model
	RaceID      uint      `gorm:"not null"`
	DriverID    uint      `gorm:"not null"`
	LapNumber   int       `gorm:"not null"`
	LapTime     time.Duration
	Position    int
	IsFastest   bool      `gorm:"default:false"`
	PitStop     bool      `gorm:"default:false"`
	PitStopTime time.Duration
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
} 