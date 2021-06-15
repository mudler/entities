// Copyright © 2021 Daniele Rondina <geaaru@sabayonlinux.org>
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
	//"fmt"
	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store Tests", func() {
	Context("Loading entities via yaml", func() {

		It("Check Store Load", func() {
			store1 := NewEntitiesStore()
			err := store1.Load("../../testing/fixtures")
			Expect(err).Should(BeNil())
			Expect(len(store1.Users)).Should(Equal(2))
			Expect(len(store1.Groups)).Should(Equal(2))
			Expect(len(store1.Shadows)).Should(Equal(2))
			Expect(len(store1.GShadows)).Should(Equal(2))
			Expect(len(store1.Groups["foo"].GetUsers())).Should(Equal(4))
		})

		It("Check Store Merge", func() {
			store2 := NewEntitiesStore()
			err := store2.Load("../../testing/fixtures")
			Expect(err).Should(BeNil())
			Expect(len(store2.Users)).Should(Equal(2))
			Expect(len(store2.Shadows)).Should(Equal(2))
			Expect(len(store2.GShadows)).Should(Equal(2))

			// Check merge
			gid := 1
			err = store2.AddEntity(Group{
				Name:     "foo",
				Password: "yy",
				Gid:      &gid,
				Users:    "one,five",
			})
			Expect(err).Should(BeNil())
			Expect(len(store2.Groups)).Should(Equal(2))
			Expect(len(store2.Groups["foo"].GetUsers())).Should(Equal(5))
			Expect(store2.Groups["foo"].Password).Should(Equal("xx"))
		})

	})
})
