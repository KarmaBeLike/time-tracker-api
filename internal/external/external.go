package external

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KarmaBeLike/time-tracker-api/config"
)

type People struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func FetchUserData(config config.Config, passportSerie, passportNumber string) (*People, error) {
	url := fmt.Sprintf(config.PeopleAPIBaseURL+"/info?passportSerie=%s&passportNumber=%s", passportSerie, passportNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user data: status code %d", resp.StatusCode)
	}

	var people People
	err = json.NewDecoder(resp.Body).Decode(&people)
	if err != nil {
		return nil, err
	}

	return &people, nil
}
