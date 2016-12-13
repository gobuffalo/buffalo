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
	"html/template"
	"os"

	"github.com/markbates/refresh/cmd"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs your Buffalo app in 'development' mode",
	Long: `Runs your Buffalo app in 'development' mode.
This includes rebuilding your application when files change.
This behavior can be changed in your .buffalo.dev.yml file.`,
	RunE: func(c *cobra.Command, args []string) error {
		os.Setenv("GO_ENV", "development")
		cfgFile := "./.buffalo.dev.yml"
		_, err := os.Stat(cfgFile)
		if err != nil {
			f, err := os.Create(cfgFile)
			if err != nil {
				return err
			}
			t, err := template.New("").Parse(nRefresh)
			err = t.Execute(f, map[string]interface{}{
				"name": "buffalo",
			})
			if err != nil {
				return err
			}
		}
		cmd.Run(cfgFile)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}
