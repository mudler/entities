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

// ParseGShadow opens the file and parses it into a map from usernames to Entries
func ParseGShadow(path string) (map[string]GShadow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ParseGShadowReader(file)
}

// ParseGShadowReader consumes the contents of r and parses it into a map from
// usernames to Entries
func ParseGShadowReader(r io.Reader) (map[string]GShadow, error) {
	lines := bufio.NewReader(r)
	entries := make(map[string]GShadow)
	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			break
		}
		name, entry, err := parseGShadowLine(string(copyBytes(line)))
		if err != nil {
			return nil, err
		}
		entries[name] = entry
	}
	return entries, nil
}

func parseGShadowLine(line string) (string, GShadow, error) {
	fs := strings.Split(line, ":")
	if len(fs) != 4 {
		return "", GShadow{}, errors.New("Unexpected number of fields in /etc/GShadow: found " + strconv.Itoa(len(fs)))
	}

	return fs[0], GShadow{fs[0], fs[1], fs[2], fs[3]}, nil
}

type GShadow struct {
	Name           string `yaml:"name"`
	Password       string `yaml:"password"`
	Administrators string `yaml:"administrators"`
	Members        string `yaml:"members"`
}

func (u GShadow) String() string {
	return strings.Join([]string{
		u.Name,
		u.Password,
		u.Administrators,
		u.Members,
	}, ":")
}

func (u GShadow) Delete(s string) error {
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

func (u GShadow) Apply(s string) error {
	current, err := ParseGShadow(s)
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
			if strings.HasPrefix(line, u.Name) { // FIXME: Drop greedy match
				lines[i] = u.String()
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(s, []byte(output), os.FileMode(permissions))
		if err != nil {
			return errors.Wrap(err, "Could not write")
		}

	} else {
		// Add it
		f, err := os.OpenFile(s, os.O_APPEND|os.O_WRONLY, os.FileMode(permissions))
		if err != nil {
			return errors.Wrap(err, "Could not read")
		}

		defer f.Close()

		if _, err = f.WriteString(u.String() + "\n"); err != nil {
			return errors.Wrap(err, "Could not write")
		}

	}

	return nil
}
