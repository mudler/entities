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
	"regexp"
	"sort"
	"strconv"
	"time"

	. "github.com/mudler/entities/pkg/entities"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func createStore(specsdirs []string) (*EntitiesStore, error) {
	store := NewEntitiesStore()

	// Load sepcs
	for _, d := range specsdirs {
		err := store.Load(d)
		if err != nil {
			return store, errors.New(
				"Error on load specs from directory " + d + ": " + err.Error())
		}
	}

	return store, nil
}

func filterMatch(filter, field string) bool {
	if filter == "" {
		return true
	}

	r := regexp.MustCompile(filter)
	if r != nil {
		if r.MatchString(field) {
			return true
		}
		return false
	}
	return true
}

func listGroups(file, order, filter string, jsonOutput, groupHasShadow bool, specsdirs []string) error {
	var err error
	var mGShadows map[string]GShadow
	var mGroups map[string]Group

	if len(specsdirs) > 0 {
		store, err := createStore(specsdirs)
		if err != nil {
			return err
		}
		mGroups = store.Groups
	} else {
		file = GroupsDefault(file)

		mGroups, err = ParseGroup(file)
		if err != nil {
			return err
		}
	}

	// Sort group name
	groups := []string{}
	mGids := make(map[string]Group, 0)

	if order == "name" {
		for k, _ := range mGroups {
			if filterMatch(filter, k) {
				groups = append(groups, k)
			}
		}
		sort.Strings(groups)
	} else {

		gidList := []int{}

		for k, _ := range mGroups {
			if filterMatch(filter, k) {
				gid := fmt.Sprintf("%d", *mGroups[k].Gid)
				mGids[gid] = mGroups[k]
				gidList = append(gidList, *mGroups[k].Gid)
			}
		}

		sort.Ints(gidList)
		for _, g := range gidList {
			groups = append(groups, fmt.Sprintf("%d", g))
		}
	}

	if jsonOutput {
		res := []Group{}

		for _, group := range groups {
			gName := group
			if order == "id" {
				gName = mGids[group].Name
			}

			res = append(res, mGroups[gName])
		}

		data, _ := json.Marshal(res)
		fmt.Println(string(data))

	} else {

		// TODO: handle the file as an option
		if groupHasShadow {
			mGShadows, err = ParseGShadow(GShadowDefault(""))
			if err != nil {
				return err
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{
			Left:   true,
			Top:    true,
			Right:  true,
			Bottom: true,
		})
		headers := []string{
			"Group Name", "Encrypted Password", "Group ID", "Users",
		}

		if groupHasShadow {
			headers = append(headers, "Has GShadow")
		}
		table.SetHeader(headers)

		for _, group := range groups {
			gName := group
			if order == "id" {
				gName = mGids[group].Name
			}
			row := []string{
				mGroups[gName].Name,
				mGroups[gName].Password,
				fmt.Sprintf("%d", *mGroups[gName].Gid),
				mGroups[gName].Users,
			}

			if groupHasShadow {
				_, hasShadow := mGShadows[gName]

				row = append(row, fmt.Sprintf("%v", hasShadow))
			}

			table.Append(row)

		}

		table.Render()
	}

	return nil
}

func listShadows(file, order, filter string, jsonOutput, humanReadable bool, specsdirs []string) error {
	var err error
	var mShadows map[string]Shadow

	if len(specsdirs) > 0 {
		store, err := createStore(specsdirs)
		if err != nil {
			return err
		}
		mShadows = store.Shadows

	} else {
		file = ShadowDefault(file)
		mShadows, err = ParseShadow(file)
		if err != nil {
			return err
		}
	}

	// Sort group name
	shadows := []string{}

	for k, _ := range mShadows {
		if filterMatch(filter, k) {
			shadows = append(shadows, k)
		}
	}
	sort.Strings(shadows)

	if jsonOutput {
		res := []Shadow{}

		for _, s := range shadows {

			lastChanged := mShadows[s].LastChanged
			if humanReadable {
				i, err := strconv.Atoi(lastChanged)
				if err != nil {
					return err
				}

				unixSec := int64(i) * 24 * 60 * 60
				lctime := time.Unix(unixSec, 0)
				lastChanged = lctime.Format("2006-01-02T15:04:05Z")
			}

			expire := mShadows[s].Expire
			if humanReadable && expire != "" {
				i, err := strconv.Atoi(expire)
				if err != nil {
					return err
				}

				unixSec := int64(i) * 24 * 60 * 60
				extime := time.Unix(unixSec, 0)
				expire = extime.Format("2006-01-02T15:04:05Z")
			}

			shadow := mShadows[s]
			shadow.Expire = expire
			shadow.LastChanged = lastChanged

			res = append(res, shadow)
		}

		data, _ := json.Marshal(res)
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
			"Username", "Encrypted Password", "Last Password Change",
			"Minimum Changed", "Maximun Changed", "Warning Expiration",
			"Inactive", "Expire",
		})

		for _, s := range shadows {

			pass := mShadows[s].Password
			if len(pass) > 60 {
				pass = pass[0:60] + "\n" + pass[60:]
			}

			lastChanged := mShadows[s].LastChanged
			if humanReadable {
				i, err := strconv.Atoi(lastChanged)
				if err != nil {
					return err
				}

				unixSec := int64(i) * 24 * 60 * 60
				lctime := time.Unix(unixSec, 0)
				lastChanged = lctime.Format("2006-01-02T15:04:05Z")
			}

			expire := mShadows[s].Expire
			if humanReadable && expire != "" {
				i, err := strconv.Atoi(expire)
				if err != nil {
					return err
				}

				unixSec := int64(i) * 24 * 60 * 60
				extime := time.Unix(unixSec, 0)
				expire = extime.Format("2006-01-02T15:04:05Z")
			}

			table.Append([]string{
				mShadows[s].Username,
				pass,
				lastChanged,
				mShadows[s].MinimumChanged,
				mShadows[s].MaximumChanged,
				mShadows[s].Warn,
				mShadows[s].Inactive,
				expire,
			})

		}

		table.Render()
	}

	return nil
}

