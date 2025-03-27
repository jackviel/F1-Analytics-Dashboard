package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	OpenF1BaseURL  = "https://api.openf1.org/v1"
	RateLimitDelay = 5 * time.Second // 5 second delay between requests
	MaxRetries     = 3               // Maximum number of retries for rate-limited requests
)

type OpenF1Service struct {
	client      *http.Client
	lastRequest time.Time
	mu          sync.Mutex
	requestChan chan struct{}
	cache       struct {
		sync.RWMutex
		drivers map[string][]Driver
		teams   map[string][]Team
		races   map[string][]Race
		session *Session
	}
}

func NewOpenF1Service() *OpenF1Service {
	service := &OpenF1Service{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		requestChan: make(chan struct{}, 1),
	}
	service.cache.drivers = make(map[string][]Driver)
	service.cache.teams = make(map[string][]Team)
	service.cache.races = make(map[string][]Race)
	return service
}

// rateLimit ensures we don't exceed the API rate limit
func (s *OpenF1Service) rateLimit() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if diff := now.Sub(s.lastRequest); diff < RateLimitDelay {
		time.Sleep(RateLimitDelay - diff)
	}
	s.lastRequest = time.Now()
}

// GetCurrentSeason returns the current F1 season
func GetCurrentSeason() int {
	now := time.Now()
	year := now.Year()
	// If we're in the first few months of the year, return previous year's season
	if now.Month() < 3 {
		return year - 1
	}
	return year
}

