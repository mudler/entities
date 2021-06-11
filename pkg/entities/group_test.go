// Copyright © 2020 Ettore Di Giacinto <mudler@gentoo.org>
//                  Daniele Rondina <geaaru@sabayonlinux.org>
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

package entities_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Group", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/group/group", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/group/update.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Group).Name).Should(Equal("sddm"))

			err = entity.Apply(tmpFile.Name(), false)
			Expect(err).Should(BeNil())

			dat, err := ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`nm-openconnect:x:979:
sddm:xx:1:one,two,tree
openvpn:x:977:
nm-openvpn:x:976:
minetest:x:975:
abrt:x:974:
geoclue:x:973:
ntp:x:123:
`))
		})

		It("Adds and deletes an entry", func() {
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/group/group", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/group/group.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Group).Name).Should(Equal("foo"))

			entity.Apply(tmpFile.Name(), false)

			dat, err := ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`nm-openconnect:x:979:
sddm:x:978:
openvpn:x:977:
nm-openvpn:x:976:
minetest:x:975:
abrt:x:974:
geoclue:x:973:
ntp:x:123:
foo:xx:1:one,two,tree
`))

			entity, err = p.ReadEntity("../../testing/fixtures/group/group_add.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Group).Name).Should(Equal("foo"))

			entity.Apply(tmpFile.Name(), false)

			dat, err = ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`nm-openconnect:x:979:
sddm:x:978:
openvpn:x:977:
nm-openvpn:x:976:
minetest:x:975:
abrt:x:974:
geoclue:x:973:
ntp:x:123:
foo:xx:1:one,two,tree,four
`))

			entity.Delete(tmpFile.Name())
			dat, err = ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`nm-openconnect:x:979:
sddm:x:978:
openvpn:x:977:
nm-openvpn:x:976:
minetest:x:975:
abrt:x:974:
geoclue:x:973:
ntp:x:123:
`))
		})
	})
})
