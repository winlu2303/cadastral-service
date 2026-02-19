package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cadastral-service/internal/config"
	"cadastral-service/internal/models"
	"cadastral-service/internal/repository"
)

type Service struct {
	repo *repository.Repository
	cfg  *config.Config
}

type ExternalServerResponse struct {
	Result bool    `json:"result"`
	Delay  float64 `json:"delay"`
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

// ProcessQuery is proccess request asynchron
func (s *Service) ProcessQuery(query *models.Query) {
	// update status on processing
	if err := s.repo.UpdateQuery(nil, query.ID, "processing", nil); err != nil {
		log.Printf("Failed to update query status: %v", err)
		return
	}

	// imitate sending on external server
	result, err := s.callExternalServer(query)
	if err != nil {
		log.Printf("Failed to call external server: %v", err)
		s.repo.UpdateQuery(nil, query.ID, "failed", nil)
		return
	}

	// update a result
	if err := s.repo.UpdateQuery(nil, query.ID, "completed", &result); err != nil {
		log.Printf("Failed to update query result: %v", err)
	}
}

// callExternalServer is call for a external emulate server
func (s *Service) callExternalServer(query *models.Query) (bool, error) {
	// make a data for sending
	requestData := map[string]interface{}{
		"cadastral_number": query.CadastralNumber,
		"latitude":         query.Latitude,
		"longitude":        query.Longitude,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false, err
	}

	// send request on an external server
	externalServerURL := s.cfg.ExternalServerURL
	if externalServerURL == "" {
		externalServerURL = "http://localhost:" + s.cfg.Port + "/api/result"
	}

	req, err := http.NewRequest("POST", externalServerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 65 * time.Second, // 60s + 5 for buffer
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("external server returned status: %d", resp.StatusCode)
	}

	// lets parcing an answer
	var response ExternalServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	return response.Result, nil
}
