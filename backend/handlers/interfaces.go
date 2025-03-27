package handlers

import "github.com/f1-analytics/services"

// OpenF1Service defines the interface for interacting with the OpenF1 API
type OpenF1Service interface {
	GetDrivers(season *int, meetingKey *int, sessionKey *int, teamName *string) ([]services.Driver, error)
	GetTeams() ([]services.Team, error)
	GetRaces() ([]services.Race, error)
	GetRaceResults(raceID string) ([]services.RaceResult, error)
	GetCurrentSession() (*services.Session, error)
}
