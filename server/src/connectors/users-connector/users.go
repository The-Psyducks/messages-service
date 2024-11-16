// coverage:ignore
package users_connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type profile struct {
	Username  string `json:"username"`
	ImagePath string `json:"picture_path"`
}

type UsersConnector struct {
}

type profileUserResponse struct {
	Profile *profile `json:"profile"`
}

func (uc *UsersConnector) GetUserNameAndImage(id string, header string) (string, string, error) {

	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", errors.New("error creating request:" + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", errors.New("error against user service:" + err.Error())
	}
	// bind response to struct

	var response profileUserResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", "", err

	}
	if response.Profile == nil {
		return "", "", errors.New("error getting user profile. user not found")
	}

	return response.Profile.Username, response.Profile.ImagePath, nil
}

func NewUsersConnector() Interface {
	return &UsersConnector{}
}

func (uc *UsersConnector) CheckUserExists(id string, header string) (bool, error) {
	log.Printf("Checking user exists")
	if os.Getenv("MOCK_USERS_SERVICE") == "true" {
		return true, nil
	}
	url := "http://" + os.Getenv("USERS_HOST") + "/users/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, errors.New("error creating request:" + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.New("error against user service:" + err.Error())
	}
	log.Println("user service response: ", resp.StatusCode)
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	//case http.Status
	case http.StatusNotFound:
		return false, nil
	default:
		log.Println("Error consulting users: ", resp.StatusCode, resp.Body)
		return false, fmt.Errorf("error consulting user: %d", resp.StatusCode)

	}
}
