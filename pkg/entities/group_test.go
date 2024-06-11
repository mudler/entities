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
	"path/filepath"
	"sync"

	"github.com/gofrs/flock"
	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Group", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
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

			dat, err := os.ReadFile(tmpFile.Name())
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
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
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

			dat, err := os.ReadFile(tmpFile.Name())
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

			dat, err = os.ReadFile(tmpFile.Name())
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
			dat, err = os.ReadFile(tmpFile.Name())
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

		It("works with locks", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
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

			baseName := filepath.Base(tmpFile.Name())
			fileLock := flock.New(fmt.Sprintf("/var/lock/%s.lock", baseName))
			defer os.Remove(fileLock.Path())
			locked, err := fileLock.TryLock()
			Expect(err).To(BeNil())
			Expect(locked).To(BeTrue())
			defer fileLock.Close()

			err = entity.Apply(tmpFile.Name(), false)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Failed locking file"))
		})
	})
})
