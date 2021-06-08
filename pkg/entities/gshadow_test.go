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

package entities_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GShadow", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/gshadow/gshadow", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/gshadow/update.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(GShadow).Name).Should(Equal("postmaster"))

			err = entity.Apply(tmpFile.Name(), false)
			Expect(err).Should(BeNil())

			dat, err := ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`systemd-bus-proxy:!::
systemd-coredump:!::
systemd-journal-gateway:!::
systemd-journal-remote:!::
systemd-journal-upload:!::
systemd-network:!::
systemd-resolve:!::
systemd-timesync:!::
netdev:!::
avahi:!::
avahi-autoipd:!::
mail:!::
postmaster:foo:barred:baz
ldap:!::
`))
		})

		It("Adds and deletes an entry", func() {
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/gshadow/gshadow", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/gshadow/gshadow.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(GShadow).Name).Should(Equal("test"))

			entity.Apply(tmpFile.Name(), false)

			dat, err := ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`systemd-bus-proxy:!::
systemd-coredump:!::
systemd-journal-gateway:!::
systemd-journal-remote:!::
systemd-journal-upload:!::
systemd-network:!::
systemd-resolve:!::
systemd-timesync:!::
netdev:!::
avahi:!::
avahi-autoipd:!::
mail:!::
postmaster:!::
ldap:!::
test:!:foo,bar:foo,baz
`))

			entity.Delete(tmpFile.Name())
			dat, err = ioutil.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`systemd-bus-proxy:!::
systemd-coredump:!::
systemd-journal-gateway:!::
systemd-journal-remote:!::
systemd-journal-upload:!::
systemd-network:!::
systemd-resolve:!::
systemd-timesync:!::
netdev:!::
avahi:!::
avahi-autoipd:!::
mail:!::
postmaster:!::
ldap:!::
`))
		})
	})
})
