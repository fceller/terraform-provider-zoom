package client

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestClient_GetUser(t *testing.T) {
	testCases := []struct {
		testName     string
		email        string
		expectErr    bool
		expectedResp *User
	}{
		{
			testName:  "user exists",
			email:     "user@gmail.com",
			expectErr: false,
			expectedResp: &User{
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
		},
		{
			testName:     "user does not exist",
			email:        "user@gmail.com",
			expectErr:    true,
			expectedResp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("", 2)
			accountId := os.Getenv("ZOOM_ACCOUNT_ID")
			clientId := os.Getenv("ZOOM_CLIENT_ID")
			clientSecret := os.Getenv("ZOOM_CLIENT_SECRET")
			client.GenerateToken(accountId, clientId, clientSecret)
			user, err := client.GetUser(tc.email)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, user)
		})
	}
}

func TestClient_NewItem(t *testing.T) {
	testCases := []struct {
		testName  string
		user      *User
		expectErr bool
	}{
		{
			testName: "user creation successful",
			user: &User{
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
			expectErr: false,
		},
		{
			testName: "user already exists",
			user: &User{
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("", 2)
			accountId := os.Getenv("ZOOM_ACCOUNT_ID")
			clientId := os.Getenv("ZOOM_CLIENT_ID")
			clientSecret := os.Getenv("ZOOM_CLIENT_SECRET")
			client.GenerateToken(accountId, clientId, clientSecret)
			_, err := client.NewUser(tc.user)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			user, err := client.GetUser(tc.user.Email)
			assert.NoError(t, err)
			assert.Equal(t, tc.user, user)
		})
	}
}

func TestClient_UpdateUser(t *testing.T) {
	testCases := []struct {
		testName  string
		user      *User
		expectErr bool
	}{
		{
			testName: "user exists",
			user: &User{
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
			expectErr: false,
		},
		{
			testName: "user does not exist",
			user: &User{
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				RoleName:   "Member",
				Status:     "active",
				Department: "devops",
				JobTitle:   "Engineer",
				Location:   "Delhi",
			},
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("", 2)
			accountId := os.Getenv("ZOOM_ACCOUNT_ID")
			clientId := os.Getenv("ZOOM_CLIENT_ID")
			clientSecret := os.Getenv("ZOOM_CLIENT_SECRET")
			client.GenerateToken(accountId, clientId, clientSecret)
			err := client.UpdateUser(tc.user.Email, tc.user)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			user, err := client.GetUser(tc.user.Email)
			assert.NoError(t, err)
			assert.Equal(t, tc.user, user)
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	testCases := []struct {
		testName  string
		email     string
		expectErr bool
	}{
		{
			testName:  "user exists",
			email:     "user@gmail.com",
			expectErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("", 2)
			accountId := os.Getenv("ZOOM_ACCOUNT_ID")
			clientId := os.Getenv("ZOOM_CLIENT_ID")
			clientSecret := os.Getenv("ZOOM_CLIENT_SECRET")
			client.GenerateToken(accountId, clientId, clientSecret)
			err := client.DeleteUser(tc.email, "pending")
			log.Println(err)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			_, err = client.GetUser(tc.email)
			assert.Error(t, err)
		})
	}
}
