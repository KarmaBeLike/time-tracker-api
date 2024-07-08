package external

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PeopleAPIClient struct {
	BaseURL string
}

type PeopleResponse struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func NewPeopleAPIClient(baseURL string) *PeopleAPIClient {
	return &PeopleAPIClient{BaseURL: baseURL}
}

func (client *PeopleAPIClient) GetPersonInfo(passportNumber string) (*PeopleResponse, error) {
	url := fmt.Sprintf("%s/info?passportNumber=%s", client.BaseURL, passportNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get person info, status code: %d", resp.StatusCode)
	}

	var peopleResponse PeopleResponse
	err = json.NewDecoder(resp.Body).Decode(&peopleResponse)
	if err != nil {
		return nil, err
	}

	return &peopleResponse, nil
}
