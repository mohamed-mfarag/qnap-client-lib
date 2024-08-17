# QNAP Package Documentation

This package provides functionalities to interact with the QNAP Container Station API. It allows you to manage containers, applications, and volumes on a QNAP device.
The package was build using reverse engineering for QNAP Container Station due to the unavaibility of an up to date API documenation, please consider this while using the client as code might change without notice after QNAP updates.

## Supported QNAP Version
QNAP NAS QuTS hero h5.1.8 | Container Station 3.0.7.891 (2024/05/09)


## Index

1. [Client](#client)
   - [NewClient](#newclient)
   - [doRequest](#dorequest)

2. [Authentication](#authentication)
   - [SignIn](#signin)
   - [SignOut](#signout)

3. [Container Management](#container-management)
   - [GetContainers](#getcontainers)
   - [CreateContainer](#createcontainer)
   - [InspectContainer](#inspectcontainer)
   - [DeleteContainer](#deletecontainer)
   - [StartContainer](#startcontainer)
   - [StopContainer](#stopcontainer)

4. [Application Management](#application-management)
   - [CreateApplication](#createapplication)
   - [InspectApplication](#inspectapplication)
   - [DeleteApplication](#deleteapplication)
   - [StartApplication](#startapplication)
   - [StopApplication](#stopapplication)

5. [Volume Management](#volume-management)
   - [CreateVolume](#createvolume)
   - [ListVolumes](#listvolumes)
   - [InspectVolume](#inspectvolume)
   - [DeleteVolume](#deletevolume)

6. [Miscellaneous](#miscellaneous)
   - [GetContainerStationOverview](#getcontainerstationoverview)
   - [GetTaskStatus](#gettaskstatus)

---

## 1. Client

The `Client` struct provides methods for interacting with the QNAP Container Station API.

### `NewClient`

```go
func NewClient(host, username, password *string) (*Client, error)
```

**Description:**

Creates a new QNAP client instance. Initializes the HTTP client with a timeout and sets the host URL. If `username` and `password` are provided, it attempts to sign in and retrieve an authentication token.

**Parameters:**
- `host`: Optional. The URL of the QNAP host. Defaults to `"http://localhost:19090"` if not provided.
- `username`: Optional. The username for authentication.
- `password`: Optional. The password for authentication.

**Returns:**
- `*Client`: The new client instance.
- `error`: An error if the client creation or authentication fails.

### `doRequest`

```go
func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, string, error)
```

**Description:**

Sends an HTTP request using the client's HTTP client. Handles the authentication token and returns the response body and token.

**Parameters:**
- `req`: The HTTP request to be sent.
- `authToken`: Optional. The authentication token to be included in the request header.

**Returns:**
- `[]byte`: The response body.
- `string`: The authentication token extracted from the response.
- `error`: An error if the request fails.

## 2. Authentication

Methods for handling authentication with the QNAP API.

### `SignIn`

```go
func (c *Client) SignIn() (*AuthResponse, error)
```

**Description:**

Obtains a new authentication token for the user. Requires the `Username` and `Password` fields to be set in the `Client`'s `Auth` field.

**Returns:**
- `*AuthResponse`: The authentication response containing the username and token.
- `error`: An error if authentication fails.

### `SignOut`

```go
func (c *Client) SignOut(authToken *string) error
```

**Description:**

Revokes the authentication token for the user.

**Parameters:**
- `authToken`: The authentication token to be revoked.

**Returns:**
- `error`: An error if the sign-out operation fails.

## 3. Container Management

Methods for managing containers on the QNAP Container Station.

### `GetContainers`

```go
func (c *Client) GetContainers() ([]Container, error)
```

**Description:**

Returns a list of all containers currently managed by Container Station.

**Returns:**
- `[]Container`: A list of containers.
- `error`: An error if the request fails.

### `CreateContainer`

```go
func (c *Client) CreateContainer(container NewContainerSpec, authToken *string) (*ContainerInfo, error)
```

**Description:**

Creates a new container. Ensures no container with the same name exists unless the operation is "recreate".

**Parameters:**
- `container`: The specifications for the new container.
- `authToken`: Optional. The authentication token.

**Returns:**
- `*ContainerInfo`: Information about the created container.
- `error`: An error if the creation fails.

### `InspectContainer`

```go
func (c *Client) InspectContainer(containerID string, containerType string, authToken *string) (*ContainerInfo, error)
```

**Description:**

Returns specific details about a container identified by its ID and type.

**Parameters:**
- `containerID`: The ID of the container.
- `containerType`: The type of the container.
- `authToken`: Optional. The authentication token.

**Returns:**
- `*ContainerInfo`: Details about the container.
- `error`: An error if the inspection fails.

### `DeleteContainer`

```go
func (c *Client) DeleteContainer(containerID string, containerType string, containerVolumeRemove bool, authToken *string) (bool, error)
```

**Description:**

Deletes a container. Optionally removes associated volumes.

**Parameters:**
- `containerID`: The ID of the container to delete.
- `containerType`: The type of the container.
- `containerVolumeRemove`: Whether to remove associated volumes.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the deletion is successful.
- `error`: An error if the deletion fails.

### `StartContainer`

```go
func (c *Client) StartContainer(containerID string, containerType string, authToken *string) (bool, error)
```

**Description:**

Starts a container identified by its ID and type.

**Parameters:**
- `containerID`: The ID of the container to start.
- `containerType`: The type of the container.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the container started successfully.
- `error`: An error if the start operation fails.

### `StopContainer`

```go
func (c *Client) StopContainer(containerID string, containerType string, authToken *string) (bool, error)
```

**Description:**

Stops a container identified by its ID and type.

**Parameters:**
- `containerID`: The ID of the container to stop.
- `containerType`: The type of the container.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the container stopped successfully.
- `error`: An error if the stop operation fails.

## 4. Application Management

Methods for managing applications on the QNAP Container Station.

### `CreateApplication`

```go
func (c *Client) CreateApplication(application NewAppReqModel, authToken *string) (*AppRespModel, error)
```

**Description:**

Creates a new application. Ensures no application with the same name exists.

**Parameters:**
- `application`: The specifications for the new application.
- `authToken`: Optional. The authentication token.

**Returns:**
- `*AppRespModel`: Information about the created application.
- `error`: An error if the creation fails.

### `InspectApplication`

```go
func (c *Client) InspectApplication(applicationName string, authToken *string) (*AppRespModel, error)
```

**Description:**

Returns specific details about an application identified by its name.

**Parameters:**
- `applicationName`: The name of the application.
- `authToken`: Optional. The authentication token.

**Returns:**
- `*AppRespModel`: Details about the application.
- `error`: An error if the inspection fails.

### `DeleteApplication`

```go
func (c *Client) DeleteApplication(applicationName string, containerVolumeRemove bool, authToken *string) (bool, error)
```

**Description:**

Deletes an application. Optionally removes associated volumes.

**Parameters:**
- `applicationName`: The name of the application to delete.
- `containerVolumeRemove`: Whether to remove associated volumes.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the deletion is successful.
- `error`: An error if the deletion fails.

### `StartApplication`

```go
func (c *Client) StartApplication(applicationName string, authToken *string) (bool, error)
```

**Description:**

Starts an application identified by its name.

**Parameters:**
- `applicationName`: The name of the application to start.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the application started successfully.
- `error`: An error if the start operation fails.

### `StopApplication`

```go
func (c *Client) StopApplication(applicationName string, authToken *string) (bool, error)
```

**Description:**

Stops an application identified by its name.

**Parameters:**
- `applicationName`: The name of the application to stop.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the application stopped successfully.
-

 `error`: An error if the stop operation fails.

## 5. Volume Management

Methods for managing volumes on the QNAP Container Station.

### `CreateVolume`

```go
func (c *Client) CreateVolume(volumeName string, authToken *string) (*VolumeRespModel, error)
```

**Description:**

Creates a new volume with the specified name.

**Parameters:**
- `volumeName`: The name of the new volume.
- `authToken`: Optional. The authentication token.

**Returns:**
- `*VolumeRespModel`: Information about the created volume.
- `error`: An error if the creation fails.

### `ListVolumes`

```go
func (c *Client) ListVolumes(authToken *string) (*VolumesRespModel, error)
```

**Description:**

Lists all volumes currently managed by Container Station.

**Parameters:**
- `authToken`: Optional. The authentication token.

**Returns:**
- `*VolumesRespModel`: A list of volumes.
- `error`: An error if the request fails.

### `InspectVolume`

```go
func (c *Client) InspectVolume(volumeName string, authToken *string) (*VolumeRespModel, error)
```

**Description:**

Returns specific details about a volume identified by its name.

**Parameters:**
- `volumeName`: The name of the volume to inspect.
- `authToken`: Optional. The authentication token.

**Returns:**
- `*VolumeRespModel`: Details about the volume.
- `error`: An error if the inspection fails.

### `DeleteVolume`

```go
func (c *Client) DeleteVolume(volumeName string, authToken *string) (bool, error)
```

**Description:**

Deletes a volume identified by its name.

**Parameters:**
- `volumeName`: The name of the volume to delete.
- `authToken`: Optional. The authentication token.

**Returns:**
- `bool`: `true` if the deletion is successful.
- `error`: An error if the deletion fails.

## 6. Miscellaneous

### `GetContainerStationOverview`

```go
func (c *Client) GetContainerStationOverview() (*ContainerStationOverview, error)
```

**Description:**

Returns an overview of all containers and applications running inside Container Station.

**Returns:**
- `*ContainerStationOverview`: Overview of containers and applications.
- `error`: An error if the request fails.

### `GetTaskStatus`

```go
func (c *Client) GetTaskStatus(taskID string) (string, error)
```

**Description:**

Gets the status of a task identified by its task ID.

**Parameters:**
- `taskID`: The ID of the task.

**Returns:**
- `string`: The status of the task (e.g., "completed").
- `error`: An error if the request fails.
