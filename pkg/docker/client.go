/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package docker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/RasaHQ/rasaxctl/pkg/status"
	rtypes "github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
)

type Docker struct {
	Client       *client.Client
	ctx          context.Context
	Namespace    string
	Log          logr.Logger
	Spinner      *status.SpinnerMessage
	ProjectPath  string
	kubeadmToken string
	Kind         KindSpec
	Flags        *rtypes.RasaXCtlFlags
}

type KindSpec struct {
	ControlPlaneHost string
	Version          string
}

const (
	kindImagePrefix string = "kindest/node:"
)

func (d *Docker) New() error {
	d.ctx = context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	d.Client = cli

	d.Log.Info("Initializing Docker client")

	return nil
}

func (d *Docker) prepareKindJoinConfiguration() (string, error) {

	token, err := d.getKubeadmToken()
	if err != nil {
		return "", err
	}
	d.kubeadmToken = token

	file := fmt.Sprintf("/tmp/rasaxctl-kind-joinconfig-%s.yaml", d.Namespace)

	joinConfiguration := bootstrapv1.JoinConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "JoinConfiguration",
			APIVersion: "kubeadm.k8s.io/v1beta2",
		},
		Discovery: bootstrapv1.Discovery{
			BootstrapToken: &bootstrapv1.BootstrapTokenDiscovery{
				Token:                    token,
				APIServerEndpoint:        fmt.Sprintf("%s:6443", d.Kind.ControlPlaneHost),
				UnsafeSkipCAVerification: true,
			},
			Timeout:           &metav1.Duration{Duration: time.Minute * 5},
			TLSBootstrapToken: token,
		},
		NodeRegistration: bootstrapv1.NodeRegistrationOptions{
			CRISocket: "unix:///run/containerd/containerd.sock",
			KubeletExtraArgs: map[string]string{
				"fail-swap-on": "false",
				"node-labels":  fmt.Sprintf("rasaxctl-project=%s", d.Namespace),
			},
			Name: fmt.Sprintf("kind-%s", d.Namespace),
			Taints: []v1.Taint{{
				Key:    "rasaxctl",
				Value:  "true",
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

	d.Log.V(1).Info("Creating a kubeadm join configuration", "configuration", joinConfiguration)

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

	d.Log.Info("Copying join configuration to a kind container", "container", container.ID)
	return nil
}

func (d *Docker) getKubeadmToken() (string, error) {
	token := new(bytes.Buffer)
	execSpec, err := d.Client.ContainerExecCreate(d.ctx, d.Kind.ControlPlaneHost, types.ExecConfig{
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

	d.Log.Info("Getting join token", "token", r)

	return r, nil
}

func (d *Docker) deleteKubeadmToken() error {
	execSpec, err := d.Client.ContainerExecCreate(d.ctx, d.Kind.ControlPlaneHost, types.ExecConfig{
		WorkingDir:   "/",
		Cmd:          []string{"kubeadm", "token", "delete", d.kubeadmToken},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}
	if err := d.Client.ContainerExecStart(d.ctx, execSpec.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	d.Log.Info("Removing kubeadm join token", "token", d.kubeadmToken)

	return nil
}

func (d *Docker) joinKindNodeToKubernetesCluster(container container.ContainerCreateCreatedBody) error {
	execSpec, err := d.Client.ContainerExecCreate(d.ctx, container.ID, types.ExecConfig{
		WorkingDir:   "/",
		Cmd:          []string{"kubeadm", "join", "--config", "config.yaml", "--skip-phases=preflight", "-v", "6"},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}
	cmdReader, err := d.Client.ContainerExecAttach(d.ctx, execSpec.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	cmdLogs := bufio.NewReader(cmdReader.Reader)
	for {
		line, _, err := cmdLogs.ReadLine()
		if len(line) > 0 {
			d.Log.V(1).Info("Joining a kind node to the cluster", "details", string(line))
		}
		if err != nil {
			break
		}
	}
	cmdReader.Close()

	d.Spinner.Message("Waiting for kind node to join to the cluster")
	for {
		status, err := d.Client.ContainerExecInspect(d.ctx, execSpec.ID)
		if err != nil {
			panic(err)
		}
		d.Log.Info("Waiting for kind node to join to the cluster", "container", container.ID, "running", status.Running, "exitCode", status.ExitCode)
		if status.ExitCode != 0 {
			return errors.Errorf("Can't join kind node to the cluster")
		}
		if !status.Running {
			break
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

func (d *Docker) getKindNetwork() (string, error) {

	inspect, err := d.Client.ContainerInspect(d.ctx, d.Kind.ControlPlaneHost)
	if err != nil {
		return "", err
	}

	return inspect.HostConfig.NetworkMode.NetworkName(), nil
}

func (d *Docker) CreateKindNode(hostname string) (container.ContainerCreateCreatedBody, error) {
	kindImage := fmt.Sprintf("%s%s", kindImagePrefix, d.Kind.Version)

	d.Log.Info("Pulling image", "image", kindImage)
	imagePull, err := d.Client.ImagePull(d.ctx, kindImage, types.ImagePullOptions{})
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	imagePullLogs := bufio.NewReader(imagePull)
	for {
		line, _, err := imagePullLogs.ReadLine()
		if len(line) > 0 {
			d.Log.V(1).Info("Pulling image", "details", string(line))
		}
		if err != nil {
			break
		}
	}
	imagePull.Close()

	d.Log.Info("Creating a kind node", "node", hostname, "image", kindImage)

	hostConfig := &container.HostConfig{
		Privileged:  true,
		SecurityOpt: []string{"apparmor=unconfined", "seccomp=unconfined"},
		Tmpfs:       map[string]string{"/run": "", "/tmp": ""},
	}

	hostConfig.ExtraHosts = []string{
		"rasa.localhost:host-gateway",
		"host.docker.internal:host-gateway",
	}

	hostConfig.Mounts = []mount.Mount{
		{
			Source:   "/lib/modules",
			Target:   "/lib/modules",
			Type:     mount.TypeBind,
			ReadOnly: true,
		},
	}

	if d.ProjectPath != "" {
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Source: d.ProjectPath,
			Target: d.ProjectPath,
			Type:   mount.TypeBind,
		})

	}

	kindNetwork, err := d.getKindNetwork()
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	resp, err := d.Client.ContainerCreate(d.ctx,
		&container.Config{
			Image:    kindImage,
			Tty:      false,
			Hostname: hostname,
			Volumes:  map[string]struct{}{"/var": {}},
			Env:      []string{"container=docker"},
		},
		hostConfig, &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{kindNetwork: {}},
		}, hostname)
	if err != nil {
		return resp, err
	}

	if err := d.Client.ContainerStart(d.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return resp, err
	}

	if err := d.copyJoinConfigurationToContainer(resp); err != nil {
		return resp, err
	}

	if err := d.joinKindNodeToKubernetesCluster(resp); err != nil {
		return resp, err
	}

	if err := d.deleteKubeadmToken(); err != nil {
		return resp, err
	}
	return resp, nil
}

func (d *Docker) StopKindNode(hostname string) error {
	timeout := time.Minute * 1
	if err := d.Client.ContainerStop(d.ctx, hostname, &timeout); err != nil {
		return err
	}
	return nil
}

func (d *Docker) StartKindNode(hostname string) error {
	if err := d.Client.ContainerStart(d.ctx, hostname, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *Docker) DeleteKindNode(hostname string) error {
	if err := d.Client.ContainerRemove(d.ctx, hostname, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return err
	}
	return nil
}