func listUsers(file, order, filter string, jsonOutput, userHasShadow bool, specsdirs []string) error {
	var err error
	var mUsers map[string]UserPasswd

	if len(specsdirs) > 0 {
		store, err := createStore(specsdirs)
		if err != nil {
			return err
		}
		mUsers = store.Users

	} else {
		file = UserDefault(file)
		mUsers, err = ParseUser(file)
		if err != nil {
			return err
		}
	}

	// Sort group name
	users := []string{}
	mUids := make(map[string]UserPasswd, 0)

	if order == "name" {
		for k, _ := range mUsers {
			if filterMatch(filter, k) {
				users = append(users, k)
			}
		}
		sort.Strings(users)
	} else {

		uidList := []int{}

		for k, _ := range mUsers {
			if filterMatch(filter, k) {
				uid := fmt.Sprintf("%d", mUsers[k].Uid)
				mUids[uid] = mUsers[k]
				uidList = append(uidList, mUsers[k].Uid)
			}
		}

		sort.Ints(uidList)
		for _, u := range uidList {
			users = append(users, fmt.Sprintf("%d", u))
		}
	}

	if jsonOutput {
		res := []UserPasswd{}

		for _, user := range users {
			uName := user
			if order == "id" {
				uName = mUids[user].Username
			}

			res = append(res, mUsers[uName])
		}

		data, _ := json.Marshal(res)
		fmt.Println(string(data))

	} else {

		// TODO: handle the file as an option
		mShadows, err := ParseShadow(ShadowDefault(""))
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{
			Left:   true,
			Top:    true,
			Right:  true,
			Bottom: true,
		})
		headers := []string{
			"Username", "Encrypted Password", "User ID", "Group ID", "Info",
			"Homedir", "Shell",
		}
		if userHasShadow {
			headers = append(headers, "With Shadow")
		}

		table.SetHeader(headers)

		for _, user := range users {
			uName := user
			if order == "id" {
				uName = mUids[user].Username
			}

			_, hasShadow := mShadows[uName]

			row := []string{
				mUsers[uName].Username,
				mUsers[uName].Password,
				fmt.Sprintf("%d", mUsers[uName].Uid),
				fmt.Sprintf("%d", mUsers[uName].Gid),
				mUsers[uName].Info,
				mUsers[uName].Homedir,
				mUsers[uName].Shell,
			}

			if userHasShadow {
				row = append(row, fmt.Sprintf("%v", hasShadow))
			}

			table.Append(row)
		}

		table.Render()
	}

	return nil
}

