package usersConnector

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func CheckUserExists(id string, header string) (bool, error) {
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

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		log.Println("Error consulting users: ", resp.StatusCode, resp.Body)
		return false, fmt.Errorf("error consulting user: %d", resp.StatusCode)

	}
}
