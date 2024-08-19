package qnap

// ContainerStationTaskResponse represents the response structure for Container Station tasks.
type ContainerStationTaskResponse struct {
	Data struct {
		TaskID string `json:"taskID"`
	} `json:"data"`
}

// ContainerStationOverview represents the response structure for Container Station overview.
type ContainerStationOverview struct {
	Data struct {
		App []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"app"`

		Container []struct {
			ID     string  `json:"id"`
			Name   string  `json:"name"`
			Type   string  `json:"type"`
			Status string  `json:"status"`
			CPU    float64 `json:"cpu"`
			Memory float64 `json:"memory"`
		} `json:"container"`
	} `json:"data"`
}

// Container represents the structure for a container.
type Container struct {
	ID                    string                   `json:"id"`
	Name                  string                   `json:"name"`
	Type                  string                   `json:"type"`
	Image                 string                   `json:"image"`
	ImageID               string                   `json:"imageid"`
	Status                string                   `json:"status"`
	Project               string                   `json:"project"`
	Runtime               string                   `json:"runtime"`
	MemLimit              int32                    `json:"memorylimit"`
	CpuLimit              int32                    `json:"cpulimit"`
	Cpupin                int32                    `json:"cpupin"`
	UUID                  string                   `json:"uuid"`
	UsedByInternalService string                   `json:"usedbyinternalservice"`
	Privileged            bool                     `json:"privileged"`
	CPU                   float32                  `json:"cpu"`
	Memory                float32                  `json:"memory"`
	TX                    int32                    `json:"tx"`
	RX                    int32                    `json:"rx"`
	Read                  int32                    `json:"read"`
	Write                 int32                    `json:"write"`
	Created               string                   `json:"created"`
	StartedAt             string                   `json:"startedat"`
	CMD                   []string                 `json:"cmd"`
	PortBindings          []containersPortBindings `json:"portbindings"`
	Networks              []containersNetworks     `json:"networks"`
}

// containersPortBindings represents the structure for port bindings of a container.
type containersPortBindings struct {
	Host        int32  `json:"host"`
	Container   int32  `json:"container"`
	Protocol    string `json:"protocol"`
	HostIP      string `json:"hostip"`
	ContainerIP string `json:"containerip"`
}

// containersNetworks represents the structure for networks of a container.
type containersNetworks struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	IpAddress   string `json:"ipaddress"`
	MacAddress  string `json:"macaddress"`
	Gateway     string `json:"gateway"`
	NetworkType string `json:"networktype"`
	IsStaticIP  bool   `json:"isstaticip"`
}

// NewContainerSpec represents the structure for creating a new container.
type NewContainerSpec struct {
	Type          string            `json:"type"`
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	AutoRemove    bool              `json:"autoremove"`
	Cmd           []string          `json:"cmd"`
	Entrypoint    []string          `json:"entrypoint"`
	Tty           bool              `json:"tty"`
	OpenStdin     bool              `json:"openstdin"`
	Pull          bool              `json:"pull"`
	Network       string            `json:"network"`
	NetworkType   string            `json:"networktype"`
	Hostname      string            `json:"hostname"`
	DNS           []string          `json:"dns"`
	Env           map[string]string `json:"env"`
	Labels        map[string]string `json:"labels"`
	Runtime       string            `json:"runtime"`
	Privileged    bool              `json:"privileged"`
	Operation     string            `json:"operation"`
	IPAddress     string            `json:"ipAddress"`
	Devices       []Devices         `json:"devices"`
	Volumes       []Volumes         `json:"volumes"`
	PortBindings  []PortBindings    `json:"portbindings"`
	Cpupin        Cpupin            `json:"cpupin"`
	RestartPolicy RestartPolicy     `json:"restartpolicy"`
}

// RestartPolicy represents the structure for the restart policy of a container.
type RestartPolicy struct {
	Name              string `json:"name"`
	MaximumRetryCount int32  `json:"maximumretrycount"`
}

// Cpupin represents the structure for CPU pinning of a container.
type Cpupin struct {
	CPUIDs string `json:"cpuids"`
	Type   string `json:"type"`
}

// PortBindings represents the structure for port bindings of a container.
type PortBindings struct {
	Host        int32  `json:"host"`
	Container   int32  `json:"container"`
	Protocol    string `json:"protocol"`
	HostIP      string `json:"hostip"`
	ContainerIP string `json:"containerip"`
}

// Volumes represents the structure for volumes of a container.
type Volumes struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Container   string `json:"container"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Permission  string `json:"permission"`
}

// Devices represents the structure for devices of a container.
type Devices struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

// ContainerInfo represents the response structure for container information.
type ContainerInfo struct {
	Data struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		Image          string `json:"image"`
		ImageID        string `json:"imageID"`
		Status         string `json:"status"`
		CPULimit       int32  `json:"cpuLimit"`
		MemLimit       int32  `json:"memLimit"`
		MemReservation int32  `json:"memReservation"`
		Cpupin         struct {
			CPUIDs string `json:"cpuids"`
			Type   string `json:"type"`
		} `json:"cpupin"`
		Networks []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			IPAddress   string `json:"ipAddress"`
			DisplayName string `json:"displayName"`
			MacAddress  string `json:"macAddress"`
			Gateway     string `json:"gateway"`
			NetworkType string `json:"networkType"`
			IsStaticIP  bool   `json:"isStaticIP"`
		} `json:"networks"`
		Project      string   `json:"project"`
		UUID         string   `json:"uuid"`
		Runtime      string   `json:"runtime"`
		Created      string   `json:"created"`
		StartedAt    string   `json:"startedAt"`
		FinishedAt   string   `json:"finishedAt"`
		Cmd          []string `json:"cmd"`
		DNS          []string `json:"dns"`
		ExposedPorts []string `json:"exposedPorts"`
		Pid          int32    `json:"pid"`
		PortBindings []struct {
			Host        int32  `json:"host"`
			Container   int32  `json:"container"`
			Protocol    string `json:"protocol"`
			HostIP      string `json:"hostIP"`
			ContainerIP string `json:"containerIP"`
		} `json:"portBindings"`
		Devices []struct {
			Name       string `json:"name"`
			Permission string `json:"permission"`
		} `json:"devices"`
		RestartPolicy struct {
			Name              string `json:"name"`
			MaximumRetryCount int32  `json:"maximumRetryCount"`
		} `json:"restartPolicy"`
		Entrypoint []string          `json:"entrypoint"`
		Privileged bool              `json:"privileged"`
		Env        map[string]string `json:"env"`
		Labels     map[string]string `json:"labels"`
		Volumes    []struct {
			Type        string `json:"type"`
			Name        string `json:"name"`
			Container   string `json:"container"`
			Source      string `json:"source"`
			Destination string `json:"destination"`
			Permission  string `json:"permission"`
		} `json:"volumes"`
		AutoRemove   bool    `json:"autoRemove"`
		Hostname     string  `json:"hostname"`
		CPU          float32 `json:"cpu"`
		Memory       float64 `json:"memory"`
		Tx           float32 `json:"tx"`
		Rx           float32 `json:"rx"`
		Read         float32 `json:"read"`
		Write        float32 `json:"write"`
		Tty          bool    `json:"tty"`
		OpenStdin    bool    `json:"openStdin"`
		DockerStatus struct {
			Running    bool   `json:"running"`
			Paused     bool   `json:"paused"`
			Restarting bool   `json:"restarting"`
			Dead       bool   `json:"dead"`
			ExitCode   int32  `json:"exitCode"`
			StartedAt  string `json:"startedAt"`
			FinishedAt string `json:"finishedAt"`
			Health     string `json:"health"`
		} `json:"dockerStatus"`
	} `json:"data"`
}

