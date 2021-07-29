package k8s

import (
	"context"
	"fmt"

	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/RasaHQ/rasaxctl/pkg/utils/cloud"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	kubeconfig    string
	clientset     *kubernetes.Clientset
	Namespace     string
	Helm          HelmSpec
	Log           logr.Logger
	BackendType   types.KubernetesBackendType
	CloudProvider *cloud.Provider
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

	backendType, err := k.detectBackend()
	if err != nil {
		return err
	}
	k.BackendType = backendType

	k.Log.Info("Initializing Kubernetes client")

	return nil
}

func (k *Kubernetes) GetRasaXURL() (string, error) {

	nginxValues := k.Helm.Values["nginx"]
	ingressValues := k.Helm.Values["ingress"]
	nginxIsEnabled := nginxValues.(map[string]interface{})["enabled"].(bool)
	ingressIsEnabled := ingressValues.(map[string]interface{})["enabled"].(bool)
	nginxServiceType := nginxValues.(map[string]interface{})["service"].(map[string]interface{})["type"].(string)
	rasaXScheme := k.Helm.Values["rasax"].(map[string]interface{})["scheme"].(string)
	serviceName := fmt.Sprintf("%s-nginx", k.Helm.ReleaseName)

	url := "UNKNOWN"

	if nginxServiceType == "LoadBalancer" && nginxIsEnabled && (k.BackendType != types.KubernetesBackendLocal || k.CloudProvider.Name != types.CloudProviderUnknown) {

		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		ipAddress := service.Status.LoadBalancer.Ingress[0].IP
		port := service.Spec.Ports[0].Port

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, ipAddress, port)

		return url, nil
	} else if nginxServiceType == "NodePort" && nginxIsEnabled && k.BackendType == types.KubernetesBackendLocal && k.CloudProvider.Name != types.CloudProviderUnknown {

		service, err := k.clientset.CoreV1().Services(k.Namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			return url, err
		}
		port := service.Spec.Ports[0].NodePort

		url = fmt.Sprintf("%s://%s:%d", rasaXScheme, k.CloudProvider.ExternalIP, port)
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
			k.Log.V(1).Info("Deployment has replica number set to 0",
				"statefulset", deployment.Name, "replicas", deployment.Status.Replicas, "readyReplicas", deployment.Status.ReadyReplicas)
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
			k.Log.V(1).Info("Statefulset has replica number set to 0",
				"statefulset", statefulset.Name, "replicas", statefulset.Status.Replicas, "readyReplicas", statefulset.Status.ReadyReplicas)
			return false, nil
		}
	}

	return true, nil
}

func (k *Kubernetes) CreateNamespace() error {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: k.Namespace,
			Labels: map[string]string{
				"rasaxctl": "true",
			},
		},
	}
	_, err := k.clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	switch t := err.(type) {
	case *errors.StatusError:
		if t.ErrStatus.Code == 409 {
			k.Log.V(1).Info("Namespace already exists", "namespace", k.Namespace)
			return nil
		}
		return err
	}

	k.Log.V(1).Info("Create namespace", "namespace", k.Namespace)
	return nil
}

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

func (k *Kubernetes) GetPods() (*v1.PodList, error) {
	pods, err := k.clientset.CoreV1().Pods(k.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}
