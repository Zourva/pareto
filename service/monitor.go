package service

import (
	"sync"
)

type StatusConf struct {
	//Interval, in seconds, to refresh and publish service status,
	//optional with a minimum value of 1 second.
	//If not provided, or set to 0, status publishing is disabled.
	Interval uint32 `json:"interval"`

	//Threshold number of intervals before out-of-sync
	Threshold uint32 `json:"threshold"`

	//Endpoint used to export service status periodically, optional.
	//If not provided, the default format is used: {service name}/status.
	//Changed: use centered topic to aggregate service status.
	//endpoint string
}

// Status defines heartbeat info published by a service.
type Status struct {
	Name   string      `json:"name"`             //name of the service
	State  State       `json:"state"`            //state of the service
	Time   uint64      `json:"time"`             //report timestamp in milliseconds
	Health *StatusConf `json:"health,omitempty"` //if provided, check timeout
}

// StatusList defines all services status info.
type StatusList struct {
	Services []*Status `json:"services"`
}

func getDefaultStatusConf() *StatusConf {
	return &StatusConf{
		Interval:  StatusReportInterval,
		Threshold: StatusLostThreshold,
		//endpoint:  EndpointServiceStatus,
	}
}

//func getDefaultStatusEndpoint(name string) string {
//	return fmt.Sprintf("%s/status", name)
//}

// Monitor monitors status of services
// and that of some infrastructure, and
// report alerts if required.
type Monitor struct {
	registry *RegistryManager
}

func (m *Monitor) GetStatus(name string) *Status {
	reg := m.registry.get(name)
	if reg == nil {
		return nil
	}

	return &Status{
		Name:  reg.name,
		State: reg.state,
		Time:  reg.updateTime,
	}
}

func (m *Monitor) GetStatusList() StatusList {
	var list StatusList
	m.registry.services.Range(func(key, value any) bool {
		reg := value.(*registry)

		list.Services = append(list.Services, &Status{
			Name:  reg.name,
			State: reg.state,
			Time:  reg.updateTime,
			//Health: &StatusConf{
			//	Interval:  uint32(reg.interval),
			//	Threshold: uint32(reg.threshold),
			//},
		})

		return true
	})

	return list
}

// GetNotServicing returns names of services from the given names
// that are not ready yet.
func (m *Monitor) GetNotServicing(filtered []string) []string {
	var result []string
	for _, name := range filtered {
		s := m.registry.get(name)
		if s == nil || s.state != Servicing {
			result = append(result, name)
		}
	}

	return result
}

var monLock sync.Mutex
var monitor *Monitor

// NewMonitor creates an instance of Monitor.
// Monitor itself is not a service, however,
// it needs to inquiry service registry to
// get service status, so it creates a service registry
// internally to manage all services registered.
func NewMonitor(registry string) *Monitor {
	s := &Monitor{
		registry: NewRegistryManager(registry),
	}

	return s
}

func GetMonitor() *Monitor {
	return monitor
}

// EnableMonitor enables service monitor by creating
// and attaching a service registry manager to the
// given service registry address.
func EnableMonitor(registry string) *Monitor {
	monLock.Lock()
	defer monLock.Unlock()

	if monitor == nil {
		monitor = NewMonitor(registry)
	}

	Start(monitor.registry)

	return monitor
}

// DisableMonitor disables the service monitor
// if it is enabled already.
func DisableMonitor() {
	if monitor == nil {
		return
	}

	monLock.Lock()
	defer monLock.Unlock()

	Stop(monitor.registry)

	monitor = nil
}
