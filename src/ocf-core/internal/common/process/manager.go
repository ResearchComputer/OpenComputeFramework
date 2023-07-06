package process

type ProcessManager struct {
	processes []*Process
}

var pm *ProcessManager

func NewProcessManager() *ProcessManager {
	if pm == nil {
		pm = &ProcessManager{}
	}
	return pm
}

func (pm *ProcessManager) StartProcess(command string, envs string, args []string) {
	process := NewProcess(command, envs, args...)
	process = process.Start()
	pm.processes = append(pm.processes, process)
}

func (pm *ProcessManager) StopAllProcesses() {
	for _, process := range pm.processes {
		process.Kill()
	}
}
