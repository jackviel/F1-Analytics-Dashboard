package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/f1-analytics/models"
	"github.com/f1-analytics/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RaceHandler struct {
	openF1Service OpenF1Service
	db            *gorm.DB
}

func NewRaceHandler(openF1Service OpenF1Service, db *gorm.DB) *RaceHandler {
	return &RaceHandler{
		openF1Service: openF1Service,
		db:            db,
	}
}

// GetRaces returns all F1 races
func (h *RaceHandler) GetRaces(c *gin.Context) {
	var races []models.Race
	result := h.db.Find(&races)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch races from database",
		})
		return
	}

	// If no races in database, fetch from OpenF1 API and store them
	if len(races) == 0 {
		apiRaces, err := h.openF1Service.GetRaces()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch races from API",
			})
			return
		}

		// Convert API races to database models
		for _, apiRace := range apiRaces {
			race := models.Race{
				Name:    apiRace.Name,
				Season:  services.GetCurrentSeason(), // Use helper function from services package
				Round:   apiRace.Round,
				Circuit: models.Circuit{Name: apiRace.Circuit},
				Date:    time.Now(), // TODO: Parse apiRace.Date properly
				Status:  apiRace.Status,
				// Set other fields as needed
			}
			if err := h.db.Create(&race).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to store races in database",
				})
				return
			}
			races = append(races, race)
		}
	}

	c.JSON(http.StatusOK, races)
}

// GetRace returns a specific race by ID
func (h *RaceHandler) GetRace(c *gin.Context) {
	raceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid race ID",
		})
		return
	}

	var race models.Race
	result := h.db.First(&race, raceID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Race not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch race from database",
		})
		return
	}

	c.JSON(http.StatusOK, race)
}

// GetRaceResults returns results for a specific race
func (h *RaceHandler) GetRaceResults(c *gin.Context) {
	raceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid race ID",
		})
		return
	}

	var race models.Race
	result := h.db.Preload("Drivers").First(&race, raceID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Race not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch race from database",
		})
		return
	}

	// Get race results from the join table
	var results []models.RaceDriver
	if err := h.db.Where("race_id = ?", raceID).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch race results",
		})
		return
	}

	c.JSON(http.StatusOK, results)
}
