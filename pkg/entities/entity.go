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

import "strings"

const (
	ENTITY_ENV_DEF_GROUPS        = "ENTITY_DEFAULT_GROUPS"
	ENTITY_ENV_DEF_PASSWD        = "ENTITY_DEFAULT_PASSWD"
	ENTITY_ENV_DEF_SHADOW        = "ENTITY_DEFAULT_SHADOW"
	ENTITY_ENV_DEF_GSHADOW       = "ENTITY_DEFAULT_GSHADOW"
	ENTITY_ENV_DEF_DYNAMIC_RANGE = "ENTITY_DYNAMIC_RANGE"
)

// Entity represent something that needs to be applied to a file

type Entity interface {
	GetKind() string
	String() string
	Delete(s string) error
	Create(s string) error
	Apply(s string, safe bool) error
}

func entityIdentifier(s string) string {
	fs := strings.Split(s, ":")
	if len(fs) == 0 {
		return ""
	}

	return fs[0]
}
