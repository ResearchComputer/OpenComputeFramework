package process

import (
	"bufio"
	"fmt"
	"os"
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

func readStuff(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println("Performed Scan")
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func (pm *ProcessManager) StartProcess(command string, envs string, args []string) {
	process := NewProcess(command, envs, args...)
	process = process.Start()
	// stream the output
	pm.processes = append(pm.processes, process)
}

func (pm *ProcessManager) StopAllProcesses() {
	for _, process := range pm.processes {
		process.Kill()
	}
}

func StartSubProcess(cmd string) {
	pm := NewProcessManager()
	pm.StartProcess(cmd, "", nil)
}
