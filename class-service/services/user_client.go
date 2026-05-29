package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

// GetUserProfile requests the user-service /profile endpoint using the provided raw JWT token.
// Returns the `data` object from the user-service response as a map.
func GetUserProfile(rawToken string) (map[string]interface{}, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8080"
	}

	req, err := http.NewRequest("GET", userServiceURL+"/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+rawToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch profile from user-service")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	data, ok := body["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid profile response format")
	}

	return data, nil
}
