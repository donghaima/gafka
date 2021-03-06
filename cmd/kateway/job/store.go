// Package job implements the schedulable message(job) underlying storage.
package job

// JobStore is the backend storage layer for jobs(schedulable message).
type JobStore interface {

	// Name returns the underlying storage name.
	Name() string

	Start() error
	Stop()

	// CreateJobQueue creates a storage container where jobs will persist.
	CreateJobQueue(shardId int, appid, topic string) (err error)

	// Add pubs a schedulable message(job) synchronously.
	Add(appid, topic string, payload []byte, due int64) (jobId string, err error)

	// Delete removes a job by jobId.
	Delete(appid, topic, jobId string) (err error)
}

var Default JobStore
