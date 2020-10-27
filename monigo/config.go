package monigo

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type WebConfig struct {
	Port int `json:"port"`
}

type DBConfig struct {
	File string `json:"file"`
}

type LogConfig struct {
	File           string `json:"file"`
	FileCritical   string `json:"file_critical"`
	Level          string `json:"level"`
	MaxSizeMB      int    `json:"max_size_mb"`
	MaxSizeBackups int    `json:"max_size_backups"`
	MaxAgeDays     int    `json:"max_age_days"`
}

type ResourceConfig struct {
	Frequency int64 `json:"frequency"`
}

type ResourcesConfig struct {
	OldDataLifeTime int64 `json:"old_data_life_time"`
	RAM     ResourceConfig `json:"ram"`
	CPU     ResourceConfig `json:"cpu"`
	Network ResourceConfig `json:"network"`

}

type Configuration struct {
	Web      WebConfig       `json:"web"`
	DB       DBConfig        `json:"db"`
	Log      LogConfig       `json:"log"`
	Resource ResourcesConfig `json:"resource"`
}

var Config *Configuration

func LoadConfig(path string) error {
	LogInfo("Loading configuration file")
	dt, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			bytes, _ := json.MarshalIndent(Configuration{
				Web: WebConfig{
					Port: 10000,
				},
				DB: DBConfig{
					File: "monigo.db",
				},
				Log: LogConfig{
					File:           "monigo.log",
					FileCritical:   "monigo.critical.log",
					Level:          "info",
					MaxSizeMB:      10,
					MaxSizeBackups: 5,
					MaxAgeDays:     7,
				},
				Resource: ResourcesConfig{
					OldDataLifeTime: 30,
					RAM: ResourceConfig{
						Frequency: 60,
					},
					CPU: ResourceConfig{
						Frequency: 60,
					},
					Network: ResourceConfig{
						Frequency: 60,
					},
				},
			}, "", "  ")

			_ = ioutil.WriteFile(path, bytes, os.ModePerm)
			LogInfo("Configuration file has empty. Creating")
		} else {
			LogCritical("Reading config %s failed with: %s", path, err)
			return err
		}
	}

	Config = &Configuration{}
	if err := json.Unmarshal(dt, Config); err != nil {
		LogCritical("Unmarshal config %s failed with: %s", path, err)
		return err
	}
	LogInfo("Configuration file was loaded")

	if Config.DB.File == "" {
		Config.DB.File = "monigo.db"
	}

	if Config.Resource.OldDataLifeTime == 0 {
		Config.Resource.OldDataLifeTime = 30
	}

	if Config.Resource.RAM.Frequency == 0 {
		Config.Resource.RAM.Frequency = 60
	}

	if Config.Resource.CPU.Frequency == 0 {
		Config.Resource.CPU.Frequency = 60
	}

	if Config.Resource.Network.Frequency == 0 {
		Config.Resource.Network.Frequency = 60
	}



	return nil
}
