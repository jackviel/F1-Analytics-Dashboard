package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name            string    `gorm:"not null;unique"`
	Nationality     string    `gorm:"not null"`
	Founded         time.Time
	BaseLocation    string
	TeamPrincipal   string
	TechnicalChief  string
	Chassis         string
	Engine          string
	FirstEntry      time.Time
	WorldTitles     int       `gorm:"default:0"`
	RaceWins        int       `gorm:"default:0"`
	Poles           int       `gorm:"default:0"`
	FastestLaps     int       `gorm:"default:0"`
	Podiums         int       `gorm:"default:0"`
	Points          int       `gorm:"default:0"`
	Active          bool      `gorm:"default:true"`
	LogoURL         string
	Website         string
	Drivers         []Driver  `gorm:"foreignKey:TeamID"`
	Races           []Race    `gorm:"many2many:race_teams;"`
}

// RaceTeam represents the many-to-many relationship between teams and races
type RaceTeam struct {
	TeamID      uint      `gorm:"primaryKey"`
	RaceID      uint      `gorm:"primaryKey"`
	Points      float64
	Position    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
} 