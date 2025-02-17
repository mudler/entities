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

	"github.com/gofrs/flock"
	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/simple/passwd", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/simple/update.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(UserPasswd).Username).Should(Equal("root"))

			err = entity.Apply(tmpFile.Name(), false)
			Expect(err).Should(BeNil())

			dat, err := os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`root:x:0:0:Foo!:/home/foo:/bin/bash
bin:x:1:1:bin:/bin:/bin/false
daemon:x:2:2:daemon:/sbin:/bin/false
adm:x:3:4:adm:/var/adm:/bin/false
lp:x:4:7:lp:/var/spool/lpd:/bin/false
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
unbound:x:999:955:added by portage for unbound:/etc/unbound:/sbin/nologin
gpsd:x:139:14:added by portage for gpsd:/dev/null:/sbin/nologin
`))
		})

		It("Adds and deletes an entry", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/simple/passwd", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/simple/user.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(UserPasswd).Username).Should(Equal("foo"))

			err = entity.Apply(tmpFile.Name(), false)
			Expect(err).Should(BeNil())

			dat, err := os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/bin/false
daemon:x:2:2:daemon:/sbin:/bin/false
adm:x:3:4:adm:/var/adm:/bin/false
lp:x:4:7:lp:/var/spool/lpd:/bin/false
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
unbound:x:999:955:added by portage for unbound:/etc/unbound:/sbin/nologin
gpsd:x:139:14:added by portage for gpsd:/dev/null:/sbin/nologin
foo:pass:0:0:Foo!:/home/foo:/bin/bash
`))

			entity.Delete(tmpFile.Name())
			dat, err = os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/bin/false
daemon:x:2:2:daemon:/sbin:/bin/false
adm:x:3:4:adm:/var/adm:/bin/false
lp:x:4:7:lp:/var/spool/lpd:/bin/false
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
unbound:x:999:955:added by portage for unbound:/etc/unbound:/sbin/nologin
gpsd:x:139:14:added by portage for gpsd:/dev/null:/sbin/nologin
`))
		})

		It("Read broken file", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			expectedMap := map[string]UserPasswd{
				"root": {
					Username: "root",
					Password: "x",
					Uid:      0,
					Gid:      0,
					Group:    "",
					Info:     "Foo!",
					Homedir:  "/home/foo",
					Shell:    "/bin/bash",
				},
				"brokenuid": {
					Username: "brokenuid",
					Password: "x",
					Uid:      0,
					Gid:      100,
					Group:    "",
					Info:     "group",
					Homedir:  "/home/broken",
					Shell:    "/bin/bash",
				},
				"brokengid": {
					Username: "brokengid",
					Password: "x",
					Uid:      100,
					Gid:      100,
					Group:    "",
					Info:     "group",
					Homedir:  "/home/broken",
					Shell:    "/bin/bash",
				},
			}

			dat := `root:x:0:0:Foo!:/home/foo:/bin/bash
brokenuid:x::100:group:/home/broken:/bin/bash
brokengid:x:100::group:/home/broken:/bin/bash
`

			tmpFile.WriteString(dat)
			tmpFile.Close()

			m, err := ParseUser(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(m).Should(Equal(expectedMap))

		})

		It("Works with locks", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/simple/passwd", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/simple/update.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(UserPasswd).Username).Should(Equal("root"))

			baseName := filepath.Base(tmpFile.Name())
			fileLock := flock.New(fmt.Sprintf("/run/lock/%s.lock", baseName))
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
