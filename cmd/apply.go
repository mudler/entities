// Copyright © 2019 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.
package cmd

import (
	. "github.com/mudler/entities/pkg/entities"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies an entity",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Applies a entity yaml file to your system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &Parser{}

		safe, _ := cmd.Flags().GetBool("safe")

		entity, err := p.ReadEntity(args[0])
		if err != nil {
			return err
		}

		return entity.Apply(entityFile, safe)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	var flags = applyCmd.Flags()
	flags.Bool("safe", false,
		"Avoid to override existing entity if it has difference or if the id is used in a different way.")
}
