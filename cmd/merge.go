// Copyright © 2021 Daniele Rondina <geaaru@sabayonlinux.org>
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
	"errors"
	"fmt"

	. "github.com/mudler/entities/pkg/entities"

	"github.com/spf13/cobra"
)

func mergeEntity(store, currentStore *EntitiesStore,
	entityName, usersFile, groupsFile, shadowFile, gShadowFile string) error {

	var err error
	found := false

	// Merge user if present.
	if u, ok := store.Users[entityName]; ok {
		found = true

		var newEntity Entity = u

		if cu, ok := currentStore.Users[entityName]; ok {
			// POST: the entity is already present. I merge it.
			newEntity, err = cu.Merge(u)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on merge user %s: %s", cu.Username, err.Error()))
			}
		}

		err = newEntity.Apply(usersFile, false)
		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Error on apply user %s: %s", entityName, err.Error()))
		}

		fmt.Println(fmt.Sprintf(
			"Merged users %s.", entityName))
	}

	if g, ok := store.Groups[entityName]; ok {
		found = true
		var newEntity Entity = g

		if cu, ok := currentStore.Groups[entityName]; ok {
			// POST: the entity is already present. I merge it
			newEntity, err = cu.Merge(g)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on merge group %s: %s", entityName, err.Error()))
			}
		}

		err = newEntity.Apply(groupsFile, false)
		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Error on apply group %s: %s", entityName, err.Error()))
		}

		fmt.Println(fmt.Sprintf(
			"Merged group %s.", entityName))
	}

	if s, ok := store.Shadows[entityName]; ok {
		found = true

		var newEntity Entity = s

		if cs, ok := currentStore.Shadows[entityName]; ok {
			// POST: the entity is already present. I merge it
			newEntity, err = cs.Merge(s)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on merge shadow %s: %s", entityName, err.Error()))
			}
		}

		err = newEntity.Apply(shadowFile, false)
		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Error on apply shadow %s: %s", entityName, err.Error()))
		}

		fmt.Println(fmt.Sprintf(
			"Merged shadow %s.", entityName))
	}

	if s, ok := store.GShadows[entityName]; ok {
		found = true
		var newEntity Entity = s

		if cs, ok := currentStore.GShadows[entityName]; ok {
			// POST: the entity is already present. I merge it
			newEntity, err = cs.Merge(s)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on merge gshadow %s: %s", entityName, err.Error()))
			}
		}

		err = newEntity.Apply(gShadowFile, false)
		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Error on apply gshadow %s: %s", entityName, err.Error()))
		}

		fmt.Println(fmt.Sprintf(
			"Merged gshadow %s.", entityName))

	}

	if !found {
		return errors.New(
			fmt.Sprintf("No entities found with name %s.", entityName))
	}

	return nil
}

func mergeAllEntities(store, currentStore *EntitiesStore,
	usersFile, groupsFile, shadowFile, gShadowFile string) error {
	entities := make(map[string]bool, 0)

	for k, _ := range store.Users {
		entities[k] = true
	}

	for k, _ := range store.Groups {
		entities[k] = true
	}

	for k, _ := range store.Shadows {
		entities[k] = true
	}

	for k, _ := range store.GShadows {
		entities[k] = true
	}

	if len(entities) == 0 {
		return errors.New("No entities to merge")
	}

	for k, _ := range entities {
		err := mergeEntity(store, currentStore, k,
			usersFile, groupsFile, shadowFile, gShadowFile)
		if err != nil {
			return err
		}
	}

	return nil
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge one or more entities.",
	Long: `
Merge one or more entities read from a specified directory to
the existing system. If the entity is already present it merges
entities without override uid/gid or password.

To read /etc/shadow and /etc/gshadow requires root permissions.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")
		if len(specsdirs) == 0 {
			return errors.New("At least one specs directory is needed.")
		}

		entity, _ := cmd.Flags().GetString("entity")
		all, _ := cmd.Flags().GetBool("all")

		if entity == "" && !all {
			return errors.New("You need choice an entity or to use --all.")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")
		usersFile, _ := cmd.Flags().GetString("users-file")
		groupsFile, _ := cmd.Flags().GetString("groups-file")
		shadowFile, _ := cmd.Flags().GetString("shadow-file")
		gShadowFile, _ := cmd.Flags().GetString("gshadow-file")
		entity, _ := cmd.Flags().GetString("entity")
		all, _ := cmd.Flags().GetBool("all")

		store := NewEntitiesStore()
		currentStore := NewEntitiesStore()

		// Load sepcs
		for _, d := range specsdirs {
			err := store.Load(d)
			if err != nil {
				return errors.New(
					"Error on load specs from directory " + d + ": " + err.Error())
			}
		}

		// Retrieve current information
		err := getCurrentStatus(currentStore,
			usersFile, groupsFile, shadowFile, gShadowFile,
		)
		if err != nil {
			return errors.New(
				"Error on retrieve current entities status: " + err.Error(),
			)
		}

		if entity != "" {
			err = mergeEntity(store, currentStore, entity,
				usersFile, groupsFile, shadowFile, gShadowFile,
			)
		} else if all {
			err = mergeAllEntities(store, currentStore, usersFile,
				groupsFile, shadowFile, gShadowFile,
			)
		}

		if err != nil {
			return err
		}

		fmt.Println("All done.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	var flags = mergeCmd.Flags()
	flags.StringArrayP("specs-dir", "s", []string{},
		"Define the directory where read entities specs. At least one directory is needed.")
	flags.BoolP("all", "a", false,
		"Merge all entities available on the specified directories.",
	)
	flags.StringP("entity", "e", "",
		"Entity name to merge as groups,user,shadow and gshadow.")
	flags.String("users-file", UserDefault(""), "Define custom users file.")
	flags.String("groups-file", GroupsDefault(""), "Define custom groups file.")
	flags.String("shadow-file", ShadowDefault(""), "Define custom shadow file.")
	flags.String("gshadow-file", GShadowDefault(""), "Define custom gshadow file.")
}
