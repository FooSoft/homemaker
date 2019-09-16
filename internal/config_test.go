package internal_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/FooSoft/homemaker/internal"
)

var _ = Describe("Config", func() {
	It("Should set the environment", func() {
		c := Config{
			File: "file",
			//SelectedTask: "selectedTask",
			SrcDir:  "srcDir",
			DstDir:  "dstDir",
			Variant: "variant",
		}

		c.SetEnv()

		Expect(os.Getenv("HM_CONFIG")).To(Equal(c.File))
		//Expect(os.Getenv("HM_TASK")).To(Equal(c.SelectedTask))
		Expect(os.Getenv("HM_SRC")).To(Equal(c.SrcDir))
		Expect(os.Getenv("HM_DEST")).To(Equal(c.DstDir))
		Expect(os.Getenv("HM_VARIANT")).To(Equal(c.Variant))
	})

})
