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

func (c *Client) requestRaw(method string, path string, data url.Values) (*http.Response, io.ReadCloser, error) {
	fullURL := baseURL + path

	if method == http.MethodGet {
		if len(data) > 0 {
			fullURL += "?" + data.Encode()
		}
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, nil, err
	}

	if len(data) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		dataBytes := []byte(data.Encode())
		req.Body = io.NopCloser(bytes.NewReader(dataBytes))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp, resp.Body, nil
}

func (c *Client) request(method string, path string, data url.Values, result interface{}) error {
	_, responseBody, err := c.requestRaw(method, path, data)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	if result != nil {
		body, err := io.ReadAll(responseBody)
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
