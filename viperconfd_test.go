package config_test

import (
	"os"

	"github.com/na4ma4/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("ViperConfD test", func() {
	It("is thread-safe", func() {
		v := config.NewViperConfD("test", "test/conf.d")

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

		stringTest := v.GetString("test.string")

		Expect(stringTest).To(Equal("string"))
	})

	It("loading from file", func() {
		v := config.NewViperConfD("test-project", "test/conf.d")

		stringTest := v.GetString("category1.string")
		Expect(stringTest).To(Equal("foobar"))
	})

	It("loading from specified file", func() {
		v := config.NewViperConfD("test", "test/conf.d", "test/test-project.toml")

		stringTest := v.GetString("category1.string")
		Expect(stringTest).To(Equal("foobar"))
	})

	It("writing to a file", func() {
		tempfile, err := os.CreateTemp("", "*-dummy-file.toml")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tempfile.Name())

		v := config.NewViperConfD("test", "", tempfile.Name())

		v.SetString("category.test", "barfoo")

		err = v.Save()
		Expect(err).NotTo(HaveOccurred())

		b, err := os.ReadFile(tempfile.Name())
		Expect(err).NotTo(HaveOccurred())

		expectedOutput := "[category]\ntest = 'barfoo'\n"

		Expect(string(b)).To(Equal(expectedOutput))
	})

	It("importing system viper", func() {
		viper.SetDefault("sesame.open", "open.sesame")
		viper.SetDefault("system.test.duration", "30s")
		viper.Set("fooman", "barwoman")

		viper.SetDefault("system.default", "default")

		v := config.NewViperConfDFromViper(viper.GetViper(), "test/conf.d")

		v.Set("system.default", "override")

		Expect(v.GetString("sesame.open")).To(Equal("open.sesame"))
		Expect(v.GetString("fooman")).To(Equal("barwoman"))
		Expect(v.GetDuration("system.test.duration").String()).To(Equal("30s"))
		Expect(v.GetString("system.default")).To(Equal("override"))
	})

	It("can set a default value", func() {
		v := config.NewViperConfDFromViper(viper.GetViper(), "test/conf.d")

		if vp, ok := v.(*config.ViperConfD); ok {
			vp.SetDefault("some-key-with-default", "custom-default-value")
		}

		Expect(v.GetString("some-key-with-default")).To(Equal("custom-default-value"))

		v.Set("some-key-with-default", "new-value")

		Expect(v.GetString("some-key-with-default")).To(Equal("new-value"))
	})

	It("will load the conf.d files", func() {
		v := config.NewViperConfDFromViper(viper.GetViper(), "test/conf.d")

		Expect(v.GetInt("category2.int")).To(Equal(8335))
	})

	It("will load in order and overwrite specified options", func() {
		v := config.NewViperConfDFromViper(viper.GetViper(), "test/conf.d")

		Expect(v.GetString("category3.first")).To(Equal("foo"))
		Expect(v.GetString("category3.second")).To(Equal("foobar"))
	})
})
