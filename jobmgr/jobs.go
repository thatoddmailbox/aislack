package jobmgr

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type JobID int64
type JobStatusString string

const (
	JobStatusQueued    JobStatusString = "queued"
	JobStatusStarted   JobStatusString = "started"
	JobStatusCompleted JobStatusString = "completed"
	JobStatusFailed    JobStatusString = "failed"
)

type Job struct {
	ID         int               `json:"id"`
	Status     JobStatusString   `json:"status"`
	Priority   int               `json:"priority"`
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
	Created    int               `json:"created"`
	Started    *int              `json:"started"`
	Completed  *int              `json:"completed"`
	UserID     int               `json:"userID"`
}

type jobResponse struct {
	Status string  `json:"status"`
	Job    Job     `json:"job"`
	Result *string `json:"result"`
	// Artifacts []data.Artifact `json:"artifacts"`
}

func (c *Client) GetJob(jobID JobID) (Job, error) {
	var response jobResponse
	err := c.request(http.MethodGet, "jobs/get", url.Values{
		"id": []string{strconv.Itoa(int(jobID))},
	}, &response)
	if err != nil {
		return Job{}, err
	}

	return response.Job, nil
}

func (c *Client) StartJob(name string, parameters map[string]string) (JobID, error) {
	parametersData, err := json.Marshal(parameters)
	if err != nil {
		return 0, err
	}

	var response createdResponse
	err = c.request(http.MethodPost, "jobs/start", url.Values{
		"name":       []string{name},
		"parameters": []string{string(parametersData)},
	}, &response)
	if err != nil {
		return 0, err
	}

	return JobID(response.Created), nil
}