// Session represents a Formula 1 session (practice, qualifying, race)
type Session struct {
	SessionKey  int       `json:"session_key"`
	MeetingKey  int       `json:"meeting_key"`
	SessionName string    `json:"session_name"`
	CountryName string    `json:"country_name"`
	Year        int       `json:"year"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
}

// GetCurrentSession fetches the current or most recent F1 session
func (s *OpenF1Service) GetCurrentSession() (*Session, error) {
	s.rateLimit()

	// Check cache first
	s.cache.RLock()
	if s.cache.session != nil {
		session := *s.cache.session
		s.cache.RUnlock()
		return &session, nil
	}
	s.cache.RUnlock()

	url := fmt.Sprintf("%s/sessions", OpenF1BaseURL)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}
	defer resp.Body.Close()

	var sessions []Session
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return nil, fmt.Errorf("failed to decode sessions: %w", err)
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found")
	}

	// Find the most recent session by comparing dates
	now := time.Now()
	var mostRecent *Session
	for i := range sessions {
		session := &sessions[i]
		// Skip future sessions
		if session.DateStart.After(now) {
			continue
		}

		// Check if this session has driver data
		s.rateLimit()
		url := fmt.Sprintf("%s/drivers?session_key=%d", OpenF1BaseURL, session.SessionKey)
		driverResp, err := s.client.Get(url)
		if err != nil {
			continue
		}
		defer driverResp.Body.Close()

		var drivers []Driver
		if err := json.NewDecoder(driverResp.Body).Decode(&drivers); err != nil {
			continue
		}

		// Skip sessions with no driver data
		if len(drivers) == 0 {
			continue
		}

		// If we haven't found a session yet, or if this session is more recent
		if mostRecent == nil || session.DateStart.After(mostRecent.DateStart) {
			mostRecent = session
		}
	}

	if mostRecent == nil {
		return nil, fmt.Errorf("no recent sessions with driver data found")
	}

	// Cache the result
	s.cache.Lock()
	s.cache.session = mostRecent
	s.cache.Unlock()

	return mostRecent, nil
}

// Driver represents a Formula 1 driver
type Driver struct {
	SessionKey    int    `json:"session_key"`
	MeetingKey    int    `json:"meeting_key"`
	BroadcastName string `json:"broadcast_name"`
	CountryCode   string `json:"country_code"`
	FirstName     string `json:"first_name"`
	FullName      string `json:"full_name"`
	HeadshotURL   string `json:"headshot_url"`
	LastName      string `json:"last_name"`
	DriverNumber  int    `json:"driver_number"`
	TeamColor     string `json:"team_colour"`
	TeamName      string `json:"team_name"`
	NameAcronym   string `json:"name_acronym"`
}

// makeRequest makes an HTTP request with rate limiting and retries
func (s *OpenF1Service) makeRequest(url string) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Acquire request token
	s.requestChan <- struct{}{}
	defer func() { <-s.requestChan }() // Release token when done

	for i := 0; i < MaxRetries; i++ {
		s.rateLimit()
		resp, err = s.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		// If we get a 429, wait and retry
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			time.Sleep(RateLimitDelay * time.Duration(i+1)) // Exponential backoff
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded, last error: %v", err)
}

// GetDrivers fetches F1 drivers with optional filtering
func (s *OpenF1Service) GetDrivers(season *int, meetingKey *int, sessionKey *int, teamName *string) ([]Driver, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf("%d-%d-%d-%s",
		getIntValue(season),
		getIntValue(meetingKey),
		getIntValue(sessionKey),
		getStringValue(teamName))

	// Check cache first
	s.cache.RLock()
	if drivers, ok := s.cache.drivers[cacheKey]; ok {
		s.cache.RUnlock()
		return drivers, nil
	}
	s.cache.RUnlock()

	// If no session key provided, get the most recent session
	if sessionKey == nil {
		sessions, err := s.GetCurrentSession()
		if err != nil {
			return nil, fmt.Errorf("failed to get current session: %w", err)
		}
		sessionKey = &sessions.SessionKey
	}

	// Try to get drivers for the specified session
	url := fmt.Sprintf("%s/drivers?session_key=%d", OpenF1BaseURL, *sessionKey)

	// Add additional query parameters if provided
	params := make([]string, 0)
	if season != nil {
		params = append(params, fmt.Sprintf("season=%d", *season))
	}
	if meetingKey != nil {
		params = append(params, fmt.Sprintf("meeting_key=%d", *meetingKey))
	}
	if teamName != nil {
		params = append(params, fmt.Sprintf("team_name=%s", *teamName))
	}

	// Add query parameters to URL
	if len(params) > 0 {
		url += "&" + strings.Join(params, "&")
	}

	fmt.Printf("Fetching drivers from URL: %s\n", url)
	resp, err := s.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch drivers: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	fmt.Printf("Response body: %s\n", string(body))

	var drivers []Driver
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&drivers); err != nil {
		return nil, fmt.Errorf("failed to decode drivers: %w", err)
	}

	fmt.Printf("Decoded %d drivers\n", len(drivers))
	for _, driver := range drivers {
		fmt.Printf("Driver: %+v\n", driver)
	}

	// Cache the result
	s.cache.Lock()
	s.cache.drivers[cacheKey] = drivers
	s.cache.Unlock()

	return drivers, nil
}

// GetTeams fetches current F1 teams
func (s *OpenF1Service) GetTeams() ([]Team, error) {
	// Check cache first
	s.cache.RLock()
	if teams, ok := s.cache.teams["current"]; ok {
		s.cache.RUnlock()
		return teams, nil
	}
	s.cache.RUnlock()

	session, err := s.GetCurrentSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get current session: %w", err)
	}

	url := fmt.Sprintf("%s/teams?session_key=%d", OpenF1BaseURL, session.SessionKey)
	resp, err := s.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}
	defer resp.Body.Close()

	var teams []Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("failed to decode teams: %w", err)
	}

	// Cache the result
	s.cache.Lock()
	s.cache.teams["current"] = teams
	s.cache.Unlock()

	return teams, nil
}

// GetRaces fetches race information
func (s *OpenF1Service) GetRaces() ([]Race, error) {
	// Check cache first
	s.cache.RLock()
	if races, ok := s.cache.races["current"]; ok {
		s.cache.RUnlock()
		return races, nil
	}
	s.cache.RUnlock()

	session, err := s.GetCurrentSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get current session: %w", err)
	}

	url := fmt.Sprintf("%s/races?session_key=%d", OpenF1BaseURL, session.SessionKey)
	resp, err := s.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch races: %w", err)
	}
	defer resp.Body.Close()

	var races []Race
	if err := json.NewDecoder(resp.Body).Decode(&races); err != nil {
		return nil, fmt.Errorf("failed to decode races: %w", err)
	}

	// Cache the result
	s.cache.Lock()
	s.cache.races["current"] = races
	s.cache.Unlock()

	return races, nil
}

// GetCircuits fetches circuit information
func (s *OpenF1Service) GetCircuits() ([]Circuit, error) {
	url := fmt.Sprintf("%s/circuits", OpenF1BaseURL)
	resp, err := s.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch circuits: %w", err)
	}
	defer resp.Body.Close()

	var circuits []Circuit
	if err := json.NewDecoder(resp.Body).Decode(&circuits); err != nil {
		return nil, fmt.Errorf("failed to decode circuits: %w", err)
	}

	return circuits, nil
}

// GetRaceResults fetches results for a specific race
func (s *OpenF1Service) GetRaceResults(raceID string) ([]RaceResult, error) {
	url := fmt.Sprintf("%s/race_results?race_id=%s", OpenF1BaseURL, raceID)
	resp, err := s.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch race results: %w", err)
	}
	defer resp.Body.Close()

	var results []RaceResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode race results: %w", err)
	}

	return results, nil
}

// Data structures matching OpenF1 API response
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Race struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Circuit string `json:"circuit"`
	Date    string `json:"date"`
	Round   int    `json:"round"`
	Status  string `json:"status"`
}

type Circuit struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Country  string `json:"country"`
}

type RaceResult struct {
	Position     int     `json:"position"`
	DriverNumber int     `json:"driver_number"`
	DriverName   string  `json:"driver_name"`
	Team         string  `json:"team"`
	Laps         int     `json:"laps"`
	Time         string  `json:"time"`
	Points       float64 `json:"points"`
}

// Helper functions for cache keys
func getIntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
