package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	Email          string        `json:"email,omitempty"`
	FirstName      string        `json:"first_name,omitempty"`
	LastName       string        `json:"last_name,omitempty"`
	Status         string        `json:"status,omitempty"`
	Type           int           `json:"type,omitempty"`
	Pmi            int           `json:"pmi,omitempty"`
	UsePmi         *bool         `json:"use_pmi,omitempty"`
	Timezone       string        `json:"timezone,omitempty"`
	Language       string        `json:"language,omitempty"`
	VanityName     string        `json:"vanity_name,omitempty"`
	HostKey        string        `json:"host_key,omitempty"`
	CmsUserId      string        `json:"cms_user_id,omitempty"`
	Company        string        `json:"company,omitempty"`
	GroupId        string        `json:"group_id,omitempty"`
	Manager        string        `json:"manager,omitempty"`
	Pronouns       string        `json:"pronouns,omitempty"`
	PhoneNumbers   []PhoneNumber `json:"phone_numbers,omitempty"`
	PronounsOption int           `json:"pronouns_option,omitempty"`
	RoleName       string        `json:"role_name,omitempty"`
	Department     string        `json:"dept,omitempty"`
	JobTitle       string        `json:"job_title,omitempty"`
	Location       string        `json:"location,omitempty"`
	Id             string        `json:"id,omitempty"`
}

type UserInfo struct {
	Id string `json:"id,omitempty"`
}

type PhoneNumber struct {
	Country string `json:"country,omitempty"`
	Code    string `json:"code,omitempty"`
	Number  string `json:"number,omitempty"`
	Label   string `json:"label,omitempty"`
}

type Client struct {
	authToken      string
	TimeoutMinutes int
	httpClient     *http.Client
}

type AuthInfo struct {
	AccessToken string `json:"access_token,omitempty"`
}

func NewClient(authToken string, timeoutMinutes int) *Client {
	return &Client{
		authToken:      authToken,
		TimeoutMinutes: timeoutMinutes,
		httpClient:     &http.Client{},
	}
}

func (c *Client) NewUser(user *User) (string, error) {
	userInfo, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	body := strings.NewReader(fmt.Sprintf("{\"action\":\"create\",\"user_info\":" + string(userInfo) + "}"))
	res, err := c.httpRequest("POST", body, "")
	if err != nil {
		log.Println("[CREATE ERROR]: ", err)
		return "", err
	}
	created := &UserInfo{}
	err = json.Unmarshal(res, &created)
	return created.Id, nil
}

func (c *Client) GetUser(id string) (*User, error) {
	body, err := c.httpRequest("GET", &strings.Reader{}, fmt.Sprintf("/%v", id))
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	user := &User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	return user, nil
}

func (c *Client) UpdateUser(id string, user *User) error {
	userInfo, err := json.Marshal(user)
	if err != nil {
		return err
	}
	body := strings.NewReader(string(userInfo))
	_, err = c.httpRequest("PATCH", body, fmt.Sprintf("/%v", id))
	if err != nil {
		log.Println("[UPDATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) DeleteUser(id, status string) error {
	var err error
	if status == "pending" {
		_, err = c.httpRequest("DELETE", &strings.Reader{}, fmt.Sprintf("/%s", id))
	} else {
		_, err = c.httpRequest("DELETE", &strings.Reader{}, fmt.Sprintf("/%s?action=delete", id))
	}
	if err != nil {
		log.Println("[DELETE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) ChangeUserStatus(id, action string) error {
	action = fmt.Sprintf("{\"action\":\"%s\"}", action)
	body := strings.NewReader(action)
	_, err := c.httpRequest("PUT", body, fmt.Sprintf("/%s/status", id))
	if err != nil {
		log.Println("[DEACTIVATE/ACTIVATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) ChangeEmail(oldEmail, newEmail string) error {
	body := strings.NewReader(fmt.Sprintf("{\"email\":\"" + newEmail + "\"}"))
	_, err := c.httpRequest("PUT", body, fmt.Sprintf("/%v/email", oldEmail))
	if err != nil {
		log.Println("[UPDATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) GenerateToken(accountId, clientId, clientSecret string) error {
	data := url.Values{}
	data.Set("grant_type", "account_credentials")
	data.Set("account_id", accountId)

	req, err := http.NewRequest(http.MethodPost, "https://zoom.us/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(string(respBody) + fmt.Sprintf(", StatusCode: %v", resp.StatusCode))
	}

	info := &AuthInfo{}
	err = json.Unmarshal(respBody, &info)
	if err != nil {
		return err
	}

	c.authToken = info.AccessToken
	return nil
}

func (c *Client) httpRequest(method string, body *strings.Reader, path string) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://api.zoom.us/v2/users"+path), body)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	authtoken := "Bearer " + c.authToken
	req.Header.Add("Authorization", authtoken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return respBody, nil
	}
	return nil, fmt.Errorf(string(respBody) + fmt.Sprintf(", StatusCode: %v", resp.StatusCode))
}

func (c *Client) IsRetry(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "429") == true {
			return true
		}
	}
	return false
}
