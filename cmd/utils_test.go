package cmd

import (
	"testing"

	// mock

	"github.com/RasaHQ/rasactl/pkg/k8s/fake"
	"github.com/RasaHQ/rasactl/pkg/rasactl"
	"github.com/golang/mock/gomock"
)

func TestParseArgs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	initLog()
	initConfig()
	getNamespace()

	m := fake.NewMockKubernetesInterface(ctrl)
	m.EXPECT().GetNamespaces().Return([]string{}, nil)
	m.EXPECT().SetNamespace(gomock.Any())

	rasaCtl = &rasactl.RasaCtl{
		Log:              log,
		Flags:            rasactlFlags,
		KubernetesClient: m,
	}

	args := []string{}
	namespace = ""
	parseArgs(args, 0, 0)

}
