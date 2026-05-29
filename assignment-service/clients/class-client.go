package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	redispkg "assignment-service/pkg/redis"
)

type ClassClient struct {
	BaseURL string
	Client  *http.Client
}

func NewClassClient() *ClassClient {
	return &ClassClient{
		BaseURL: os.Getenv("CLASS_SERVICE_URL"),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ClassClient) GetStudentsByClassID(ctx context.Context, classID string, token string) ([]redispkg.NotificationRecipient, string, error) {
	url := fmt.Sprintf("%s/classes/%s/students", c.BaseURL, classID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("class-service returned status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			ClassName string                           `json:"class_name"`
			Students  []redispkg.NotificationRecipient `json:"students"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	return result.Data.Students, result.Data.ClassName, nil
}
