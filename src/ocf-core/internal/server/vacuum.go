package server

import (
	"ocfcore/internal/common"
	"time"
)

func DisconnectionDetection(tolerance time.Duration) {
	new_workers := workerHub.Workers[:0]
	// if a worker in the workerhub does not update its status for a certain amount of time, remove it from the workerhub
	for idx, worker := range workerHub.Workers {
		if worker.LastUpdated >= time.Now().Unix()-int64(tolerance.Seconds()) {
			new_workers = append(new_workers, workerHub.Workers[idx])
		} else {
			common.Logger.Debug("Worker " + worker.WorkerID + " disconnected")
		}
	}
	workerHub.Workers = new_workers
}
