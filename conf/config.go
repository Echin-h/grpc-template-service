package conf

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"grpc-template-service/pkg/colorful"
	"grpc-template-service/pkg/fsx"
	"os"
)

var serverConfig = &GlobalConfig{}

func LoadConfig(configPath ...string) {
	if len(configPath) == 0 || configPath[0] == "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./conf")
	} else {
		viper.SetConfigName(configPath[0])
	}

	loadConfig := func() {
		newConf := new(GlobalConfig)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("Config Read failed: " + err.Error())
			os.Exit(1)
		}
		err = viper.Unmarshal(newConf)
		if err != nil {
			fmt.Println("Config Unmarshal failed: " + err.Error())
			os.Exit(1)
		}
		serverConfig = newConf
	}

	loadConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config fileHandle changed: ", e.Name)
		loadConfig()
	})

	viper.WatchConfig()
}

// GenYamlConfig generate config file if not exist or force generate
func GenYamlConfig(path string, force bool) error {
	if !fsx.FileExist(path) || force {
		data, _ := yaml.Marshal(&GlobalConfig{MODE: "debug"})
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			return errors.New("Generate file with error: " + err.Error())
		}
		fmt.Println(colorful.Green("Config file `config.yaml` generate success in " + path))
	} else {
		return errors.New(path + " already exist, use -f to Force coverage")
	}
	return nil
}

func Get() *GlobalConfig {
	return serverConfig
}
