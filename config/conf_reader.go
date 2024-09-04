package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/mgolfam/gogutils/filemanager"
	"github.com/mgolfam/gogutils/glog"
)

func LoadConfig(configPath, logLevel string, configStruct interface{}) error {
	confStr, err := filemanager.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(confStr), &configStruct); err != nil {
		log.Fatal(err)
	}

	glog.Log(configStruct)

	// configure log leve
	glog.LogLevel.Label = strings.ToUpper(logLevel)
	glog.LogLevel.Load()

	return nil
}
