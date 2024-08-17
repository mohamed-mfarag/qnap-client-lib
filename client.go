package qnap

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:19090"

// Client represents the QNAP client
type Client struct {
	HostURL    string       // The URL of the QNAP host
	HTTPClient *http.Client // The HTTP client used for making requests
	Token      string       // The authentication token
	Auth       AuthStruct   // The authentication credentials
}

// AuthStruct represents the authentication credentials
type AuthStruct struct {
	Username string `json:"username"` // The username
	Password string `json:"password"` // The password
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Username string `json:"username"` // The username
	Token    string `json:"token"`    // The authentication token
}

// NewClient creates a new QNAP client
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL, // Set the default host URL
	}

	if host != nil {
		c.HostURL = *host // Override the default host URL if provided
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn() // Sign in and get the authentication response
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token // Set the authentication token

	return &c, nil
}

// doRequest sends an HTTP request and returns the response body, token, and error
func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, string, error) {
	var token string // Declare token variable

	if authToken != nil {
		token = *authToken
		req.Header.Set("Authorization", "Bearer "+strings.Split(token, "=")[1])
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", token)
	}

	res, err := c.HTTPClient.Do(req) // Send the HTTP request
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body) // Read the response body
	if err != nil {
		return nil, "", err
	}

	token = strings.Split(string(res.Header.Get("Set-Cookie")), ";")[0] // Extract the token from the response header

	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, token, err
}
