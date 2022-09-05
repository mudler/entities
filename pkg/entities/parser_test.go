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
		It("understands the gshadow kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/gshadow/gshadow.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(GShadow).Name).Should(Equal("test"))
		})
	})
})
