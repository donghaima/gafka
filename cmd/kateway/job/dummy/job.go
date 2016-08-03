package dummy

import (
	"time"

	"github.com/funkygao/gafka/cmd/kateway/job"
)

type dummy struct{}

func New() job.JobStore {
	return &dummy{}
}

func (this *dummy) Add(cluster, topic string, payload []byte, delay time.Duration) (jobId int64, err error) {
	return
}

func (this *dummy) Delete(cluster, jobId int64) (err error) {
	return
}

func (this *dummy) Name() string {
	return "dummy"
}

func (this *dummy) Start() error {
	return nil
}

func (this *dummy) Stop() {}
