package qnap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Data struct {
	Containers []Container `json:"items"`
}
type RemoveContainer struct {
	Data RemoveContainerData `json:"data"`
}

type RemoveContainerData struct {
	Items         []Item `json:"items"`
	RemoveVolumes bool   `json:"removeVolumes"`
}

type ChangeContainer struct {
	Data ChangeContainerData `json:"data"`
}

type ChangeContainerData struct {
	Items []Item `json:"items"`
}

type Item struct {
	CID   string `json:"cid"`
	CType string `json:"ctype"`
}

var TaskStatusCompleted = "completed"
var ContainerStatusRunning = "running"

// GetContainers returns a list of containers
func (c *Client) GetContainers() ([]Container, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/containers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var parsedData struct {
		Data Data `json:"data"`
	}

	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		return nil, err
	}

	return parsedData.Data.Containers, nil
}

// CreateContainer creates a new container
func (c *Client) CreateContainer(container NewContainerSpec, authToken *string) (*ContainerInfo, error) {
	containerName := container.Name
	containerOperation := container.Operation
	TaskStatusCompleted := "completed"

	rb, err := json.Marshal(container)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/container-station/api/v3/containers", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	containersBefore, err := c.GetContainerStationOverview()
	if err != nil {
		return nil, err
	}

	for _, containerBefore := range containersBefore.Data.Container {
		if containerBefore.Name == containerName && containerOperation != "recreate" {
			return nil, errors.New("cannot create container as a container with the same name already exists")
		}
	}

	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var response ContainerStationTaskResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for {
		taskStatus, err := c.GetTaskStatus(response.Data.TaskID)
		if err != nil {
			return nil, err
		}

		if taskStatus == TaskStatusCompleted {
			break
		}

		time.Sleep(2 * time.Second)
	}

	containersAfter, err := c.GetContainerStationOverview()
	if err != nil {
		return nil, err
	}

	var newContainerInfo *ContainerInfo
	for _, containerAfter := range containersAfter.Data.Container {
		if containerAfter.Name == containerName {
			newContainerInfo, err = c.InspectContainer(containerAfter.ID, containerAfter.Type, authToken)
			if err != nil {
				return nil, err
			}
			return newContainerInfo, nil
		}
	}

	return nil, errors.New("container is not found after creation. Possible options: QNAP container station needs more time or the container creation failed silently")
}

// InspectContainer returns specific container specifications
func (c *Client) InspectContainer(containerID string, containerType string, authToken *string) (*ContainerInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/containers/%s?id=%s",
		c.HostURL, containerType, containerID), nil)
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var containerData ContainerInfo
	err = json.Unmarshal(body, &containerData)
	if err != nil {
		return nil, err
	}

	return &containerData, nil
}

// DeleteContainer deletes a container
func (c *Client) DeleteContainer(containerID string, containerType string, containerVolumeRemove bool, authToken *string) (bool, error) {
	return c.ChangeContainerState(containerID, containerType, containerVolumeRemove, "delete", authToken)
}

// StartContainer starts a container
func (c *Client) StartContainer(containerID string, containerType string, authToken *string) (bool, error) {
	return c.ChangeContainerState(containerID, containerType, false, "start", authToken)
}

// StopContainer stops a container
func (c *Client) StopContainer(containerID string, containerType string, authToken *string) (bool, error) {
	return c.ChangeContainerState(containerID, containerType, false, "stop", authToken)
}

// ChangeContainerState changes the state of a container - used by start,stop and delete functions
func (c *Client) ChangeContainerState(containerID string, containerType string, containerVolumeRemove bool, operation string, authToken *string) (bool, error) {
	var httpOperation string
	var rb []byte
	var err error
	var url string

	if operation == "start" || operation == "stop" {
		container := ChangeContainer{
			Data: ChangeContainerData{
				Items: []Item{
					{
						CID:   containerID,
						CType: containerType,
					},
				},
			},
		}
		rb, err = json.Marshal(container)
		if err != nil {
			return false, err
		}
		httpOperation = "PUT"
		url = fmt.Sprintf("%s/container-station/api/v3/containers/%s", c.HostURL, operation)
	} else if operation == "delete" {
		container := RemoveContainer{
			Data: RemoveContainerData{
				Items: []Item{
					{
						CID:   containerID,
						CType: containerType,
					},
				},
				RemoveVolumes: containerVolumeRemove,
			},
		}
		rb, err = json.Marshal(container)
		if err != nil {
			return false, err
		}
		httpOperation = "DELETE"
		url = fmt.Sprintf("%s/container-station/api/v3/containers", c.HostURL)
	} else {
		return false, errors.New("container operation " + operation + " not supported")
	}

	req, err := http.NewRequest(httpOperation, url, strings.NewReader(string(rb)))
	if err != nil {
		return false, err
	}

	body, _, err := c.doRequest(req, authToken)
	if err != nil {
		return false, err
	}

	var response ContainerStationTaskResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	for {
		taskStatus, err := c.GetTaskStatus(response.Data.TaskID)
		if err != nil {
			return false, err
		}

		if taskStatus == TaskStatusCompleted {
			break
		}

		time.Sleep(2 * time.Second)
	}

	containersAfter, err := c.GetContainerStationOverview()
	if err != nil {
		return false, err
	}

	for _, containerAfterItem := range containersAfter.Data.Container {
		if containerAfterItem.ID == containerID {
			switch operation {
			case "start":
				if containerAfterItem.Status == "running" {
					return true, nil
				} else {
					return false, errors.New("container operation " + operation + " failed to complete")
				}
			case "stop":
				if containerAfterItem.Status == "stopped" {
					return true, nil
				} else {
					return false, errors.New("container operation " + operation + " failed to complete")
				}
			}
		}
	}
	if operation == "delete" {
		return true, nil
	} else {
		return false, errors.New("container operation " + operation + " failed to complete, container not found..")
	}
}
