package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	setDefaultConfig()
	setEnvVariables()
	setConfigFile()
	viper.Set("BinaryArch", "darwin_amd64") // TODO: multiarch
}

func setDefaultConfig() {
	viper.SetDefault("LogLevel", "debug")
	viper.SetDefault("LogFile", "~/.sonatina/sonatina.log")
	viper.SetDefault("DeploymentsPath", "~/.sonatina/deployments")
	viper.SetDefault("DeploymentsFilename", "deployments.json")
	viper.SetDefault("ManagerConnector", "json")
	viper.SetDefault("TestFilesystem", false)
	viper.SetDefault("TerraformPath", "~/.sonatina/terraform")
	viper.SetDefault("DefaultTerraformVersion", "0.12.24")
	viper.SetDefault("DefaultFlavour", "default")
	viper.SetDefault("Editor", "vi")
}

func setEnvVariables() {
	viper.SetEnvPrefix("SONATINA_")
}

func setConfigFile() {
	viper.SetConfigName("config")
	homeConfigPath, err := homedir.Expand("~/.sonatina")
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't expand homedir")
	}
	viper.AddConfigPath(homeConfigPath)
	viper.AddConfigPath("/etc/sonatina")
}
