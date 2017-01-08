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

package generate

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

//ActionCmd is the cmd that generates actions.
var ActionCmd = &cobra.Command{
	Use:     "action [name] [actionName...]",
	Aliases: []string{"a"},
	Short:   "Generates new action(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("you should provide action name and handler name at least")
		}

		_, err := os.Stat("actions")
		if err != nil {
			return errors.New("actions directory not found, ensure you're inside your buffalo folder")
		}

		//checks: file exists
		return nil
	},
}
