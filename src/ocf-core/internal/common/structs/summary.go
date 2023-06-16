package structs

type CardStatus struct {
	CardID          string  `json:"card_id"`
	Status          string  `json:"status"`
	Serving         string  `json:"serving"`
	PowerUsage      float64 `json:"power_usage"`
	GPUUtilization  float64 `json:"gpu_utilization"`
	UsedMemory      float64 `json:"used_memory"`
	AvailableMemory float64 `json:"available_memory"`
	LastUpdated     int64   `json:"last_updated"`
	GPUSpecifier    string  `json:"gpu_specifier"`
}

type StatusSummary struct {
	Status map[string]CardStatus `json:"status"`
}

type CardMetrics struct {
	PowerUsage      float64 `json:"power_usage"`
	GPUUtilization  float64 `json:"gpu_utilization"`
	UsedMemory      float64 `json:"used_memory"`
	AvailableMemory float64 `json:"available_memory"`
}
