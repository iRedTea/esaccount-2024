package service

import (
	"bytes"
	"encoding/json"
	"esaccount"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Authorize(header string) (*esaccount.AuthorizedUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://auth.easystartup.su/api/user/self", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", header)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authozire: %d %s\nheader: %s", resp.StatusCode, resp.Status, header)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result esaccount.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) AuthorizeById(id int64) (*esaccount.AuthorizedUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://auth.easystartup.su/api/user/%d", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("ADMIN_TOKEN"))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authozire: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result esaccount.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


func (s *AuthService) AuthorizeAndUpdatePicture(header string, picUrl string) (*esaccount.AuthorizedUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://auth.easystartup.su/api/user/self", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", header)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authozire: %d %s\nheader: %s", resp.StatusCode, resp.Status, header)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result esaccount.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	result.PictureURL = picUrl
	requestBody, _ := json.Marshal(result)

	req, err = http.NewRequest("POST", "https://auth.easystartup.su/api/user/self", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", header)
	req.Header.Add("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}