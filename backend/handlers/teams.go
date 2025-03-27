package handlers

import (
	"net/http"
	"strconv"

	"github.com/f1-analytics/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeamHandler struct {
	openF1Service OpenF1Service
	db            *gorm.DB
}

func NewTeamHandler(openF1Service OpenF1Service, db *gorm.DB) *TeamHandler {
	return &TeamHandler{
		openF1Service: openF1Service,
		db:            db,
	}
}

// GetTeams returns all F1 teams
func (h *TeamHandler) GetTeams(c *gin.Context) {
	var teams []models.Team
	result := h.db.Find(&teams)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch teams from database",
		})
		return
	}

	// If no teams in database, fetch from OpenF1 API and store them
	if len(teams) == 0 {
		apiTeams, err := h.openF1Service.GetTeams()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch teams from API",
			})
			return
		}

		teams = []models.Team{} // Reset teams slice
		// Convert API teams to database models
		for _, apiTeam := range apiTeams {
			team := models.Team{
				Name:         apiTeam.Name,
				Active:       true,
				Nationality:  "Unknown", // Required field, set default
				BaseLocation: "Unknown",
				WorldTitles:  0,
				RaceWins:     0,
				Poles:        0,
				FastestLaps:  0,
				Podiums:      0,
				Points:       0,
			}
			if err := h.db.Create(&team).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to store teams in database",
				})
				return
			}
			teams = append(teams, team)
		}
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeam returns a specific team by ID
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	var team models.Team
	result := h.db.First(&team, teamID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Team not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch team from database",
		})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetTeamStats returns statistics for a specific team
func (h *TeamHandler) GetTeamStats(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	var team models.Team
	result := h.db.First(&team, teamID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Team not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch team from database",
		})
		return
	}

	stats := gin.H{
		"worldTitles": team.WorldTitles,
		"raceWins":    team.RaceWins,
		"poles":       team.Poles,
		"fastestLaps": team.FastestLaps,
		"podiums":     team.Podiums,
		"points":      team.Points,
	}

	c.JSON(http.StatusOK, stats)
}
