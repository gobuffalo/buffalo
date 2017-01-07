// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// task command is a forward to grift tasks
var taskCommand = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t"},
	Short:   "Runs your grift tasks",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := exec.LookPath("grift")
		if err != nil {
			return errors.New("we could not find \"grift\" in your path.\n You must first install \"grift\" in order to use the Buffalo console:\n\n $ go get github.com/markbates/grift")
		}

		_, err = os.Stat("grifts")
		if err != nil {
			return errors.New("seems there is no grift folder on your current directory, please ensure you're inside your buffalo app root")
		}

		cmd := exec.Command("grift", args...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	},
}

func init() {
	RootCmd.AddCommand(taskCommand)
}
