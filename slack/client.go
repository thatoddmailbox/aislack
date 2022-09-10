package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const baseURL = "https://slack.com/api/"

type Client struct {
	token       string
	commandLock sync.Mutex
	commandMap  map[string]CommandHandler
}

func NewClient(token string) (*Client, error) {
	return &Client{
		token: token,
	}, nil
}

func (c *Client) request(method string, path string, data interface{}) error {
	req, err := http.NewRequest(method, baseURL+path, nil)
	if err != nil {
		return err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		req.Body = io.NopCloser(bytes.NewReader(dataBytes))
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
