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

package entities_test

import (
	"fmt"
	"os"

	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GShadow", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
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

			dat, err := os.ReadFile(tmpFile.Name())
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
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
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

			dat, err := os.ReadFile(tmpFile.Name())
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
			dat, err = os.ReadFile(tmpFile.Name())
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
