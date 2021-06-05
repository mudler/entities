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
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	permbits "github.com/phayes/permbits"
	"github.com/pkg/errors"
	passwd "github.com/willdonnelly/passwd"
)

func UserDefault(s string) string {
	if s == "" {
		s = "/etc/passwd"
	}
	return s
}

type UserPasswd struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Uid      int    `yaml:"uid"`
	Gid      int    `yaml:"gid"`
	Info     string `yaml:"info"`
	Homedir  string `yaml:"homedir"`
	Shell    string `yaml:"shell"`
}

func ParseUser(path string) (map[string]UserPasswd, error) {
	ans := make(map[string]UserPasswd, 0)

	current, err := passwd.ParseFile(path)
	if err != nil {
		return ans, errors.Wrap(err, "Failed parsing passwd")
	}
	_, err = permbits.Stat(path)
	if err != nil {
		return ans, errors.Wrap(err, "Failed getting permissions")
	}

	for k, v := range current {
		uid, err := strconv.Atoi(v.Uid)
		if err != nil {
			return ans, errors.Wrap(err, "Invalid uid found")
		}

		gid, err := strconv.Atoi(v.Gid)
		if err != nil {
			return ans, errors.Wrap(err, "Invalid gid found")
		}

		ans[k] = UserPasswd{
			Username: k,
			Password: v.Pass,
			Uid:      uid,
			Gid:      gid,
			Info:     v.Gecos,
			Homedir:  v.Home,
			Shell:    v.Shell,
		}
	}

	return ans, nil
}

func (u UserPasswd) String() string {
	return strings.Join([]string{u.Username,
		u.Password,
		strconv.Itoa(u.Uid),
		strconv.Itoa(u.Gid),
		u.Info,
		u.Homedir,
		u.Shell,
	}, ":")
}

func (u UserPasswd) Delete(s string) error {
	s = UserDefault(s)
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

func (u UserPasswd) Create(s string) error {
	s = UserDefault(s)
	current, err := passwd.ParseFile(s)
	if err != nil {
		return errors.Wrap(err, "Failed parsing passwd")
	}
	if _, ok := current[u.Username]; ok {
		return errors.New("Entity already present")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}
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

func (u UserPasswd) Apply(s string) error {
	s = UserDefault(s)
	current, err := passwd.ParseFile(s)
	if err != nil {
		return errors.Wrap(err, "Failed parsing passwd")
	}
	permissions, err := permbits.Stat(s)
	if err != nil {
		return errors.Wrap(err, "Failed getting permissions")
	}

	if _, ok := current[u.Username]; ok {

		input, err := ioutil.ReadFile(s)
		if err != nil {
			return errors.Wrap(err, "Could not read input file")
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if entityIdentifier(line) == u.Username {
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
		return u.Create(s)
	}

	return nil
}
