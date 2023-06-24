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

type WorkloadTableRow struct {
	WorkloadID string   `json:"workload_id"`
	Providers  []string `json:"providers"`
}

type WorkloadTable struct {
	Workloads []WorkloadTableRow `json:"workloads"`
}

func (wt WorkloadTable) Add(workloadID string, provider string) {
	for _, workload := range wt.Workloads {
		if workload.WorkloadID == workloadID {
			workload.Providers = append(workload.Providers, provider)
			return
		}
	}
	row := WorkloadTableRow{WorkloadID: workloadID, Providers: []string{provider}}
	// todo(xiaozhe): figure out why linter complains about this
	wt.Workloads = append(wt.Workloads, row)
}

type NatsConnections struct {
	ServerID string `json:"server_id"`
}
