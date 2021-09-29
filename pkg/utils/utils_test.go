package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/RasaHQ/rasactl/pkg/utils"
)

var _ = Describe("Utils", func() {

	var hostname string = "unittest.rasactl.localhost"

	It("check if the ls command exists - CommandExists", func() {
		cmdExists := utils.CommandExists("ls")
		Expect(cmdExists).To(Equal(true))
	})

	It("add a hostname to /etc/hosts - AddHostToEtcHosts", func() {
		err := utils.AddHostToEtcHosts(hostname, "127.0.0.1")
		Expect(err).To(BeNil())
	})

	It("delete a hostname from /etc/hosts - DeleteHostToEtcHosts", func() {
		err := utils.DeleteHostToEtcHosts(hostname)
		Expect(err).To(BeNil())
	})

	It("check if debug or verbose mode is enabled", func() {
		IsDebugOrVerboseEnabled := utils.IsDebugOrVerboseEnabled()
		Expect(IsDebugOrVerboseEnabled).To(Equal(false))
	})

	It("check if URL is accessible", func() {
		IsURLAccessible := utils.IsURLAccessible("https://github.com")
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

})
