package jobmgr

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Artifact struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	MIME    string `json:"mime"`
	Size    int64  `json:"size"`
	UUID    string `json:"uuid"`
	Created int64  `json:"created"`
	JobID   int    `json:"jobID"`
}

func (c *Client) DownloadArtifact(artifact *Artifact) (io.ReadCloser, error) {
	resp, body, err := c.requestRaw(http.MethodGet, "artifacts/download", url.Values{
		"id": []string{strconv.Itoa(artifact.ID)},
	})
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// we're expecting a redirect
	// TODO: why does the go http client not follow this redirect already?
	location := resp.Header.Get("Location")
	fileResp, err := http.Get(location)
	if err != nil {
		return nil, err
	}

	return fileResp.Body, nil
}
