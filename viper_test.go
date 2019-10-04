package config_test

import (
	"io/ioutil"
	"os"

	"github.com/na4ma4/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ViperConf test", func() {

	It("is thread-safe", func() {
		v := config.NewViperConfig("test")

		v.SetString("test.string", "string")

		go func() {
			for i := 0; i < 1000; i++ {
				v.SetString("test.string", "string")
			}
		}()

		go func() {
			for i := 0; i < 1000; i++ {
				v.GetString("test.string")
			}
		}()

		string_test := v.GetString("test.string")

		Expect(string_test).To(Equal("string"))
	})

	It("loading from file", func() {
		v := config.NewViperConfig("test-project")

		string_test := v.GetString("category1.string")
		Expect(string_test).To(Equal("foobar"))
	})

	It("loading from specified file", func() {
		v := config.NewViperConfig("test", "test/test-project.toml")

		string_test := v.GetString("category1.string")
		Expect(string_test).To(Equal("foobar"))
	})

	It("writing to a file", func() {
		tempfile, err := ioutil.TempFile("", "*-dummy-file.toml")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tempfile.Name())

		v := config.NewViperConfig("test", tempfile.Name())

		v.SetString("category.test", "barfoo")

		err = v.Save()
		Expect(err).NotTo(HaveOccurred())

		b, err := ioutil.ReadFile(tempfile.Name())
		Expect(err).NotTo(HaveOccurred())

		expectedOutput := "\n[category]\n  test = \"barfoo\"\n"

		Expect(string(b)).To(Equal(expectedOutput))
	})

})
