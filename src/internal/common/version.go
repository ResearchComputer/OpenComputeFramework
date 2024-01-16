package common

type jsonVersion struct {
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	Date        string `json:"date"`
	BuildSecret string `json:"buildSecret"`
}

var JSONVersion jsonVersion
