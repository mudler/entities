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
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/mudler/entities/pkg/entities"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func writeUsers(store *EntitiesStore, targetDir string) error {
	dir := filepath.Join(targetDir, "users")

	if len(store.Users) > 0 {
		fmt.Println(fmt.Sprintf(
			"Creating %d users under the directory %s", len(store.Users), dir))

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"Error on creating directory %s: %s",
				dir, err.Error()),
			)
		}

		for k, u := range store.Users {

			file := filepath.Join(dir, fmt.Sprintf(
				"entity_user_%s.yaml", k))

			data, err := yaml.Marshal(&u)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on marshal user %s: %s",
					k, err.Error()))
			}

			err = ioutil.WriteFile(file, data, 0755)
			if err != nil {
				return errors.New(fmt.Sprint(
					"Error on write file %s: %s",
					file, err.Error()))
			}

		}
	}

	return nil
}

func writeGroups(store *EntitiesStore, targetDir string) error {
	dir := filepath.Join(targetDir, "groups")

	if len(store.Groups) > 0 {
		fmt.Println(fmt.Sprintf(
			"Creating %d groups under the directory %s", len(store.Groups), dir))

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"Error on creating directory %s: %s",
				dir, err.Error()),
			)
		}

		for k, g := range store.Groups {

			file := filepath.Join(dir, fmt.Sprintf(
				"entity_group_%s.yaml", k))

			data, err := yaml.Marshal(&g)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on marshal group %s: %s",
					k, err.Error()))
			}

			err = ioutil.WriteFile(file, data, 0755)
			if err != nil {
				return errors.New(fmt.Sprint(
					"Error on write file %s: %s",
					file, err.Error()))
			}

		}
	}

	return nil
}

func writeShadows(store *EntitiesStore, targetDir string) error {
	dir := filepath.Join(targetDir, "shadows")

	if len(store.Shadows) > 0 {
		fmt.Println(fmt.Sprintf(
			"Creating %d shadows under the directory %s", len(store.Shadows), dir))

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"Error on creating directory %s: %s",
				dir, err.Error()),
			)
		}

		for k, e := range store.Shadows {

			file := filepath.Join(dir, fmt.Sprintf(
				"entity_shadow_%s.yaml", k))

			data, err := yaml.Marshal(&e)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on marshal shadow %s: %s",
					k, err.Error()))
			}

			err = ioutil.WriteFile(file, data, 0755)
			if err != nil {
				return errors.New(fmt.Sprint(
					"Error on write file %s: %s",
					file, err.Error()))
			}

		}
	}

	return nil
}

func writeGShadows(store *EntitiesStore, targetDir string) error {
	dir := filepath.Join(targetDir, "gshadows")

	if len(store.GShadows) > 0 {
		fmt.Println(fmt.Sprintf(
			"Creating %d gshadows under the directory %s", len(store.GShadows), dir))

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"Error on creating directory %s: %s",
				dir, err.Error()),
			)
		}

		for k, e := range store.GShadows {

			file := filepath.Join(dir, fmt.Sprintf(
				"entity_gshadow_%s.yaml", k))

			data, err := yaml.Marshal(&e)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on marshal gshadow %s: %s",
					k, err.Error()))
			}

			err = ioutil.WriteFile(file, data, 0755)
			if err != nil {
				return errors.New(fmt.Sprint(
					"Error on write file %s: %s",
					file, err.Error()))
			}

		}
	}

	return nil
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump current system status in entities format",
	Long: `
Read system files and generate entities files to the specified directory.

To read /etc/shadow and /etc/gshadow requires root permissions.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		targetDir, _ := cmd.Flags().GetString("target-dir")
		if targetDir == "" {
			return errors.New("Missing mandatory target-dir.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, _ := cmd.Flags().GetString("target-dir")
		usersFile, _ := cmd.Flags().GetString("users-file")
		groupsFile, _ := cmd.Flags().GetString("groups-file")
		shadowFile, _ := cmd.Flags().GetString("shadow-file")
		gShadowFile, _ := cmd.Flags().GetString("gshadow-file")

		store := NewEntitiesStore()

		// Retrieve current information
		err := getCurrentStatus(store,
			usersFile, groupsFile, shadowFile, gShadowFile,
		)
		if err != nil {
			return errors.New(
				"Error on retrieve current entities status: " + err.Error(),
			)
		}

		err = writeUsers(store, targetDir)
		if err != nil {
			return err
		}

		err = writeGroups(store, targetDir)
		if err != nil {
			return err
		}

		err = writeShadows(store, targetDir)
		if err != nil {
			return err
		}

		err = writeGShadows(store, targetDir)
		if err != nil {
			return err
		}

		fmt.Println("All done.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	var flags = dumpCmd.Flags()
	flags.StringP("target-dir", "t", "",
		"Define the directory where dump entities files.")
	flags.String("users-file", UserDefault(""), "Define custom users file.")
	flags.String("groups-file", GroupsDefault(""), "Define custom groups file.")
	flags.String("shadow-file", ShadowDefault(""), "Define custom shadow file.")
	flags.String("gshadow-file", GShadowDefault(""), "Define custom gshadow file.")
}
