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
package k8s

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateVolume creates a volume that uses a local host path.
func (k *Kubernetes) CreateVolume(hostPath string) (string, error) {

	pv, err := k.createPV(hostPath)
	if err != nil {
		return "", err
	}

	pvc, err := k.createPVC(pv)
	if err != nil {
		return "", err
	}

	return pvc.Name, nil
}

// DeleteVolumes deletes a volume that uses a local host path.
func (k *Kubernetes) DeleteVolume() error {
	pvc := fmt.Sprintf("rasactl-pvc-%s", k.Namespace)
	if err := k.deletePVC(pvc); err != nil {
		return err
	}

	pv := fmt.Sprintf("rasactl-pv-%s", k.Namespace)
	err := k.deletePV(pv)

	return err
}

func (k *Kubernetes) createPV(hostPath string) (*apiv1.PersistentVolume, error) {

	pvSpec := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rasactl-pv-%s", k.Namespace),
			Namespace: k.Namespace,
			Labels: map[string]string{
				"rasactl": "true",
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			StorageClassName: "standard",
			AccessModes:      []apiv1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Capacity: apiv1.ResourceList{
				apiv1.ResourceStorage: resource.MustParse("2Gi"),
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: hostPath,
				},
			},
		},
	}

	pv, err := k.clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvSpec, metav1.CreateOptions{})
	if err != nil {
		return pv, err
	}

	k.Log.V(1).Info("Persistent Volume has been created",
		"name", pv.Name, "namespace", pv.Namespace, "hostPath", hostPath,
	)
	return pv, nil
}

func (k *Kubernetes) createPVC(pv *apiv1.PersistentVolume) (*apiv1.PersistentVolumeClaim, error) {

	pvcSpec := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rasactl-pvc-%s", k.Namespace),
			Namespace: k.Namespace,
			Labels: map[string]string{
				"rasactl": "true",
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse(pv.Spec.Capacity.Storage().String()),
				},
			},
			VolumeName: pv.Name,
		},
	}

	pvc, err := k.clientset.CoreV1().PersistentVolumeClaims(k.Namespace).Create(context.TODO(), pvcSpec, metav1.CreateOptions{})
	if err != nil {
		return pvc, err
	}
	k.Log.V(1).Info("Persistent Volume Claim has been created", "name", pvc.Name, "namespace", pvc.Namespace)
	return pvc, nil
}

func (k *Kubernetes) deletePV(name string) error {

	if err := k.clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	k.Log.V(1).Info("Persistent Volume has been deleted",
		"name", name, "namespace", k.Namespace)
	return nil
}

func (k *Kubernetes) deletePVC(name string) error {

	if err := k.clientset.CoreV1().PersistentVolumeClaims(k.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	k.Log.V(1).Info("Persistent Volume Claim has been deleted",
		"name", name, "namespace", k.Namespace)
	return nil
}
