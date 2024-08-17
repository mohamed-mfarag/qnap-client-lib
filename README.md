# QNAP Container Station API Client

This package provides a Go client for interacting with the QNAP Container Station API. It allows you to manage containers, applications, and volumes on your QNAP device programmatically.

## Table of Contents

- [Installation](#installation)
- [Authentication](#authentication)
- [Container Management](#container-management)
- [Application Management](#application-management)
- [Volume Management](#volume-management)
- [Miscellaneous](#miscellaneous)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Installation

To use this package, you need to have Go installed on your system. You can install the package using `go get`:

```bash
go get github.com/mohamed-mfarag/qnap-client-lib
```

## Authentication

### `SignIn`

```go
func (c *Client) SignIn() (*AuthResponse, error)
```

Obtains a new authentication token. Requires the `Username` and `Password` fields to be set in the `Client`'s `Auth` field.

### `SignOut`

```go
func (c *Client) SignOut(authToken *string) error
```

Revokes the authentication token.

## Container Management

### `GetContainers`

```go
func (c *Client) GetContainers() ([]Container, error)
```

Returns a list of all containers.

### `CreateContainer`

```go
func (c *Client) CreateContainer(container NewContainerSpec, authToken *string) (*ContainerInfo, error)
```

Creates a new container. Ensures no container with the same name exists unless the operation is "recreate".

### `InspectContainer`

```go
func (c *Client) InspectContainer(containerID string, containerType string, authToken *string) (*ContainerInfo, error)
```

Returns specific details about a container identified by its ID and type.

### `DeleteContainer`

```go
func (c *Client) DeleteContainer(containerID string, containerType string, containerVolumeRemove bool, authToken *string) (bool, error)
```

Deletes a container. Optionally removes associated volumes.

### `StartContainer`

```go
func (c *Client) StartContainer(containerID string, containerType string, authToken *string) (bool, error)
```

Starts a container.

### `StopContainer`

```go
func (c *Client) StopContainer(containerID string, containerType string, authToken *string) (bool, error)
```

Stops a container.

## Application Management

### `CreateApplication`

```go
func (c *Client) CreateApplication(application NewAppReqModel, authToken *string) (*AppRespModel, error)
```

Creates a new application.

### `InspectApplication`

```go
func (c *Client) InspectApplication(applicationName string, authToken *string) (*AppRespModel, error)
```

Returns specific details about an application identified by its name.

### `DeleteApplication`

```go
func (c *Client) DeleteApplication(applicationName string, containerVolumeRemove bool, authToken *string) (bool, error)
```

Deletes an application. Optionally removes associated volumes.

### `StartApplication`

```go
func (c *Client) StartApplication(applicationName string, authToken *string) (bool, error)
```

Starts an application.

### `StopApplication`

```go
func (c *Client) StopApplication(applicationName string, authToken *string) (bool, error)
```

Stops an application.

## Volume Management

### `CreateVolume`

```go
func (c *Client) CreateVolume(volumeName string, authToken *string) (*VolumeRespModel, error)
```

Creates a new volume.

### `ListVolumes`

```go
func (c *Client) ListVolumes(authToken *string) (*VolumesRespModel, error)
```

Lists all volumes.

### `InspectVolume`

```go
func (c *Client) InspectVolume(volumeName string, authToken *string) (*VolumeRespModel, error)
```

Returns specific details about a volume identified by its name.

### `DeleteVolume`

```go
func (c *Client) DeleteVolume(volumeName string, authToken *string) (bool, error)
```

Deletes a volume.

## Miscellaneous

### `GetContainerStationOverview`

```go
func (c *Client) GetContainerStationOverview() (*ContainerStationOverview, error)
```

Returns an overview of all containers and applications.

### `GetTaskStatus`

```go
func (c *Client) GetTaskStatus(taskID string) (string, error)
```

Gets the status of a task identified by its task ID.

## Examples

Here is an example of how to use the QNAP client:

### Create and List Containers

```go
package main

import (
	"fmt"
)


func main() {

	host := "https://qnap.example.com"
	username := "username"
	password := "password"

	client, _ := NewClient(&host, &username, &password)
	
	authToken, _ := client.SignIn()

	// list all volumes
	volumes, err := client.ListVolumes(&authToken.Token)
	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, volume := range volumes.Data.Items {
		fmt.Println("Volume Name: ", volume.Name)
	}

	volumeName := "test-volume"
	// create volume
	volumeName = "test-volume"
	volume, err := client.CreateVolume(volumeName, &authToken.Token)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Volume has been created with name: ", volume.Data.VolumeInfo.Name)

	// delete volume
	_, err := client.DeleteVolume(volumeName, &authToken.Token)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Container Operations
	// Create the raw JSON data for creating the container
	containerData := NewContainerSpec{
		Name:        "demo",
		Type:        "docker",
		Image:       "ubuntu:latest",
		Network:     "bridge",
		NetworkType: "default",
		PortBindings: []PortBindings{
			{
				Host:        49123,
				Container:   8888,
				Protocol:    "TCP",
				HostIP:      "0.0.0.0",
				ContainerIP: "",
			},
		},
		RestartPolicy: RestartPolicy{
			Name:              "always",
			MaximumRetryCount: 0,
		},
		Volumes: []Volumes{
			{
				Type: "new",
			},
		},
	}
	// Create the container
	container, err := client.CreateContainer(containerData, &authToken.Token)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Container has been created and is in status " + container.Data.Status)


	// stop container
	client.StopContainer(container.Data.ID, container.Data.Type, &authToken.Token)

	//start container
	client.StartContainer(container.Data.ID, container.Data.Type, &authToken.Token)

	// Remove the container
	client.DeleteContainer(container.Data.ID, container.Data.Type, true, &authToken.Token)

	// Application Operations
	newApp := NewAppReqModel{
		Name: "postgresql-test",
		Yml:  "version: '3'\nservices:\n  postgres:\n    image: postgres:15.1\n    restart: unless-stopped\n    ports:\n      - 127.0.0.1:5432:5432\n    volumes:\n      - postgres_db:/var/lib/postgresql/data\n    environment:\n      POSTGRES_USER: postgres_qnap_user\n      POSTGRES_PASSWORD: postgres_qnap_pwd\n\n  phppgadmin:\n    image: qnapsystem/phppgadmin:7.13.0-1\n    restart: on-failure\n    ports:\n      - 7070:80\n    depends_on:\n      - postgres\n    environment:\n      PHP_PG_ADMIN_SERVER_HOST: postgres\n      PHP_PG_ADMIN_SERVER_PORT: 5432\n\nvolumes:\n  postgres_db:\n",
	}

	//Create the application
	_, err := client.CreateApplication(newApp, &authToken.Token)
	if err != nil {
		fmt.Println("Error:", err)
	}

	//Stop the application
	client.StopApplication("postgresql-test", &authToken.Token)
	time.Sleep(10 * time.Second)

	//Start the application
	client.StartApplication("postgresql-test", &authToken.Token)
	time.Sleep(10 * time.Second)

	//Remove the application
	client.DeleteApplication("postgresql-test", true, &authToken.Token)

}
```

## Contributing

Contributions are welcome! Please submit pull requests or open issues on the GitHub repository.

## License

This package is licensed under the Apache License. See the [LICENSE](LICENSE) file for details.
