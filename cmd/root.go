package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/cmd/operation"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sonatina",
	Short: "A terraform based framework",
	Long:  `TO DO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(fs afero.Fs) {
	common.Fs = fs

	err := rootCmd.Execute()
	if err != nil {

		st, ok := err.(stackTracer)
		if ok {
			fmt.Printf("%+v\n\n", st)
			logrus.Fatalf("%v", st)
		} else {
			logrus.Fatalf("%v", err)
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sonatina.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.SilenceUsage = true

	// Register CRUD subcommands
	rootCmd.AddCommand(operation.Clone)
	rootCmd.AddCommand(operation.Create)
	rootCmd.AddCommand(operation.Delete)
	rootCmd.AddCommand(operation.Init)
	rootCmd.AddCommand(operation.List)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".sonatina" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sonatina")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
