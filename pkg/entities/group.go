/*
Copyright © 2020 Ettore Di Giacinto <mudler@mocaccino.org>
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

package entities

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofrs/flock"
	permbits "github.com/phayes/permbits"
	"github.com/pkg/errors"
)

func GroupsDefault(s string) string {
	if s == "" {
		s = os.Getenv(ENTITY_ENV_DEF_GROUPS)
		if s == "" {
			s = "/etc/group"
		}
	}
	return s
}

// ParseGroup opens the file and parses it into a map from usernames to Entries
func ParseGroup(path string) (map[string]Group, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ParseGroupReader(file)
}

// ParseGroupReader consumes the contents of r and parses it into a map from
// usernames to Entries
func ParseGroupReader(r io.Reader) (map[string]Group, error) {
	lines := bufio.NewReader(r)
	entries := make(map[string]Group)
	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			break
		}
		name, entry, err := parseGroupLine(string(copyBytes(line)))
		if err != nil {
			return nil, err
		}
		entries[name] = entry
	}
	return entries, nil
}

func parseGroupLine(line string) (string, Group, error) {
	fs := strings.Split(line, ":")
	if len(fs) != 4 {
		return "", Group{}, errors.New(
			"Unexpected number of fields in /etc/group: found " + strconv.Itoa(len(fs)) +
				" - " + line)
	}

	gid, err := strconv.Atoi(fs[2])
	if err != nil {
		return "", Group{}, errors.New("Expected int for gid")
	}
	return fs[0], Group{fs[0], fs[1], &gid, fs[3]}, nil
}

func groupGetFreeGid(path string) (int, error) {
	result := -1
	allGroups, _ := ParseGroup(path)
	groupSet := make(map[int]struct{})

	for _, groupID := range allGroups {
		groupSet[*groupID.Gid] = struct{}{}
	}

	for i := HumanIDMin; i <= HumanIDMax; i++ {
		if _, found := groupSet[i]; found {
			continue // uid in use, skip it
		}
		result = i // found a free one, stop here
		break
	}

	if result == -1 {
		return result, errors.New("no available gid in range")
	}

	return result, nil
}

type Group struct {
	Name     string `yaml:"group_name"`
	Password string `yaml:"password"`
	Gid      *int   `yaml:"gid"`
	Users    string `yaml:"users"`
}

func (u Group) GetKind() string { return GroupKind }

func (u Group) prepare(s string) (Group, error) {
	if u.Gid != nil && *u.Gid < 0 {
		// POST: dynamic group
		gid, err := groupGetFreeGid(s)
		if err != nil {
			return u, err
		}
		u.Gid = &gid
	}

	return u, nil
}

func (u Group) String() string {
	var gid string
	if u.Gid == nil {
		gid = ""
	} else {
		gid = strconv.Itoa(*u.Gid)
	}
	return strings.Join([]string{u.Name,
		u.Password,
		gid,
		u.Users,
	}, ":")
}

func (u Group) Delete(s string) error {
	s = GroupsDefault(s)
	d, err := RetryForDuration()
	if err != nil {
		return errors.Wrap(err, "Failed getting delay")
	}

	baseName := filepath.Base(s)
	if _, ok := os.Stat("/run/lock/"); os.IsNotExist(ok) {
		_ = os.Mkdir("/run/lock/", 0755)
	}
	fileLock := flock.New(fmt.Sprintf("/run/lock/%s.lock", baseName))
	defer os.Remove(fileLock.Path())
	defer fileLock.Close()
	lockCtx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	i, err := RetryIntervalDuration()
	if err != nil {
		return errors.Wrap(err, "Failed getting interval")
	}
	locked, err := fileLock.TryLockContext(lockCtx, i)
	if err != nil || !locked {
		return errors.Wrap(err, "Failed locking file")
	}
	input, err := os.ReadFile(s)
	if err != nil {
		return errors.Wrap(err, "Could not read input file")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}

	// Drop the line which match the identifier. Don't look at the content as in other cases
	lines := strings.Split(string(input), "\n")
	var toremove int
	for i, line := range lines {
		if entityIdentifier(line) == u.Name {
			toremove = i
		}
	}

	// Remove the element at index i from a.
	copy(lines[toremove:], lines[toremove+1:]) // Shift a[i+1:] left one index.
	lines[len(lines)-1] = ""                   // Erase last element (write zero value).
	lines = lines[:len(lines)-1]               // Truncate slice.

	output := strings.Join(lines, "\n")

	err = os.WriteFile(s, []byte(output), os.FileMode(permissions))
	if err != nil {
		return errors.Wrap(err, "Could not write")
	}

	return nil
}

func (u Group) Create(s string) error {
	s = GroupsDefault(s)
	d, err := RetryForDuration()
	if err != nil {
		return errors.Wrap(err, "Failed getting delay")
	}

	baseName := filepath.Base(s)
	if _, ok := os.Stat("/run/lock/"); os.IsNotExist(ok) {
		_ = os.Mkdir("/run/lock/", 0755)
	}
	fileLock := flock.New(fmt.Sprintf("/run/lock/%s.lock", baseName))
	defer os.Remove(fileLock.Path())
	defer fileLock.Close()
	lockCtx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	i, err := RetryIntervalDuration()
	if err != nil {
		return errors.Wrap(err, "Failed getting interval")
	}
	locked, err := fileLock.TryLockContext(lockCtx, i)
	if err != nil || !locked {
		return errors.Wrap(err, "Failed locking file")
	}

	u, err = u.prepare(s)
	if err != nil {
		return errors.Wrap(err, "Failed entity preparation")
	}

	current, err := ParseGroup(s)
	if err != nil {
		return errors.Wrap(err, "Failed parsing passwd")
	}
	if _, ok := current[u.Name]; ok {
		return errors.New("Entity already present")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}
	// Add it
	f, err := os.OpenFile(s, os.O_APPEND|os.O_WRONLY, os.FileMode(permissions))
	if err != nil {
		return errors.Wrap(err, "Could not read")
	}

	defer f.Close()

	if _, err = f.WriteString(u.String() + "\n"); err != nil {
		return errors.Wrap(err, "Could not write")
	}
	return nil
}

func Unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		// Ignore invalid string. Workaround to broken /etc/groups generated by
		// previous version of entities
		if entry == "" {
			continue
		}

		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (u Group) Apply(s string, safe bool) error {
	if u.Name == "" {
		return errors.New("Empty group name")
	}

	s = GroupsDefault(s)

	u, err := u.prepare(s)
	if err != nil {
		return errors.Wrap(err, "Failed entity preparation")
	}

	current, err := ParseGroup(s)
	if err != nil {
		return errors.Wrap(err, "Failed parsing passwd")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}

	if safe && u.Gid != nil {
		// Avoid this check if the gid is not
		// present. For example for the specs where
		// we add users to a group.

		mGids := make(map[int]*Group)

		// Create gids to check gid mismatch
		// Maybe could be done always.
		for _, e := range current {
			mGids[*e.Gid] = &e
		}

		if e, present := mGids[*u.Gid]; present {
			if e.Name != u.Name {
				return errors.Wrap(err,
					fmt.Sprintf("Gid %d is already used on group %s",
						*u.Gid, u.Name))
			}
		}
	}

	if _, ok := current[u.Name]; ok {
		d, err := RetryForDuration()
		if err != nil {
			return errors.Wrap(err, "Failed getting delay")
		}

		baseName := filepath.Base(s)
		if _, ok := os.Stat("/run/lock/"); os.IsNotExist(ok) {
			_ = os.Mkdir("/run/lock/", 0755)
		}
		fileLock := flock.New(fmt.Sprintf("/run/lock/%s.lock", baseName))
		defer os.Remove(fileLock.Path())
		defer fileLock.Close()
		lockCtx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		i, err := RetryIntervalDuration()
		if err != nil {
			return errors.Wrap(err, "Failed getting interval")
		}
		locked, err := fileLock.TryLockContext(lockCtx, i)
		if err != nil || !locked {
			return errors.Wrap(err, "Failed locking file")
		}

		input, err := os.ReadFile(s)
		if err != nil {
			return errors.Wrap(err, "Could not read input file")
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if entityIdentifier(line) == u.Name {
				// Merge the groups, don't override the whole user.
				_, g, err := parseGroupLine(lines[i])
				if err != nil {
					return errors.Wrap(err, "Failed parsing current group")
				}
				if len(g.Users) > 0 {
					currentUsers := strings.Split(g.Users, ",")
					if u.Users != "" {
						currentUsers = append(currentUsers, strings.Split(u.Users, ",")...)
					}
					u.Users = strings.Join(Unique(currentUsers), ",")
				}

				if !safe {
					if len(u.Password) == 0 {
						u.Password = g.Password
					}
					if u.Gid == nil {
						u.Gid = g.Gid
					}
				} else {
					// Maintain existing group id and password
					u.Gid = g.Gid
					u.Password = g.Password
				}

				lines[i] = u.String()
			}
		}
		output := strings.Join(lines, "\n")
		err = os.WriteFile(s, []byte(output), os.FileMode(permissions))
		if err != nil {
			return errors.Wrap(err, "Could not write")
		}

	} else {
		return u.Create(s)
	}

	return nil
}
