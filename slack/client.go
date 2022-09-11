package slack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		token:      token,
		commandMap: map[string]CommandHandler{},
	}, nil
}

func (c *Client) request(method string, path string, data interface{}, result interface{}) error {
	fullURL := baseURL + path

	// kinda janky
	if strings.HasPrefix(path, "https://") {
		fullURL = path
	}

	if data != nil {
		if method == http.MethodGet {
			fullURL += "?"
			values := url.Values{}
			for k, v := range data.(map[string]string) {
				values.Set(k, v)
			}
			fullURL += values.Encode()
		}
	}

	req, err := http.NewRequest(method, fullURL, nil)
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

	if result != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// log.Println(string(body))

		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
		// err = json.NewDecoder(resp.Body).Decode(&result)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}
