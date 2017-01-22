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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var prevB *cobra.Command

	for i := 0; i < 7; i++ {
		var currB = cobra.Command{}
		var subcmd, subsubcmd string

		currB.Short = "buffalo?"
		currB.Long = "buffalo!"

		if i == 0 {
			subcmd = "buffalo."
		} else if i == 1 || i == 5 {
			subcmd = "Buffalo"
		} else {
			subcmd = "buffalo"
		}

		if i == 2 || i == 6 {
			subsubcmd = "Buffalo"
		} else if i == 1 {
			subsubcmd = "buffalo."
		} else {
			subsubcmd = "buffalo"
		}
		subcmd += fmt.Sprintf(" %s", subsubcmd)

		currB.Use = subcmd

		currB.RunE = func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || !strings.EqualFold(args[0], subcmd) {
				return fmt.Errorf("Check your command. Are you sure it is not %s?", subsubcmd)
			}
			return nil
		}

		if i == 0 {
			currB.RunE = func(cmd *cobra.Command, args []string) error {
				if len(args) != 0 {
					return fmt.Errorf("Check your command. Are you sure it is not %s?", subcmd)
				}
				fmt.Println("Buffalo buffalo Buffalo buffalo buffalo buffalo Buffalo buffalo.")
				return nil
			}
		}

		if i != 0 {
			currB.AddCommand(prevB)
		}
		prevB = &currB

	}
	RootCmd.AddCommand(prevB)
}
