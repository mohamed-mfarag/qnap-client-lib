package qnap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CreateVolume - Create new volume
func (c *Client) CreateVolume(volumeName string, authToken *string) (*VolumeRespModel, error) {

	volume := fmt.Sprintf(`{"name":"%s"}`, volumeName)

	// rb, err := json.Marshal(volume)
	// if err != nil {
	// 	return nil, err
	// }
	rb := []byte(volume)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/container-station/api/v3/volumes", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	// Do request with volume name to inspect
	volumesBefore, err := c.ListVolumes(authToken)
	if err != nil {
		return nil, err
	}

	// Check if the returned data includes the requested name
	for _, volumeBefore := range volumesBefore.Data.Items {
		if volumeBefore.Name == volumeName {
			// Terminate if an volume with the same name already exists
			return nil, errors.New("can't create volume as an volume with the same name already exists")
		}
	}

	// If volume does not exist, then request creation
	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var postiveResponse struct {
		Data struct{} `json:"data"`
	}
	var negativeResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &postiveResponse)
	if err != nil {
		if err.Error() == "json: cannot unmarshal object into Go value of type qnap*" {
			err = json.Unmarshal(body, &negativeResponse)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(negativeResponse.Message)
		}
		return nil, err
	}

	// Do request with volume name to inspect
	volumesAfter, err := c.InspectVolume(volumeName, authToken)
	if err != nil {
		return nil, err
	}
	return volumesAfter, nil
}

// ListVolumes - List all volumes
func (c *Client) ListVolumes(authToken *string) (*VolumesRespModel, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/volumes", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var parsedData *VolumesRespModel
	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		return nil, err
	}

	return parsedData, nil
}

// InspectVolume - Inspect a volume
func (c *Client) InspectVolume(volumeName string, authToken *string) (*VolumeRespModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/container-station/api/v3/volumes/%s/inspect", c.HostURL, volumeName), nil)
	if err != nil {
		return nil, err
	}

	body, _, err := c.doRequest(req, authToken)
	if err != nil {
		return nil, err
	}

	var volumeResp *VolumeRespModel
	err = json.Unmarshal(body, &volumeResp)
	if err != nil {
		return nil, err
	}

	return volumeResp, nil
}

// DeleteVolume - Delete an volume
func (c *Client) DeleteVolume(volumeName string, authToken *string) (bool, error) {
	volumeToRemove := struct {
		Data struct {
			Items []struct {
				Name string `json:"name"`
			} `json:"items"`
		} `json:"data"`
	}{
		Data: struct {
			Items []struct {
				Name string `json:"name"`
			} `json:"items"`
		}{Items: []struct {
			Name string `json:"name"`
		}{{Name: volumeName}}}}

	// Marshal the payload to JSON
	rb, err := json.Marshal(volumeToRemove)
	if err != nil {
		return false, err
	}

	// Create a DELETE request to remove the volume
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/container-station/api/v3/volumes", c.HostURL), strings.NewReader(string(rb)))
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
	return true, nil
}
