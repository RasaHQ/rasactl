package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	fd "github.com/RasaHQ/rasactl/pkg/docker/fake"
	fh "github.com/RasaHQ/rasactl/pkg/helm/fake"
	fk "github.com/RasaHQ/rasactl/pkg/k8s/fake"
	"github.com/RasaHQ/rasactl/pkg/rasactl"
	"github.com/RasaHQ/rasactl/pkg/types"
)

func TestParseArgs(t *testing.T) {

	initLog()
	initConfig()

	testCases := []struct {
		name                string
		namespace           string
		expectedNamespace   string
		minArgs             int
		maxArgs             int
		args                []string
		getNamespacesReturn []string
		expectedArgs        []string
		flags               *types.RasaCtlFlags
		arg0IsNs            bool
	}{
		{
			"no arguments", // name
			"",             // namespace
			"",             // expectedNamespace
			0,              // minArgs
			0,              // maxArgs
			[]string{},     // args
			[]string{},     // getNamespacesReturn
			[]string{},     // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, rasactl start", // name
			"",                              // namespace
			"",                              // expectedNamespace
			1,                               // minArgs
			1,                               // maxArgs
			[]string{},                      // args
			[]string{"test-deployment"},     // getNamespacesReturn
			[]string{"test-deployment"},     // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"only one arg, rasactl command deployment-name", // name
			"",                          // namespace
			"",                          // expectedNamespace
			1,                           // minArgs
			1,                           // maxArgs
			[]string{"test-deployment"}, // args
			[]string{""},                // getNamespacesReturn
			[]string{"test-deployment"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"no deployment, one arg, rasactl command [deployment-name] arg1", // name
			"",                   // namespace
			"",                   // expectedNamespace
			1,                    // minArgs
			2,                    // maxArgs
			[]string{"arg1"},     // args
			[]string{},           // getNamespacesReturn
			[]string{"", "arg1"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, one arg, rasactl command [deployment-name] arg1", // name
			"",                                  // namespace
			"",                                  // expectedNamespace
			1,                                   // minArgs
			2,                                   // maxArgs
			[]string{"arg1"},                    // args
			[]string{"test-deployment"},         // getNamespacesReturn
			[]string{"test-deployment", "arg1"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, defined deployment, one arg, rasactl command deployment-name arg1", // name
			"",                                     // namespace
			"",                                     // expectedNamespace
			1,                                      // minArgs
			2,                                      // maxArgs
			[]string{"defined-deployment", "arg1"}, // args
			[]string{"test-deployment"},            // getNamespacesReturn
			[]string{"defined-deployment", "arg1"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, two args, rasactl command deployment-name arg1 [arg2]", // name
			"",                                     // namespace
			"",                                     // expectedNamespace
			1,                                      // minArgs
			3,                                      // maxArgs
			[]string{"defined-deployment", "arg1"}, // args
			[]string{"test-deployment"},            // getNamespacesReturn
			[]string{"test-deployment", "defined-deployment", "arg1"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, defined deployment, two args, rasactl command deployment-name arg1 [arg2]", // name
			"",                                     // namespace
			"",                                     // expectedNamespace
			1,                                      // minArgs
			3,                                      // maxArgs
			[]string{"defined-deployment", "arg1"}, // args
			[]string{"test-deployment"},            // getNamespacesReturn
			[]string{"defined-deployment", "arg1", ""}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			true, // arg0IsNs
		},
		{
			"one deployment, defined deployment, three args, rasactl command deployment-name arg1 arg2", // name
			"", // namespace
			"", // expectedNamespace
			1,  // minArgs
			3,  // maxArgs
			[]string{"defined-deployment", "arg1", "arg2"}, // args
			[]string{"test-deployment"},                    // getNamespacesReturn
			[]string{"defined-deployment", "arg1", "arg2"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, two args, rasactl command [deployment-name] arg1 arg2", // name
			"",                          // namespace
			"",                          // expectedNamespace
			2,                           // minArgs
			3,                           // maxArgs
			[]string{"arg1", "arg2"},    // args
			[]string{"test-deployment"}, // getNamespacesReturn
			[]string{"test-deployment", "arg1", "arg2"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, new deployment, rasactl start deployment --create", // name
			"",                          // namespace
			"",                          // expectedNamespace
			1,                           // minArgs
			1,                           // maxArgs
			[]string{"deployment"},      // args
			[]string{"test-deployment"}, // getNamespacesReturn
			[]string{"deployment"},      // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      true,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, new deployment, rasactl start deployment --project", // name
			"",                          // namespace
			"",                          // expectedNamespace
			1,                           // minArgs
			1,                           // maxArgs
			[]string{"deployment"},      // args
			[]string{"test-deployment"}, // getNamespacesReturn
			[]string{"deployment"},      // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     true,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"one deployment, new deployment, rasactl start deployment --project-path=test", // name
			"",                          // namespace
			"",                          // expectedNamespace
			1,                           // minArgs
			1,                           // maxArgs
			[]string{"deployment"},      // args
			[]string{"test-deployment"}, // getNamespacesReturn
			[]string{"deployment"},      // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "test",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"two deployment, one arg, rasactl command deployment", // name
			"",                     // namespace
			"deployment",           // expectedNamespace
			1,                      // minArgs
			1,                      // maxArgs
			[]string{"deployment"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"deployment"},                           // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			true, // arg0IsNs
		},
		{
			"two deployment, one arg rasactl command [deployment]", // name
			"",         // namespace
			"",         // expectedNamespace
			1,          // minArgs
			1,          // maxArgs
			[]string{}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{""}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"defined deployment, two deployment, one arg rasactl command [deployment]", // name
			"defined-deployment", // namespace
			"",                   // expectedNamespace
			1,                    // minArgs
			1,                    // maxArgs
			[]string{},           // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"defined-deployment"},                   // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"defined deployment, two deployment, two args rasactl command [deployment] arg1 arg2", // name
			"defined-deployment",     // namespace
			"",                       // expectedNamespace
			2,                        // minArgs
			3,                        // maxArgs
			[]string{"arg1", "arg2"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"defined-deployment"},                   // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"defined deployment, two deployment, two args rasactl command [deployment] arg1 [arg2]", // name
			"defined-deployment", // namespace
			"",                   // expectedNamespace
			1,                    // minArgs
			3,                    // maxArgs
			[]string{"arg1"},     // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"defined-deployment", "arg1"},           // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"defined deployment, two deployment, two args required,  rasactl command deployment arg1", // name
			"defined-deployment",           // namespace
			"",                             // expectedNamespace
			1,                              // minArgs
			2,                              // maxArgs
			[]string{"deployment", "arg1"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"deployment", "arg1"},                   // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			true, // arg0IsNs
		},
		{
			"defined deployment, two deployments, rasactl command deployment arg1", // name
			"defined-deployment",           // namespace
			"",                             // expectedNamespace
			1,                              // minArgs
			2,                              // maxArgs
			[]string{"deployment", "arg1"}, // args
			[]string{"test-deployment1", "test-deployment2"},     // getNamespacesReturn
			[]string{"defined-deployment", "deployment", "arg1"}, // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"two deployments, rasactl command deployment", // name
			"",                     // namespace
			"",                     // expectedNamespace
			0,                      // minArgs
			0,                      // maxArgs
			[]string{"deployment"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"deployment"},                           // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"two deployments, rasactl command deployment arg1 arg2", // name
			"",                                     // namespace
			"",                                     // expectedNamespace
			2,                                      // minArgs
			3,                                      // maxArgs
			[]string{"deployment", "arg1", "arg2"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"deployment", "arg1", "arg2"},           // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"two deployments, rasactl command deployment arg1 [arg2]", // name
			"",                             // namespace
			"",                             // expectedNamespace
			1,                              // minArgs
			3,                              // maxArgs
			[]string{"deployment", "arg1"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"deployment", "arg1", ""},               // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"two deployments, rasactl command [deployment] arg1 [arg2]", // name
			"",               // namespace
			"",               // expectedNamespace
			1,                // minArgs
			3,                // maxArgs
			[]string{"arg1"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"", "arg1", ""},                         // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
		{
			"detected deployment, two deployments, rasactl command [deployment] arg1 [arg2]", // name
			"test-ns",        // namespace
			"test-ns",        // expectedNamespace
			1,                // minArgs
			3,                // maxArgs
			[]string{"arg1"}, // args
			[]string{"test-deployment1", "test-deployment2"}, // getNamespacesReturn
			[]string{"test-ns", "arg1"},                      // expectedArgs
			&types.RasaCtlFlags{
				Start: types.RasaCtlStartFlags{
					Create:      false,
					ProjectPath: "",
					Project:     false,
				},
			}, // flags
			false, // arg0IsNs
		},
	}

	for _, testCase := range testCases {
		// Here, we capture the range variable and force it into the scope of this block. If we don't do this, when the
		// subtest switches contexts (because of t.Parallel), the testCase value will have been updated by the for loop
		// and will be the next testCase!
		testCase := testCase

		t.Run(testCase.name, func(subT *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mk := fk.NewMockKubernetesInterface(ctrl)
			md := fd.NewMockInterface(ctrl)
			mh := fh.NewMockInterface(ctrl)

			mk.EXPECT().GetNamespaces().Return(testCase.getNamespacesReturn, nil)
			mk.EXPECT().IsNamespaceExist(gomock.Any()).Return(testCase.arg0IsNs, nil).AnyTimes()

			switch {
			case testCase.expectedNamespace == "":
				mk.EXPECT().SetNamespace(gomock.Any()).AnyTimes()
				md.EXPECT().SetNamespace(gomock.Any()).AnyTimes()
				mh.EXPECT().SetNamespace(gomock.Any()).AnyTimes()
			case testCase.expectedNamespace != "":
				mk.EXPECT().SetNamespace(gomock.Eq(testCase.expectedNamespace)).AnyTimes()
				md.EXPECT().SetNamespace(gomock.Eq(testCase.expectedNamespace)).AnyTimes()
				mh.EXPECT().SetNamespace(gomock.Eq(testCase.expectedNamespace)).AnyTimes()
			}

			rasaCtl = &rasactl.RasaCtl{
				Log:              log,
				Flags:            testCase.flags,
				KubernetesClient: mk,
				DockerClient:     md,
				HelmClient:       mh,
			}

			args, err := parseArgs(testCase.namespace, testCase.args, testCase.minArgs, testCase.maxArgs, testCase.flags)

			switch {
			case testCase.namespace == "" && len(testCase.getNamespacesReturn) == 0:
				require.NotEmpty(t, args)
			case testCase.namespace == "" && len(testCase.getNamespacesReturn) != 0:
				require.NotEmpty(t, args)
				require.Equal(t, testCase.expectedArgs, args)
			}
			require.NoError(t, err)
		})
	}

}
