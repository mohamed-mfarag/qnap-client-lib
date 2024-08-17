package qnap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RemoveApplication struct {
	Apps         []string `json:"apps"`
	RemoveVolume bool     `json:"removeVolume"`
}
type ChangeApplication struct {
	Apps []string `json:"apps"`
}

// CreateApplication - Create new application
func (c *Client) CreateApplication(application NewAppReqModel, authToken *string) (*AppRespModel, error) {
	applicationName := application.Name
	applicationOperation := application.Operation

	rb, err := json.Marshal(application)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/container-station/api/v3/apps/compose", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	// Do request with application name to inspect
	applicationsBefore, err := c.GetContainerStationOverview()
	if err != nil {
		return nil, err
	}

	// Check if the returned data includes the requested name
	for _, applicationBefore := range applicationsBefore.Data.App {
		if applicationBefore.Name == applicationName && applicationOperation != "recreate" {
			// Terminate if an application with the same name already exists
			return nil, errors.New("can't create application as an application with the same name already exists")
		}
	}

	// If application does not exist, then request creation
	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var response ContainerStationTaskResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Loop until task is in status completed
	for {
		// Get task status
		taskStatus, err := c.GetTaskStatus(response.Data.TaskID)
		if err != nil {
			return nil, err
		}

		// Check if task is completed
		if taskStatus == "completed" {
			break
		}

		// Sleep for a while before checking again
		time.Sleep(2 * time.Second)
	}

	// Do request with application name to inspect
	applicationsAfter, err := c.GetContainerStationOverview()
	if err != nil {
		return nil, err
	}

	// Check if application is created
	for _, applicationAfter := range applicationsAfter.Data.App {
		if applicationAfter.Name == applicationName {
			newApplication, err := c.InspectApplication(applicationAfter.Name, authToken)
			if err != nil {
				return nil, err
			}
			newApplication.Data.Status = applicationAfter.Status
			return newApplication, nil
		}
	}

	return nil, errors.New("application is not found after creation. Possible options, QNAP application station needs more time or the application creation failed silently")
}

// InspectApplication - Returns specific container specifications (not inspect function)
func (c *Client) InspectApplication(applicationName string, authToken *string) (*AppRespModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/apps/%s/inspect", c.HostURL, applicationName), nil)
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, authToken)
	if err != nil {
		return nil, err
	}

	var applicationResp AppRespModel
	err = json.Unmarshal(body, &applicationResp)
	if err != nil {
		return nil, err
	}

	// Do request with application name to inspect
	applicationsAfter, err := c.GetContainerStationOverview()
	if err != nil {
		return nil, err
	}

	// Check if application is created
	for _, applicationAfter := range applicationsAfter.Data.App {
		if applicationAfter.Name == applicationName {
			applicationResp.Data.Status = applicationAfter.Status
		}
	}

	return &applicationResp, nil
}

// DeleteApplication - Delete an application
func (c *Client) DeleteApplication(applicationName string, containerVolumeRemove bool, authToken *string) (bool, error) {
	return c.ChangeApplicationState(applicationName, containerVolumeRemove, "delete", authToken)
}

func (c *Client) StartApplication(applicationName string, authToken *string) (bool, error) {
	return c.ChangeApplicationState(applicationName, false, "start", authToken)
}

func (c *Client) StopApplication(applicationName string, authToken *string) (bool, error) {
	return c.ChangeApplicationState(applicationName, false, "stop", authToken)
}

func (c *Client) ChangeApplicationState(applicationName string, containerVolumeRemove bool, operation string, authToken *string) (bool, error) {
	var httpOperation string
	var rb []byte
	var err error
	var url string
	if operation == "start" || operation == "stop" {
		applicationToChange := ChangeApplication{
			Apps: []string{applicationName},
		}
		// Marshal the payload to JSON
		rb, err = json.Marshal(applicationToChange)
		if err != nil {
			return false, err
		}
		httpOperation = "PUT"
		url = fmt.Sprintf("%s/container-station/api/v3/apps/%s", c.HostURL, operation)
	} else if operation == "delete" {
		applicationToRemove := RemoveApplication{
			Apps:         []string{applicationName},
			RemoveVolume: containerVolumeRemove,
		}
		// Marshal the payload to JSON
		rb, err = json.Marshal(applicationToRemove)
		if err != nil {
			return false, err
		}
		httpOperation = "DELETE"
		url = fmt.Sprintf("%s/container-station/api/v3/apps", c.HostURL)
	} else {
		return false, errors.New("container operation " + operation + " not supported")
	}

	// Create a DELETE request to remove the application
	req, err := http.NewRequest(httpOperation, url, strings.NewReader(string(rb)))
	if err != nil {
		return false, err
	}

	// Send the request and get the response body
	body, _, err := c.doRequest(req, authToken)
	if err != nil {
		return false, err
	}

	// Unmarshal the response body to get the task ID
	var response ContainerStationTaskResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	// Wait until the task is completed
	for {
		// Get the task status
		taskStatus, err := c.GetTaskStatus(response.Data.TaskID)
		if err != nil {
			return false, err
		}

		// Check if the task is completed
		if taskStatus == "completed" {
			break
		}

		// Sleep for a while before checking again
		time.Sleep(2 * time.Second)
	}

	// Get the updated list of applications after deletion
	applicationsAfter, err := c.GetContainerStationOverview()
	if err != nil {
		return false, err
	}

	// Check if the application is still running
	for _, applicationAfterItem := range applicationsAfter.Data.App {
		if applicationAfterItem.Name == applicationName {
			switch operation {
			case "start":
				if applicationAfterItem.Status == "running" {
					return true, nil
				} else {
					return false, errors.New("application operation " + operation + " failed to complete")
				}
			case "stop":
				if applicationAfterItem.Status == "stopped" {
					return true, nil
				} else {
					return false, errors.New("application operation " + operation + " failed to complete")
				}
			}
		}
	}
	if operation == "delete" {
		return true, nil
	} else {
		return false, errors.New("application operation " + operation + " failed to complete, application not found..")
	}
}
