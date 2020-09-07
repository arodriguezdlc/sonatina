/*
Copyright © 2020 ALBERTO RODRIGUEZ <arodriguezdlc@gmail.com>

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
	"os/exec"

	"github.com/arodriguezdlc/sonatina/cmd"
	"github.com/sirupsen/logrus"
)

func sshAgentUnix() {
	err := exec.Command("sh", "-c", "eval", "`ssh-agent`").Run()
	if err != nil {
		logrus.WithError(err).Warning("couldn't setup ssh-agent")
		return
	}

	err = exec.Command("sh", "-c", "ssh-add").Run()
	if err != nil {
		logrus.WithError(err).Warning("couldn't setup ssh-add")
		return
	}
}

func main() {
	// XXX: currently this function is necessary to enable ssh connections to git repositories.
	// Only works on unix systems.
	sshAgentUnix()

	cmd.Execute()
}
