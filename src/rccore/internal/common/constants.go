package common

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func GetrccorePath() string {
	home, err := homedir.Dir()
	if err != nil {
		Logger.Error("Could not get home directory", "error", err)
		home = "."
	}
	rccorePath := path.Join(home, ".rccore")
	if _, err := os.Stat(rccorePath); os.IsNotExist(err) {
		err := os.MkdirAll(rccorePath, 0755)
		if err != nil {
			Logger.Error("Could not create rccore directory", "error", err)
			return "."
		}
	}
	return rccorePath
}

func GetDBPath() string {
	if viper.Get("database.path") != nil {
		return viper.GetString("database.path")
	}
	home, err := homedir.Dir()
	if err != nil {
		return "./rccore.db"
	}

	dbPath := path.Join(home, ".rccore", "rccore.db")
	return dbPath
}