func listGshadows(file, order, filter string, jsonOutput bool, specsdirs []string) error {
	var err error
	var mGShadows map[string]GShadow

	if len(specsdirs) > 0 {
		store, err := createStore(specsdirs)
		if err != nil {
			return err
		}
		mGShadows = store.GShadows
	} else {
		file = GShadowDefault(file)
		mGShadows, err = ParseGShadow(file)
		if err != nil {
			return err
		}
	}

	// Sort group name
	gshadows := []string{}

	for k, _ := range mGShadows {
		if filterMatch(filter, k) {
			gshadows = append(gshadows, k)
		}
	}
	sort.Strings(gshadows)

	if jsonOutput {
		res := []GShadow{}

		for _, s := range gshadows {
			res = append(res, mGShadows[s])
		}

		data, _ := json.Marshal(res)
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
			"Group Name", "Encrypted Password",
			"Administrators", "Members",
		})

		for _, s := range gshadows {

			pass := mGShadows[s].Password
			if len(pass) > 60 {
				pass = pass[0:60] + "\n" + pass[60:]
			}

			table.Append([]string{
				mGShadows[s].Name,
				pass,
				mGShadows[s].Administrators,
				mGShadows[s].Members,
			})

		}

		table.Render()
	}

	return nil
}

var listCmd = &cobra.Command{
	Use:   "list <shadow|groups|users|gshadow>",
	Short: "Show entities availables",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Show the list of entities applied on current system.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Missing mandatory argument")
		}
		etype := args[0]
		switch etype {
		case "shadow", "groups", "users", "gshadow":
			break
		default:
			return errors.New(
				"Invalid entity type string. " +
					"First argument must contains one of this values: shadow|groups|users|gshadow.",
			)
		}

		order, _ := cmd.Flags().GetString("sort")
		if order != "name" && order != "id" {
			return errors.New(
				"Invalid order value. Admits values are: name|id.",
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ans error
		etype := args[0]
		file, _ := cmd.Flags().GetString("file")
		order, _ := cmd.Flags().GetString("sort")
		filter, _ := cmd.Flags().GetString("filter")
		jsonOutput, _ := cmd.Flags().GetBool("json")
		shadowHumanReadable, _ := cmd.Flags().GetBool("shadow-human-readable")
		userHasShadow, _ := cmd.Flags().GetBool("user-has-shadow")
		groupHasShadow, _ := cmd.Flags().GetBool("group-has-shadow")
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")

		switch etype {
		case "groups":
			ans = listGroups(file, order, filter, jsonOutput, groupHasShadow, specsdirs)
		case "shadow":
			ans = listShadows(file, order, filter, jsonOutput, shadowHumanReadable, specsdirs)
		case "users":
			ans = listUsers(file, order, filter, jsonOutput, userHasShadow, specsdirs)
		case "gshadow":
			ans = listGshadows(file, order, filter, jsonOutput, specsdirs)
		default:
			return errors.New("Unexpected entity type " + etype)
		}

		return ans
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	var flags = listCmd.Flags()
	flags.StringP("sort", "s", "name", "Sort list by: name|id")
	flags.String("filter", "", "Filter entities by name. It uses the filter as regex")
	flags.Bool("json", false, "Show in JSON format.")
	flags.Bool("shadow-human-readable", false, "Show shadow days in human readable format.")
	flags.Bool("user-has-shadow", false, "Check if exists a map of the users in the /etc/shadow file. (Available only in table format)")
	flags.Bool("group-has-shadow", false, "Check if exists a map of the users in the /etc/gshadow file. (Available only in table format)")

	flags.StringArray("specs-dir", []string{},
		"Define the directory where read entities specs in alternative to the system files.")
}
