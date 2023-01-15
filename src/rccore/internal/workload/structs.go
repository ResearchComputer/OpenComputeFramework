package workload

type Workload interface {
	start() bool
	stop() bool
}

type SingularityWorkload struct {
	BkmPath string
}

type DockerWorkload struct {
	BkmPath string
}

type RawcodeWorkload struct {
	VenvPath string
}
