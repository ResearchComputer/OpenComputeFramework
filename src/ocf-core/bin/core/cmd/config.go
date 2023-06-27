package cmd

type TOMP2PConfig struct {
	port string `json:"port", yaml:"port"`
}

type TOMVaccumConfig struct {
	Interval int `json:"interval", yaml:"interval"`
}

type QueueConfig struct {
	Port string `json:"port", yaml:"port"`
}

type TOMConfig struct {
	Path   string          `json:"path", yaml:"path"`
	Port   string          `json:"port", yaml:"port"`
	Name   string          `json:"name", yaml:"name"`
	P2p    TOMP2PConfig    `json:"p2p", yaml:"p2p"`
	Vacuum TOMVaccumConfig `json:"vacuum", yaml:"vacuum"`
	Queue  QueueConfig     `json:"queue", yaml:"queue"`
}

var defaultConfig = TOMConfig{
	Path:   "",
	Port:   "8092",
	Name:   "relay",
	P2p:    TOMP2PConfig{port: "8093"},
	Vacuum: TOMVaccumConfig{Interval: 10},
	Queue:  QueueConfig{Port: "8094"},
}
