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
