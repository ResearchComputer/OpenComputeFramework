package structs

type AvailableWorkload struct {
	Name  string   `json:"name"`
	Modes []string `json:"modes"`
}

type LoadWorkLoadInstruction struct {
	Workload        string            `json:"workload"`
	Mode            string            `json:"mode"`
	BootstrapConfig map[string]string `json:"bootstrap_config"`
}

type ProvisionModelsPlan struct {
	Instructions []LoadWorkLoadInstruction `json:"instructions"`
}

// WorkloadInstructionsHub maps workerID to LoadWorkLoadInstruction
type WorkloadInstructionsHub struct {
	Instructions map[string]ProvisionModelsPlan
}
