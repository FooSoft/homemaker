package internal

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		c Config
	)
	Describe("Setting the environment", func() {
		BeforeEach(func() {
			c = Config{
				File:    "file",
				SrcDir:  "srcDir",
				DstDir:  "dstDir",
				Variant: "variant",
			}

			c.setEnv()
		})
		It("Should set HM_CONFIG", func() {
			Expect(os.Getenv("HM_CONFIG")).To(Equal(c.File))
		})
		It("Should set HM_SRC", func() {
			Expect(os.Getenv("HM_SRC")).To(Equal(c.SrcDir))
		})
		It("Should set HM_DEST", func() {
			Expect(os.Getenv("HM_DEST")).To(Equal(c.DstDir))
		})
		It("Should set HM_VARIANT", func() {
			Expect(os.Getenv("HM_VARIANT")).To(Equal(c.Variant))
		})
	})
	Describe("When digesting the configuration", func() {
		BeforeEach(func() {
			c = Config{
				File:    "file",
				SrcDir:  "./../test/src",
				DstDir:  "./../test/dst",
				Variant: "variant",
			}
		})
		JustBeforeEach(func() {
			c.digest()
		})

		Specify("that SrcDir must be an absolute path and point to a valid directory", func() {
			Expect(c.SrcDir).To(BeADirectory())
		})
		Specify("that DstDir must be and absolute path and point to a valid directory", func() {
			Expect(c.DstDir).To(BeADirectory())
		})
		Describe("If Unlink is set", func() {
			BeforeEach(func() {
				c.Unlink = true
			})
			Specify("that Nocmd must be set", func() {
				Expect(c.Nocmds).To(BeTrue())
			})
		})
	})

})
