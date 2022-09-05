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
	. "github.com/mudler/entities/pkg/entities"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create an entity",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Create a entity to your system from yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &Parser{}

		entity, err := p.ReadEntity(args[0])
		if err != nil {
			return err
		}

		return entity.Create(entityFile)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
