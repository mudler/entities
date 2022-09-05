/*
Copyright © 2019 Ettore Di Giacinto <mudler@gentoo.org>
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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var entityFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "entities",
	Short: "Modern go identity manager for UNIX systems",
	Long: `Entities is a modern groups and user manager for Unix system. It allows to create/delete user and groups 
in a system given policies following the entities yaml format.

For example:

	$> entities apply <entity.yaml>
	$> entities delete <entity.yaml>
	$> entities create <entity.yaml>
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&entityFile, "file", "f", "", "File to manipulate ( e.g. /etc/passwd ) ")
}
