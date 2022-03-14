package utils_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/RasaHQ/rasactl/pkg/utils"
)

var _ = Describe("Utils", func() {

	It("check if the ls command exists - CommandExists", func() {
		cmdExists := utils.CommandExists("ls")
		Expect(cmdExists).To(Equal(true))
	})

	It("check if debug or verbose mode is enabled", func() {
		IsDebugOrVerboseEnabled := utils.IsDebugOrVerboseEnabled()
		Expect(IsDebugOrVerboseEnabled).To(Equal(false))
	})

	It("check if URL is accessible", func() {
		IsURLAccessible := utils.IsURLAccessible("https://google.com")
		Expect(IsURLAccessible).To(Equal(true))
	})

	It("merge maps", func() {

		map1 := map[string]interface{}{"key1": "value1"}
		map2 := map[string]interface{}{"key2": "value2"}

		expectedMap := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		mergedMaps := utils.MergeMaps(map1, map2)
		Expect(mergedMaps).To(Equal(expectedMap))
	})

	It("check version constrains", func() {
		version := utils.RasaXVersionConstrains("1.0.0", "< 1.0.0")
		Expect(version).To(Equal(false))
	})

	It("validate namespace name", func() {
		err := utils.ValidateName("name!-space")
		Expect(err).To(Not(BeNil()))
	})

	It("convert string slice to JSON", func() {
		d := [][]string{{"test:", "test"}}
		str, _ := utils.StringSliceToJSON(d)
		Expect(str).To(Equal("{\"test\":\"test\"}"))
	})

	It("check password generator", func() {
		length := 14

		password, err := utils.GenerateRandomPassword(length)
		Expect(password).To(HaveLen(length))
		Expect(password).To(BeAssignableToTypeOf("string"))
		Expect(err).To(BeNil())
	})

	Describe("read Rasa X URL from environment variables", func() {
		viper.AutomaticEnv() // read in environment variables that match
		viper.SetEnvPrefix("rasactl")

		Context("no variable set", func() {
			It("URL should be empty", func() {
				url := utils.GetRasaXURLEnv("my-deployment")

				Expect(url).To(BeEmpty())
			})
		})

		Context("set environment variables", func() {
			It("set the global env variable - RASACTL_RASA_X_URL", func() {
				os.Setenv("RASACTL_RASA_X_URL", "http://test.localhost")
				url := utils.GetRasaXURLEnv("my-deployment")

				Expect(url).To(Equal("http://test.localhost"))
			})

			It("set variable for a specific namespace - RASACTL_RASA_X_URL_MY_DEPLOYMENT", func() {
				os.Setenv("RASACTL_RASA_X_URL_MY_DEPLOYMENT", "http://my-deployment.test.localhost")
				url := utils.GetRasaXURLEnv("my-deployment")

				Expect(url).To(Equal("http://my-deployment.test.localhost"))
			})

			It("check if env variable for a namespace override the global variable", func() {
				os.Setenv("RASACTL_RASA_X_URL", "http://test.localhost")
				os.Setenv("RASACTL_RASA_X_URL_MY_DEPLOYMENT", "http://my-deployment.test.localhost")
				url := utils.GetRasaXURLEnv("my-deployment")

				Expect(url).To(Equal("http://my-deployment.test.localhost"))
			})

			It("check if the global variable is used in a case that an env variable for a namespace is not set", func() {
				os.Setenv("RASACTL_RASA_X_URL", "http://test.localhost")
				os.Setenv("RASACTL_RASA_X_URL_MY_DEPLOYMENT", "http://my-deployment.test.localhost")
				url := utils.GetRasaXURLEnv("my-deployment-1")

				Expect(url).To(Equal("http://test.localhost"))
			})
		})
	})
})
