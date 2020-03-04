// Copyright © 2020 Ettore Di Giacinto <mudler@gentoo.org>
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

package entities

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	permbits "github.com/phayes/permbits"
	"github.com/pkg/errors"
)

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
		return "", Group{}, errors.New("Unexpected number of fields in /etc/Group: found " + strconv.Itoa(len(fs)))
	}

	gid, err := strconv.Atoi(fs[2])
	if err != nil {
		return "", Group{}, errors.New("Expected int for gid")
	}
	return fs[0], Group{fs[0], fs[1], gid, fs[3]}, nil
}

type Group struct {
	Name     string `yaml:"group_name"`
	Password string `yaml:"password"`
	Gid      int    `yaml:"gid"`
	Users    string `yaml:"users"`
}

func (u Group) String() string {
	return strings.Join([]string{u.Name,
		u.Password,
		strconv.Itoa(u.Gid),
		u.Users,
	}, ":")
}

func (u Group) Delete(s string) error {
	input, err := ioutil.ReadFile(s)
	if err != nil {
		return errors.Wrap(err, "Could not read input file")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}
	lines := bytes.Replace(input, []byte(u.String()+"\n"), []byte(""), 1)

	err = ioutil.WriteFile(s, []byte(lines), os.FileMode(permissions))
	if err != nil {
		return errors.Wrap(err, "Could not write")
	}

	return nil
}

func (u Group) Create(s string) error {
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

func (u Group) Apply(s string) error {
	current, err := ParseGroup(s)
	if err != nil {
		return errors.Wrap(err, "Failed parsing passwd")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}
	if _, ok := current[u.Name]; ok {
		input, err := ioutil.ReadFile(s)
		if err != nil {
			return errors.Wrap(err, "Could not read input file")
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if strings.HasPrefix(line, u.Name) {
				lines[i] = u.String()
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(s, []byte(output), os.FileMode(permissions))
		if err != nil {
			return errors.Wrap(err, "Could not write")
		}

	} else {
		return u.Create(s)
	}

	return nil
}
