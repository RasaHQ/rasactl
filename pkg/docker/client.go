package docker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
)

type Docker struct {
	Client    *client.Client
	ctx       context.Context
	Namespace string
}

const (
	kindImage   string = "kindest/node:v1.20.2"
	kindNetwork string = "kind"
)

func (d *Docker) New() error {
	d.ctx = context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	d.Client = cli
	return nil
}

func (d *Docker) prepareKindJoinConfiguration() (string, error) {

	token, err := d.getKubeadmToken()
	if err != nil {
		return "", err
	}

	file := fmt.Sprintf("/tmp/rasaxctl-kind-joinconfig-%s.yaml", d.Namespace)

	joinConfiguration := bootstrapv1.JoinConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "JoinConfiguration",
			APIVersion: "kubeadm.k8s.io/v1beta2",
		},
		Discovery: bootstrapv1.Discovery{
			BootstrapToken: &bootstrapv1.BootstrapTokenDiscovery{
				Token:                    token,
				APIServerEndpoint:        "kind-control-plane:6443",
				UnsafeSkipCAVerification: true,
			},
			Timeout:           &metav1.Duration{time.Minute * 5},
			TLSBootstrapToken: token,
		},
		NodeRegistration: bootstrapv1.NodeRegistrationOptions{
			CRISocket:        "unix:///run/containerd/containerd.sock",
			KubeletExtraArgs: map[string]string{"fail-swap-on": "false"},
			Name:             fmt.Sprintf("kind-%s", d.Namespace),
			Taints: []v1.Taint{{
				Key:    "rasaxctl-project",
				Value:  d.Namespace,
				Effect: v1.TaintEffectNoSchedule,
			}},
		},
	}

	config, err := yaml.Marshal(joinConfiguration)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file, config, 0644); err != nil {
		return "", nil
	}
	return file, nil
}

func (d *Docker) copyJoinConfigurationToContainer(container container.ContainerCreateCreatedBody) error {
	joinConfig, err := d.prepareKindJoinConfiguration()
	if err != nil {
		return nil
	}
	srcInfo, err := archive.CopyInfoSourcePath(joinConfig, true)
	if err != nil {
		return err
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer srcArchive.Close()
	dstInfo := archive.CopyInfo{
		Path:  "/config.yaml",
		IsDir: false,
	}

	dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer preparedArchive.Close()

	if err := d.Client.CopyToContainer(d.ctx, container.ID, dstDir, preparedArchive, types.CopyToContainerOptions{}); err != nil {
		return err
	}

	os.Remove(joinConfig)

	return nil
}

func (d *Docker) getKubeadmToken() (string, error) {
	container := "kind-control-plane"
	token := new(bytes.Buffer)
	execSpec, err := d.Client.ContainerExecCreate(d.ctx, container, types.ExecConfig{
		WorkingDir:   "/",
		Cmd:          []string{"kubeadm", "token", "create", "--ttl", "180s"},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", err
	}
	at, err := d.Client.ContainerExecAttach(d.ctx, execSpec.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	token.ReadFrom(at.Reader)
	at.Close()

	r := strings.TrimSuffix(string(bytes.Trim(token.Bytes(), "\x01\x00\x00\x00\x00\x00\x00\x18")), "\n")

	return r, nil
}

func (d *Docker) joinKindNodeToKubernetesCluster(container container.ContainerCreateCreatedBody) error {
	execSpec, err := d.Client.ContainerExecCreate(d.ctx, container.ID, types.ExecConfig{
		WorkingDir: "/",
		Cmd:        []string{"kubeadm", "join", "--config", "config.yaml", "--skip-phases=preflight", "-v", "6"},
		//AttachStdout: true,
		//AttachStderr: true,
	})
	if err != nil {
		return err
	}
	if err := d.Client.ContainerExecStart(d.ctx, execSpec.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	for {
		status, err := d.Client.ContainerExecInspect(d.ctx, execSpec.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println(status.ExitCode, status.Running)
		if !status.Running {
			break
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

func (d *Docker) CreateKindNode(hostname string) (container.ContainerCreateCreatedBody, error) {
	d.Client.ImagePull(d.ctx, kindImage, types.ImagePullOptions{})

	resp, err := d.Client.ContainerCreate(d.ctx,
		&container.Config{
			Image:    kindImage,
			Tty:      false,
			Hostname: hostname,
			Volumes:  map[string]struct{}{"/var": {}},
		},
		&container.HostConfig{
			Privileged:  true,
			SecurityOpt: []string{"apparmor=unconfined", "seccomp=unconfined"},
			Tmpfs:       map[string]string{"/run": "", "/tmp": ""},
		}, nil, hostname)
	if err != nil {
		return resp, err
	}
	if err := d.Client.NetworkConnect(d.ctx, kindNetwork, resp.ID, nil); err != nil {
		return resp, err
	}

	if err := d.Client.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return resp, err
	}

	if err := d.copyJoinConfigurationToContainer(resp); err != nil {
		return resp, nil
	}

	if err := d.joinKindNodeToKubernetesCluster(resp); err != nil {
		return resp, nil
	}

	return resp, nil
}
