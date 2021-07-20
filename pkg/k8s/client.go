package k8s

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	kubeconfig string
	clientset  *kubernetes.Clientset
	Namespace  string
	Helm       HelmSpec
	Log        logr.Logger
}

type HelmSpec struct {
	Values      map[string]interface{}
	ReleaseName string
}

func (k *Kubernetes) New() error {
	k.kubeconfig = viper.GetString("kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", k.kubeconfig)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k.clientset = clientset

	return nil
}

func (k *Kubernetes) GetRasaXURL() (string, error) {

	nginxValues := k.Helm.Values["nginx"]
	ingressValues := k.Helm.Values["ingress"]
	nginxIsEnabled := nginxValues.(map[string]interface{})["enabled"].(bool)
	ingressIsEnabled := ingressValues.(map[string]interface{})["enabled"].(bool)
	nginxServiceType := nginxValues.(map[string]interface{})["service"].(map[string]interface{})["type"].(string)
	rasaXScheme := k.Helm.Values["rasax"].(map[string]interface{})["scheme"].(string)
	url := "UNKNOWN"

	if nginxServiceType == "LoadBalancer" && nginxIsEnabled {
		serviceName := fmt.Sprintf("%s-nginx", k.Helm.ReleaseName)
		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		ipAddress := service.Status.LoadBalancer.Ingress[0].IP
		port := service.Spec.Ports[0].Port

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, ipAddress, port)

		return url, nil
	} else if ingressIsEnabled {
		ingress, err := k.clientset.NetworkingV1().Ingresses(k.Namespace).Get(context.TODO(), k.Helm.ReleaseName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		host := ingress.Spec.Rules[0].Host

		url = fmt.Sprintf("%s://%s", rasaXScheme, host)
	}

	return url, nil
}

func (k *Kubernetes) GetRasaXToken() (string, error) {
	var token string

	secretName := fmt.Sprintf("%s-rasa", k.Helm.ReleaseName)
	keyName := "rasaXToken"

	secret, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return token, err
	}

	return string(secret.Data[keyName]), nil
}

func (k *Kubernetes) IsRasaXRunning() (bool, error) {

	deployments, err := k.clientset.AppsV1().Deployments(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	if len(deployments.Items) == 0 {
		return false, nil
	}

	for _, deployment := range deployments.Items {
		if deployment.Status.Replicas == 0 || deployment.Status.ReadyReplicas == 0 {
			return false, nil
		}
	}

	statefulsets, err := k.clientset.AppsV1().StatefulSets(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	if len(statefulsets.Items) == 0 {
		return false, nil
	}

	for _, statefulset := range statefulsets.Items {
		if statefulset.Status.Replicas == 0 || statefulset.Status.ReadyReplicas == 0 {
			return false, nil
		}
	}

	return true, nil
}

/*func (k *Kubernetes) ScaleUp() error {
	deployments, err := k.clientset.AppsV1().Deployments(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		deployment.GetScale
	}

	return nil
}*/
