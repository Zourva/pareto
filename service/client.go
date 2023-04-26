package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zourva/pareto/box/meta"
	"github.com/zourva/pareto/ipc"
)

// Registerer acts as a registry delegator,
// helping service instances
// interacting with service server.
// Registerer implements s1 interface.
type Registerer struct {
	//ref to the messager of the service
	messager *ipc.Messager

	//ref to the service this registerer serves.
	service Service

	exporter meta.Loop
}

// EnableStatusExport exports status of the service periodically.
func (r *Registerer) EnableStatusExport() {
	if r.exporter != nil {
		log.Warnln("already enabled status export")
		return
	}

	r.exporter = meta.NewLoop("registerer", meta.LoopConfig{
		Tick: r.service.Config().Interval * 1000,
	})

	// start the loop
	r.exporter.Run(meta.LoopRunHook{
		Working: func() error {
			err := r.messager.Publish(r.service.Config().Endpoint, r.service.Status())
			if err != nil {
				log.Warnf("export status for service %s failed: %v", r.service.Name(), err)
				return err
			}
			return nil
		},
	})
}

func (r *Registerer) DisableStatusExport() {
	r.exporter.Stop()
}

// Register registers the service to the registry server,
// and starts  a separate long-running loop to export service status
// periodically when the registration succeeded.
//
// Returns false when any error occurs, and true otherwise.
func (r *Registerer) Register(s Service) bool {
	r.service = s

	_, err := msgpack.Marshal(
		map[string]interface{}{
			"name": s.Name(),
		})
	if err != nil {
		return false
	}

	//r.messager.CallV2("/ew1/service/register", b, time.Second)

	return true
}

// Unregister unregisters the service from the registry server.
func (r *Registerer) Unregister() {
	//r.messager.CallV2("/ew1/service/deregister", []byte(s.Name()), time.Second)
}

// NewRegisterer creates a new registerer with the
// given messager as its communication channel.
func NewRegisterer(m *ipc.Messager) *Registerer {
	r := &Registerer{
		messager: m,
	}
	log.Infoln("a registerer is created")
	return r
}
