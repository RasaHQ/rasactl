package k8s

import (
	"context"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) ScaleDown() error {
	deployments, err := k.clientset.AppsV1().Deployments(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		var err error
		var scale *autoscalingv1.Scale

		k.Log.V(1).Info("Scaling down", "deployment", deployment.Name)
		scale, err = k.clientset.AppsV1().Deployments(k.Namespace).GetScale(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		scale.Spec.Replicas = 0
		_, err = k.clientset.AppsV1().Deployments(k.Namespace).UpdateScale(context.TODO(), deployment.Name, scale, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	statefulsets, err := k.clientset.AppsV1().StatefulSets(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, statefulsets := range statefulsets.Items {
		var err error
		var scale *autoscalingv1.Scale

		k.Log.V(1).Info("Scaling down", "statefulsets", statefulsets.Name)
		scale, err = k.clientset.AppsV1().StatefulSets(k.Namespace).GetScale(context.TODO(), statefulsets.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		scale.Spec.Replicas = 0
		_, err = k.clientset.AppsV1().StatefulSets(k.Namespace).UpdateScale(context.TODO(), statefulsets.Name, scale, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *Kubernetes) ScaleUp() error {
	deployments, err := k.clientset.AppsV1().Deployments(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		var err error
		var scale *autoscalingv1.Scale

		scale, err = k.clientset.AppsV1().Deployments(k.Namespace).GetScale(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if scale.Spec.Replicas != 0 {
			continue
		}
		k.Log.V(1).Info("Scaling up", "deployment", deployment.Name)
		scale.Spec.Replicas = 1
		_, err = k.clientset.AppsV1().Deployments(k.Namespace).UpdateScale(context.TODO(), deployment.Name, scale, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	statefulsets, err := k.clientset.AppsV1().StatefulSets(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, statefulsets := range statefulsets.Items {
		var err error
		var scale *autoscalingv1.Scale

		scale, err = k.clientset.AppsV1().StatefulSets(k.Namespace).GetScale(context.TODO(), statefulsets.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if scale.Spec.Replicas != 0 {
			continue
		}
		k.Log.V(1).Info("Scaling up", "statefulsets", statefulsets.Name)
		scale.Spec.Replicas = 1
		_, err = k.clientset.AppsV1().StatefulSets(k.Namespace).UpdateScale(context.TODO(), statefulsets.Name, scale, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
