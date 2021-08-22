package cmd

import (
	"testing"

	"github.com/RasaHQ/rasactl/pkg/rasactl"
)

type mockKubernetesClient interface {
}

func TestParseArgs(t *testing.T) {
	t.Parallel()

	var k8sClient mockKubernetesClient

	initLog()
	initConfig()
	namespace = ""

	rasaCtl = &rasactl.RasaCtl{
		Log:              log,
		Flags:            rasactlFlags,
		KubernetesClient: &k8sClient,
	}

	// rasactl command
	args := []string{}
	parseArgs(args, 0, 0)

	// rasactl command [deployment]
	args = []string{}
	namespace = "ns-test"
	parseArgs(args, 1, 1)

	// rasactl command arg1
	args = []string{"args"}
	parseArgs(args, 1, 1)

}
