package internal_test

import (
	"os"

	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/FooSoft/homemaker/internal"
)

var _ = Describe("Homemaker Smoke Test", func() {
	Describe("Executing homemaker run default", func() {
		var (
			err error
		)
		BeforeEach(func() {

		})
		It("Should Successfully create and delete the sample.conf link", func() {
			By("Loading configuration")
			viper.SetConfigFile("./../test/homemaker.yml")
			err = viper.ReadInConfig()
			Expect(err).To(BeNil())

			By("Unmarshalling conf to structure")
			c := &Config{}
			err = viper.Unmarshal(c)
			c.SrcDir = "./../test/src"
			c.DstDir = "./../test/dst"
			Expect(err).To(BeNil())

			By("Processing task")
			err = ProcessTask("simple-link", c)
			Expect(err).To(BeNil())

			By("Checking cretedFile exist")
			createdFile := c.DstDir + "/sample.conf"
			_, err = os.Stat(createdFile)
			Expect(err).To(BeNil())

			By("Exexute same task with unlink")
			c.Unlink = true
			err = ProcessTask("simple-link", c)
			Expect(err).To(BeNil())

			By("Checking cretedFile exist")
			_, err = os.Stat(createdFile)
			Expect(err).ToNot(BeNil())
		})
	})
})
