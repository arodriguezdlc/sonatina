package main

import (
	"github.com/spf13/viper"
)

func init() {
	setDefaultConfig()
	setEnvVariables()
	setConfigFile()
}

func setDefaultConfig() {
	viper.SetDefault("LogLevel", "debug")
	viper.SetDefault("DeploymentsPath", "~/.sonatina/deployments")
	viper.SetDefault("DeploymentsFilename", "deployments.json")
	viper.SetDefault("ManagerConnector", "json")
	viper.SetDefault("TestFilesystem", false)
}

func setEnvVariables() {
	viper.SetEnvPrefix("SONATINA_")
}

func setConfigFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath("~/.sonatina")
	viper.AddConfigPath("/etc/sonatina")
}
