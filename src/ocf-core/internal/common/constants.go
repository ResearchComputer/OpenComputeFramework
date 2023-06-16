package common

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func GetocfcorePath() string {
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

func GetDBPath() string {
	if viper.Get("database.path") != nil {
		return viper.GetString("database.path")
	}
	home, err := homedir.Dir()
	if err != nil {
		return "./ocfcore.db"
	}

	dbPath := path.Join(home, ".ocfcore", "ocfcore.db")
	return dbPath
}
