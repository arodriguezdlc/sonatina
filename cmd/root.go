package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/cmd/operation"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/arodriguezdlc/sonatina/utils"

	"github.com/spf13/viper"
)

var cfgFile string

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var rootCmd = &cobra.Command{
	Use:   "sonatina",
	Short: "A terraform based framework to work in an opinionated way.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		st, ok := err.(stackTracer)
		if ok {
			if viper.GetBool("EnableStacktrace") {
				fmt.Printf("%+v\n\n", st)
			}
			logrus.Fatalf("%v", st)
		} else {
			logrus.Fatalf("%v", err)
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SilenceUsage = true

	// Register subcommands
	rootCmd.AddCommand(operation.Apply)
	rootCmd.AddCommand(operation.Clone)
	rootCmd.AddCommand(operation.Create)
	rootCmd.AddCommand(operation.Delete)
	rootCmd.AddCommand(operation.Destroy)
	rootCmd.AddCommand(operation.Edit)
	rootCmd.AddCommand(operation.Get)
	rootCmd.AddCommand(operation.Init)
	rootCmd.AddCommand(operation.List)
	rootCmd.AddCommand(operation.Refresh)
	rootCmd.AddCommand(operation.Set)
	rootCmd.AddCommand(operation.Show)
	rootCmd.AddCommand(operation.Use)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match
	common.Fs = afero.NewOsFs()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	setLogFile()
	setLogLevel()
	if err == nil {
		logrus.Infoln("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.WithError(err).Warningln("couldn't read configuration file")
	}

	err = utils.NewFileIfNotExist(common.Fs, filepath.Join("~", ".sonatina", "config"))
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't create current file")
	}

	err = manager.InitializeManager(common.Fs, viper.GetString("ManagerConnector"))
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't initialize manager")
	}
}

func setLogFile() afero.File {
	filepath, err := homedir.Expand(viper.GetString("LogFile"))
	if err != nil {
		logrus.WithError(err).Fatal("couldn't get log file path")
	}

	logfile, err := common.Fs.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.WithError(err).Fatal("couldn't open file for logging")
	}
	logrus.SetOutput(logfile)
	return logfile
}

func setLogLevel() {
	level := viper.GetString("LogLevel")
	switch level {
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		//logrus.SetReportCaller(true)
	default:
		logrus.Fatalln("Unrecognized LogLevel: " + level)
	}

	logrus.Debugln("LogLevel: " + level)
}
