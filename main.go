/*
Copyright © 2019 ALBERTO RODRIGUEZ <arodriguezdlc@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/arodriguezdlc/sonatina/cmd"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func main() {

	setLogLevel(viper.GetString("LogLevel"))

	var fs afero.Fs
	if viper.GetBool("TestFilesystem") {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}

	err := manager.InitializeManager(fs)
	if err != nil {
		logrus.Fatalln(err)
	}

	cmd.Execute()
}

func setLogLevel(level string) {
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
