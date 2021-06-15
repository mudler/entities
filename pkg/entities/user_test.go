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

var _ = Describe("User", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
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

			dat, err := ioutil.ReadFile(tmpFile.Name())
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
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
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

			dat, err := ioutil.ReadFile(tmpFile.Name())
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
			dat, err = ioutil.ReadFile(tmpFile.Name())
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
			tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			expectedMap := map[string]UserPasswd{
				"root": UserPasswd{
					Username: "root",
					Password: "x",
					Uid:      0,
					Gid:      0,
					Group:    "",
					Info:     "Foo!",
					Homedir:  "/home/foo",
					Shell:    "/bin/bash",
				},
				"brokenuid": UserPasswd{
					Username: "brokenuid",
					Password: "x",
					Uid:      0,
					Gid:      100,
					Group:    "",
					Info:     "group",
					Homedir:  "/home/broken",
					Shell:    "/bin/bash",
				},
				"brokengid": UserPasswd{
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

	})
})
