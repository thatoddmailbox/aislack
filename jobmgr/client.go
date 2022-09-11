package jobmgr

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "http://localhost:9845/"

type createdResponse struct {
	Status  string `json:"status"`
	Created int64  `json:"created"`
}

type Client struct {
}

func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) request(method string, path string, data url.Values, result interface{}) error {
	fullURL := baseURL + path

	if method == http.MethodGet {
		if len(data) > 0 {
			fullURL += "?" + data.Encode()
		}
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		dataBytes := []byte(data.Encode())
		req.Body = io.NopCloser(bytes.NewReader(dataBytes))
	}

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
