package qnap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var data struct {
	ContainerStationOverview
}

// GetContainerStationOverview  - Returns all containers and apps running inside container station
func (c *Client) GetContainerStationOverview() (*ContainerStationOverview, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/overview", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data.ContainerStationOverview, nil
}

// Get task status
func (c *Client) GetTaskStatus(taskID string) (string, error) {
	var data TaskModel
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/tasks", c.HostURL), nil)
	if err != nil {
		return "unknown", err
	}

	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return "unknown", err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "unknown", err
	}
	for _, task := range data.Data.Items {
		if task.ID == taskID {
			return task.State, nil
		}
	}
	return "not-found", nil
}
