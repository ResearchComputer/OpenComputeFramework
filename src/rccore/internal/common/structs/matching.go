package structs

type MatchingWorkerStatus struct {
	Accelerator string `json:"accelerator"`
	Status      string `json:"status"`
}

type ExpectedRuntime struct {
	WorkerID string  `json:"worker_id"`
	Runtime  float64 `json:"runtime"`
}

type MatchingModelStatus struct {
	Workers []ExpectedRuntime `json:"expectations"`
}

type MatchingStatus struct {
	Workers   map[string]MatchingWorkerStatus `json:"workers"`
	Models    map[string]MatchingModelStatus  `json:"models"`
	Timestamp int64                           `json:"timestamp"`
}
