package slurm

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func IsSlurm() bool {
	return os.Getenv("SLURM_JOB_ID") != ""
}

func getJobId() string {
	return os.Getenv("SLURM_JOB_ID")
}

func getNodeId() string {
	return os.Getenv("SLURM_NODEID")
}

func getRemainingTimeInSeconds() int32 {
	/*
			needs to execute the following command:
			squeue -h -j $SLURM_JOB_ID -O TimeLeft | awk -F':|-' 'if (NF == 1) print $NF; \
		             else if (NF == 2) print ($1 * 60) + ($2); \
		             else if (NF == 3) print ($1 * 3600) + ($2 * 60) + $3; \
		             else if (NF == 4) print ($1 * 86400) + ($2 * 3600) + ($3 * 60) + $4'
	*/
	cmd := exec.Command("bash", "-c", "squeue -h -j $SLURM_JOB_ID -O TimeLeft | awk -F':|-' 'if (NF == 1) print $NF; else if (NF == 2) print ($1 * 60) + ($2); else if (NF == 3) print ($1 * 3600) + ($2 * 60) + $3; else if (NF == 4) print ($1 * 86400) + ($2 * 3600) + ($3 * 60) + $4'")
	output, err := cmd.Output()
	if err != nil {
		return -1
	}
	remainingTimeStr := strings.TrimSpace(string(output))
	remainingTime, err := strconv.Atoi(remainingTimeStr)
	if err != nil {
		return -1
	}
	return int32(remainingTime)
}

func GetJobInfo() map[string]any {
	info := make(map[string]any)
	info["job_id"] = getJobId()
	info["node_id"] = getNodeId()
	info["remaining_time"] = strconv.Itoa(int(getRemainingTimeInSeconds()))
	return info
}