// NewAppReqModel represents the request structure for creating a new application.
type NewAppReqModel struct {
	LastUpdated    string                   `json:"last_updated"`
	Name           string                   `json:"name"`
	Yml            string                   `json:"yml"`
	DefaultURL     NewAppReqDefaultURLModel `json:"default_url"`
	CPULimit       int32                    `json:"cpu_limit"`
	MemLimit       int32                    `json:"mem_limit"`
	MemReservation int32                    `json:"mem_reservation"`
	Operation      string                   `json:"operation"`
}

// NewAppReqDefaultURLModel represents the default URL structure for a new application.
type NewAppReqDefaultURLModel struct {
	Port    int32  `json:"port"`
	Service string `json:"service"`
}

// AppRespModel represents the response structure for an application.
type AppRespModel struct {
	Data struct {
		Yml            string                   `json:"yml"`
		CPULimit       int32                    `json:"cpuLimit"`
		MemLimit       int32                    `json:"memLimit"`
		MemReservation int32                    `json:"memReservation"`
		DefaultURL     DefaultURLModel          `json:"defaultURL"`
		Containers     []AppRespContainersModel `json:"containers"`

		Status string `json:"status"`
	} `json:"data"`
}

// AppRespContainersModel represents the structure for containers in an application.
type AppRespContainersModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DefaultURLModel struct {
	Port      int32  `json:"port"`
	WebPort   int32  `json:"webPort"`
	URL       string `json:"url"`
	Container string `json:"container"`
	Protocol  string `json:"protocol"`
	Service   string `json:"service"`
}

// TaskModel represents the structure for a task.
type TaskModel struct {
	Data struct {
		Items []struct {
			ID          string `json:"id"`
			Category    string `json:"category"`
			Cancelable  bool   `json:"cancelable"`
			Description string `json:"description"`
			Progress    int    `json:"progress"`
			Detail      string `json:"detail"`
			State       string `json:"state"`
		} `json:"items"`
	} `json:"data"`
}

// Response from Inspect command
type VolumeRespModel struct {
	Data struct {
		VolumeInfo struct {
			CreatedAt  string `json:"CreatedAt"`
			Driver     string `json:"Driver"`
			Mountpoint string `json:"Mountpoint"`
			Name       string `json:"Name"`
			Scope      string `json:"Scope"`
		} `json:"volumeInfo"`
	} `json:"data"`
}

// Response from list volumes
type VolumesRespModel struct {
	Data struct {
		Items []struct {
			Created    string `json:"created"`
			Driver     string `json:"driver"`
			MountPoint string `json:"mountPoint"`
			Name       string `json:"name"`
			Project    string `json:"project"`
			Used       bool   `json:"used"`
			Size       int    `json:"size"`
		} `json:"items"`
	} `json:"data"`
}
