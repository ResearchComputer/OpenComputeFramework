package common

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

func GetHomePath() string {
	home, err := homedir.Dir()
	if err != nil {
		Logger.Error("Could not get home directory", "error", err)
		home = "."
	}
	ocfcorePath := path.Join(home, ".ocfcore")
	if _, err := os.Stat(ocfcorePath); os.IsNotExist(err) {
		err := os.MkdirAll(ocfcorePath, 0755)
		if err != nil {
			Logger.Error("Could not create ocfcore directory", "error", err)
			return "."
		}
	}
	return ocfcorePath
}

func GetDBPath(nodeid string) string {
	return path.Join(GetHomePath(), "ocfcore."+nodeid+".db")
}
