package handlers

import (
	"net/http"
	"strconv"

	"github.com/f1-analytics/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DriverHandler struct {
	openF1Service OpenF1Service
	db            *gorm.DB
}

func NewDriverHandler(openF1Service OpenF1Service, db *gorm.DB) *DriverHandler {
	return &DriverHandler{
		openF1Service: openF1Service,
		db:            db,
	}
}

// GetDrivers returns all F1 drivers with optional filtering
func (h *DriverHandler) GetDrivers(c *gin.Context) {
	var drivers []models.Driver
	result := h.db.Find(&drivers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch drivers from database",
		})
		return
	}

	// If no drivers in database, fetch from OpenF1 API and store them
	if len(drivers) == 0 {
		// Parse query parameters
		var season *int
		if seasonStr := c.Query("season"); seasonStr != "" {
			if s, err := strconv.Atoi(seasonStr); err == nil {
				season = &s
			}
		}

		var meetingKey *int
		if meetingKeyStr := c.Query("meeting_key"); meetingKeyStr != "" {
			if mk, err := strconv.Atoi(meetingKeyStr); err == nil {
				meetingKey = &mk
			}
		}

		var sessionKey *int
		if sessionKeyStr := c.Query("session_key"); sessionKeyStr != "" {
			if sk, err := strconv.Atoi(sessionKeyStr); err == nil {
				sessionKey = &sk
			}
		}

		var teamName *string
		if team := c.Query("team"); team != "" {
			teamName = &team
		}

		apiDrivers, err := h.openF1Service.GetDrivers(season, meetingKey, sessionKey, teamName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch drivers from API",
			})
			return
		}

		// Convert API drivers to database models
		for _, apiDriver := range apiDrivers {
			driver := models.Driver{
				Number:      apiDriver.DriverNumber,
				Name:        apiDriver.BroadcastName,
				Nationality: apiDriver.CountryCode,
				Active:      true,
				// Set other fields as needed
			}
			if err := h.db.Create(&driver).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to store drivers in database",
				})
				return
			}
			drivers = append(drivers, driver)
		}
	}

	// Transform drivers to match expected response format
	type DriverResponse struct {
		DriverNumber int    `json:"driver_number"`
		Name         string `json:"name"`
		Team         string `json:"team"`
		Country      string `json:"country"`
	}

	response := make([]DriverResponse, len(drivers))
	for i, driver := range drivers {
		response[i] = DriverResponse{
			DriverNumber: driver.Number,
			Name:         driver.Name,
			Team:         driver.Team.Name,
			Country:      driver.Nationality,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetDriver returns a specific driver by ID
func (h *DriverHandler) GetDriver(c *gin.Context) {
	driverNumber, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid driver number",
		})
		return
	}

	var driver models.Driver
	result := h.db.First(&driver, "number = ?", driverNumber)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Driver not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch driver from database",
		})
		return
	}

	c.JSON(http.StatusOK, driver)
}

// GetDriverStats returns statistics for a specific driver
func (h *DriverHandler) GetDriverStats(c *gin.Context) {
	driverNumber, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid driver number",
		})
		return
	}

	var driver models.Driver
	result := h.db.First(&driver, "number = ?", driverNumber)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Driver not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch driver from database",
		})
		return
	}

	stats := gin.H{
		"wins":        driver.CareerWins,
		"podiums":     driver.CareerPodiums,
		"points":      driver.CareerPoints,
		"poles":       driver.CareerPoles,
		"fastestLaps": driver.CareerFastLaps,
	}

	c.JSON(http.StatusOK, stats)
}
