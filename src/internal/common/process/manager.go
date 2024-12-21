package process

import (
	"strings"
)

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

func (pm *ProcessManager) StartProcess(command string, envs string, critical bool, args []string) {
	process := NewProcess(command, envs, critical, args...)
	process = process.Start()
	// stream the output
	pm.processes = append(pm.processes, process)
}

func (pm *ProcessManager) StopAllProcesses() {
	for _, process := range pm.processes {
		process.Kill()
	}
}

func StartCriticalProcess(cmd string) {
	pm := NewProcessManager()
	cmdParts := strings.Fields(cmd)
	if len(cmdParts) == 0 {
		return
	}
	command := cmdParts[0]
	args := cmdParts[1:]
	pm.StartProcess(command, "", true, args)
}

func HealthCheck() bool {
	pm := NewProcessManager()
	for _, process := range pm.processes {
		// check if the process is running
		if !process.isRunning() && process.critical {
			process.completed = true
			return false
		}
	}
	return true
}
