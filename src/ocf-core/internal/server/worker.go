package server

import (
	"ocfcore/internal/common"
	"ocfcore/internal/common/structs"
	"ocfcore/internal/profiler"
	"sync"

	"github.com/gofrs/uuid"
)

type Worker struct {
	WorkerID          string                      `json:"worker_id"`
	WorkerIP          string                      `json:"worker_ip"`
	GPUSpecifier      string                      `json:"gpu_specifier"`
	GPUMemory         float32                     `json:"gpu_memory"`
	AvailableWorkload []structs.AvailableWorkload `json:"available_workload"`
	Serving           string                      `json:"serving"`
	LastUpdated       int64                       `json:"last_updated"`
}

type WorkerService struct{}

type WorkerHub struct {
	Workers []Worker
}

type WorkerGroup struct {
	Workers []Worker
	Model   string `json:"model"`
	Busy    bool   `json:"busy"`
}

type WorkerStatusResponse struct {
	WorkerID  string  `json:"worker_id"`
	Metric    string  `json:"metric"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

var workerHub WorkerHub
var workloadTable *structs.WorkloadTable
var once sync.Once

func GlobalWorkloadTable() *structs.WorkloadTable {
	once.Do(func() {
		workloadTable = &structs.WorkloadTable{}
	})
	return workloadTable
}

func (wh WorkerHub) Exists(workerID string) bool {
	for _, worker := range wh.Workers {
		if worker.WorkerID == workerID {
			return true
		}
	}
	return false
}

func (s *WorkerService) Join(WorkerIP string, GPUSpecifier string, GPUMemory float32, availableWorkload []structs.AvailableWorkload) string {
	WorkerID := uuid.Must(uuid.NewV7()).String()
	common.Logger.Infof("Worker %s joined with GPU: %s", WorkerID, GPUSpecifier)
	workerHub.Workers = append(workerHub.Workers, Worker{WorkerID, WorkerIP, GPUSpecifier, GPUMemory, availableWorkload, "", 0})
	return WorkerID
}

func (s *WorkerService) Rejoin(WorkerID string, WorkerIP string, GPUSpecifier string, GPUMemory float32, availableWorkload []structs.AvailableWorkload) string {
	common.Logger.Infof("Worker %s rejoined with GPU: %s", WorkerID, GPUSpecifier)
	workerHub.Workers = append(workerHub.Workers, Worker{WorkerID, WorkerIP, GPUSpecifier, GPUMemory, availableWorkload, "", 0})
	return WorkerID
}

func (s *WorkerService) Leave(WorkerID string) (int, error) {
	return 1, nil
}

// the returned value instructs the client if they should re-join
func (s *WorkerService) Update(Timestamp int64, Metric string, Value float64, Id string) int {
	common.Logger.Infof("Worker %s updated with %s: %f", Id, Metric, Value)
	// if the worker is not in the hub, return 1 to instruct the client to re-join
	profiler.AddPoint(Id, Metric, Timestamp, Value)
	// update last updated time for the worker
	for idx, worker := range workerHub.Workers {
		if worker.WorkerID == Id {
			workerHub.Workers[idx].LastUpdated = Timestamp
		}
	}
	if !workerHub.Exists(Id) {
		return 1
	}
	return 0
}

func (s *WorkerService) GetDesiredWorkload(WorkerID string) structs.ProvisionModelsPlan {
	if instructionsHub.Instructions == nil {
		return structs.ProvisionModelsPlan{}
	} else {
		// this is one-off, after loading the workload, the instruction should be deleted
		desiredWorkload := instructionsHub.Instructions[WorkerID]
		instructionsHub.Instructions[WorkerID] = structs.ProvisionModelsPlan{}
		return desiredWorkload
	}
}

func (s *WorkerService) UpdateServingStatus(WorkerID string, Serving string) {
	for i, worker := range workerHub.Workers {
		if worker.WorkerID == WorkerID {
			workerHub.Workers[i].Serving = Serving
		}
	}
}
