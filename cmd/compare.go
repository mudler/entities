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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	. "github.com/mudler/entities/pkg/entities"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type EntityDifference struct {
	OriginalEntity Entity `json:"originalEntity,omitempty" yaml:"originalEntity,omitempty"`
	TargetEntity   Entity `json:"targetEntity,omitempty" yaml:"targetEntity,omitempty"`
	Kind           string `json:"kind" yaml:"kind"`
	Descr          string `json:"descr,omitempty" yaml:"descr,omitempty"`
	Missing        bool   `json:"missing" yaml:"missing"`
}

func getCurrentStatus(store *EntitiesStore, usersFile, groupsFile, shadowFile, gshadowFile string) error {

	mUsers, err := ParseUser(usersFile)
	if err != nil {
		return err
	}

	mGroups, err := ParseGroup(groupsFile)
	if err != nil {
		return err
	}

	mShadows, err := ParseShadow(shadowFile)
	if err != nil {
		return err
	}

	mGShadows, err := ParseGShadow(gshadowFile)
	if err != nil {
		return err
	}

	store.Users = mUsers
	store.Groups = mGroups
	store.Shadows = mShadows
	store.GShadows = mGShadows

	return nil
}

func compare(currentStore, store *EntitiesStore, jsonOutput bool) error {

	differences := []EntityDifference{}

	// Check users: I check that all entities defined in the specs are available and equal.
	// Not in reverse.
	for name, u := range store.Users {
		cUser, ok := currentStore.GetUser(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: u,
				Missing:      true,
				Kind:         u.GetKind(),
				Descr:        fmt.Sprintf("User %s is not present.", name),
			})
			continue
		}

		if (u.Uid >= 0 && cUser.Uid != u.Uid) ||
			(u.Group == "" && cUser.Gid != u.Gid) ||
			cUser.Homedir != u.Homedir || cUser.Shell != u.Shell {
			differences = append(differences, EntityDifference{
				OriginalEntity: cUser,
				TargetEntity:   u,
				Missing:        false,
				Kind:           u.GetKind(),
				Descr:          fmt.Sprintf("User %s has difference.", name),
			})
		}
	}

	// Check groups
	for name, g := range store.Groups {
		cGroup, ok := currentStore.GetGroup(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: g,
				Missing:      true,
				Kind:         g.GetKind(),
				Descr:        fmt.Sprintf("Group %s is not present.", name),
			})
			continue
		}

		if cGroup.Password != g.Password ||
			(g.Gid != nil && *g.Gid >= 0 && cGroup.Gid != g.Gid) ||
			cGroup.Users != g.Users {
			differences = append(differences, EntityDifference{
				OriginalEntity: cGroup,
				TargetEntity:   g,
				Missing:        false,
				Kind:           g.GetKind(),
				Descr:          fmt.Sprintf("Group %s has difference.", name),
			})
		}
	}

	// Check shadow
	for name, s := range store.Shadows {
		cShadow, ok := currentStore.GetShadow(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: s,
				Missing:      true,
				Kind:         s.GetKind(),
				Descr:        fmt.Sprintf("Shadow with username %s is not present.", name),
			})
			continue
		}

		if cShadow.MinimumChanged != s.MinimumChanged ||
			cShadow.MaximumChanged != s.MaximumChanged ||
			cShadow.Warn != s.Warn ||
			cShadow.Inactive != s.Inactive ||
			cShadow.Expire != s.Expire {
			differences = append(differences, EntityDifference{
				OriginalEntity: cShadow,
				TargetEntity:   s,
				Missing:        false,
				Kind:           s.GetKind(),
				Descr:          fmt.Sprintf("Shadow with user %s has difference.", name),
			})
		}
	}

	// Check gshadow
	for name, s := range store.GShadows {
		cGShadow, ok := currentStore.GetGShadow(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: s,
				Missing:      true,
				Kind:         s.GetKind(),
				Descr:        fmt.Sprintf("GShadow with name %s is not present.", name),
			})
			continue
		}

		if cGShadow.Password != s.Password ||
			cGShadow.Administrators != s.Administrators ||
			cGShadow.Members != s.Members {
			differences = append(differences, EntityDifference{
				OriginalEntity: cGShadow,
				TargetEntity:   s,
				Missing:        false,
				Kind:           s.GetKind(),
				Descr:          fmt.Sprintf("GShadow with name %s has difference.", name),
			})
		}
	}

	if jsonOutput {
		data, _ := json.Marshal(differences)
		fmt.Println(string(data))
	} else {

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{
			Left:   true,
			Top:    true,
			Right:  true,
			Bottom: true,
		})
		table.SetColWidth(50)
		table.SetHeader([]string{
			"Kind", "Name", "Missing", "Difference",
		})
		for _, d := range differences {
			var name string
			switch d.Kind {
			case UserKind:
				name = (d.TargetEntity.(UserPasswd)).Username
			case ShadowKind:
				name = (d.TargetEntity.(Shadow)).Username
			case GroupKind:
				name = (d.TargetEntity.(Group)).Name
			case GShadowKind:
				name = (d.TargetEntity.(GShadow)).Name
			}

			table.Append([]string{
				d.Kind,
				name,
				fmt.Sprintf("%v", d.Missing),
				d.Descr,
			})

		}

		table.Render()

	}

	return nil
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare entities present with specs.",
	Long: `
Compare entities of the system with the specs available in the specified directory.

To read /etc/shadow and /etc/gshadow requires root permissions.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")
		if len(specsdirs) == 0 {
			return errors.New("At least one specs directory is needed.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")
		usersFile, _ := cmd.Flags().GetString("users-file")
		groupsFile, _ := cmd.Flags().GetString("groups-file")
		shadowFile, _ := cmd.Flags().GetString("shadow-file")
		gShadowFile, _ := cmd.Flags().GetString("gshadow-file")
		jsonOutput, _ := cmd.Flags().GetBool("json")

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

		err = compare(currentStore, store, jsonOutput)
		if err != nil {
			return errors.New(
				"Error on compare entities stores: " + err.Error(),
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)

	var flags = compareCmd.Flags()
	flags.StringArrayP("specs-dir", "s", []string{},
		"Define the directory where read entities specs. At least one directory is needed.")
	flags.String("users-file", UserDefault(""), "Define custom users file.")
	flags.String("groups-file", GroupsDefault(""), "Define custom groups file.")
	flags.String("shadow-file", ShadowDefault(""), "Define custom shadow file.")
	flags.String("gshadow-file", GShadowDefault(""), "Define custom gshadow file.")
	flags.Bool("json", false, "Show in JSON format.")
}
