// Package wesender is de officiële Go SDK voor de WeSender e-mail API.
package wesender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultBaseURL = "https://api.wesender.nl"

type WesenderError struct {
	Message string
	Status  int
}

func (e *WesenderError) Error() string {
	return fmt.Sprintf("wesender: %s (HTTP %d)", e.Message, e.Status)
}

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey, baseURL: defaultBaseURL, http: &http.Client{}}
}

type SendEmailInput struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

type EmailResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type Domain struct {
	ID       string `json:"id"`
	Domain   string `json:"domain"`
	SpfOk   bool   `json:"spfOk"`
	DkimOk  bool   `json:"dkimOk"`
	DmarcOk bool   `json:"dmarcOk"`
}

func (c *Client) SendEmail(input SendEmailInput) (*EmailResult, error) {
	var result EmailResult
	return &result, c.post("/emails", input, &result)
}

func (c *Client) ListDomains() ([]Domain, error) {
	var resp struct {
		Data []Domain `json:"data"`
	}
	return resp.Data, c.get("/domains", &resp)
}

func (c *Client) CreateDomain(domain string) (*Domain, error) {
	var result Domain
	return &result, c.post("/domains", map[string]string{"domain": domain}, &result)
}

func (c *Client) get(path string, out interface{}) error {
	return c.do("GET", path, nil, out)
}

func (c *Client) post(path string, body interface{}, out interface{}) error {
	return c.do("POST", path, body, out)
}

func (c *Client) do(method, path string, body interface{}, out interface{}) error {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}
	req, err := http.NewRequest(method, c.baseURL+path, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "wesender-go/1.0.0")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var errBody struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		return &WesenderError{Message: errBody.Error, Status: res.StatusCode}
	}
	return json.NewDecoder(res.Body).Decode(out)
}
