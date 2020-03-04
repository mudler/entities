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
	. "github.com/mudler/entities/pkg/entities"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}
		It("understands the user kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/simple/user.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(UserPasswd).Username).Should(Equal("foo"))
		})
		It("understands the shadow kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/shadow/user.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Shadow).Username).Should(Equal("foo"))
		})
		It("understands the group kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/group/group.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Group).Name).Should(Equal("foo"))
		})
	})
})
