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
	"io"
	"os"
	"time"

	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

var _ = Describe("Shadow", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}

		It("Changes an entry", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/shadow/shadow", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/shadow/update.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Shadow).Username).Should(Equal("halt"))

			err = entity.Apply(tmpFile.Name(), false)
			Expect(err).Should(BeNil())

			dat, err := os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(MatchRegexp(`halt\:\$.*\:1:2:3:4:5:6:`))
			Expect(string(dat)).To(ContainSubstring(
				`operator:*:9797:0:::::
shutdown:*:9797:0:::::
sync:*:9797:0:::::
bin:*:9797:0:::::
daemon:*:9797:0:::::
adm:*:9797:0:::::
lp:*:9797:0:::::
news:*:9797:0:::::
uucp:*:9797:0:::::
`))
		})

		It("Adds and deletes an entry", func() {
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			_, err = copy("../../testing/fixtures/shadow/shadow", tmpFile.Name())
			Expect(err).Should(BeNil())

			entity, err := p.ReadEntity("../../testing/fixtures/shadow/user.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Shadow).Username).Should(Equal("foo"))

			entity.Apply(tmpFile.Name(), false)

			dat, err := os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`halt:*:9797:0:::::
operator:*:9797:0:::::
shutdown:*:9797:0:::::
sync:*:9797:0:::::
bin:*:9797:0:::::
daemon:*:9797:0:::::
adm:*:9797:0:::::
lp:*:9797:0:::::
news:*:9797:0:::::
uucp:*:9797:0:::::
foo:$bar:1:2:3:4:5:6:
`))

			entity.Delete(tmpFile.Name())
			dat, err = os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal(
				`halt:*:9797:0:::::
operator:*:9797:0:::::
shutdown:*:9797:0:::::
sync:*:9797:0:::::
bin:*:9797:0:::::
daemon:*:9797:0:::::
adm:*:9797:0:::::
lp:*:9797:0:::::
news:*:9797:0:::::
uucp:*:9797:0:::::
`))
		})

	})

	It("test Prepare", func() {

		By("Giving a specific user", func() {
			t := time.Now()
			days := t.Unix() / 24 / 60 / 60
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			s1 := &Shadow{
				Username:    "user1",
				Password:    "$!",
				LastChanged: "now",
			}

			s1.Apply(tmpFile.Name(), false)

			dat, err := os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(Equal("user1:$!:" + fmt.Sprintf("%d", days) + "::::::\n"))

		})

		By("Giving a specific user", func() {
			t := time.Now()
			days := t.Unix() / 24 / 60 / 60
			tmpFile, err := os.CreateTemp(os.TempDir(), "pre-")
			if err != nil {
				fmt.Println("Cannot create temporary file", err)
			}

			// cleaning up by removing the file
			defer os.Remove(tmpFile.Name())

			s1 := &Shadow{
				Username:    "user1",
				Password:    "pass",
				LastChanged: "now",
			}

			s1.Apply(tmpFile.Name(), false)

			dat, err := os.ReadFile(tmpFile.Name())
			Expect(err).Should(BeNil())
			Expect(string(dat)).To(MatchRegexp(`user1\:\$.*\:` + fmt.Sprintf("%d", days) + ":.*"))

		})

	})

})
