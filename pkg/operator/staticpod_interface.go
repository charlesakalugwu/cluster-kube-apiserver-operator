package operator

import (
	"k8s.io/client-go/tools/cache"

	operatorv1 "github.com/openshift/api/operator/v1"
	operatorconfigclientv1alpha1 "github.com/openshift/cluster-kube-apiserver-operator/pkg/generated/clientset/versioned/typed/kubeapiserver/v1alpha1"
	operatorclientinformers "github.com/openshift/cluster-kube-apiserver-operator/pkg/generated/informers/externalversions"
)

type staticPodOperatorClient struct {
	informers operatorclientinformers.SharedInformerFactory
	client    operatorconfigclientv1alpha1.KubeapiserverV1alpha1Interface
}

func (c *staticPodOperatorClient) Informer() cache.SharedIndexInformer {
	return c.informers.Kubeapiserver().V1alpha1().KubeAPIServerOperatorConfigs().Informer()
}

func (c *staticPodOperatorClient) Get() (*operatorv1.OperatorSpec, *operatorv1.StaticPodOperatorStatus, string, error) {
	instance, err := c.informers.Kubeapiserver().V1alpha1().KubeAPIServerOperatorConfigs().Lister().Get("instance")
	if err != nil {
		return nil, nil, "", err
	}

	return &instance.Spec.OperatorSpec, &instance.Status.StaticPodOperatorStatus, instance.ResourceVersion, nil
}

func (c *staticPodOperatorClient) UpdateStatus(resourceVersion string, status *operatorv1.StaticPodOperatorStatus) (*operatorv1.StaticPodOperatorStatus, error) {
	original, err := c.informers.Kubeapiserver().V1alpha1().KubeAPIServerOperatorConfigs().Lister().Get("instance")
	if err != nil {
		return nil, err
	}
	copy := original.DeepCopy()
	copy.ResourceVersion = resourceVersion
	copy.Status.StaticPodOperatorStatus = *status

	ret, err := c.client.KubeAPIServerOperatorConfigs().UpdateStatus(copy)
	if err != nil {
		return nil, err
	}

	return &ret.Status.StaticPodOperatorStatus, nil
}

// TODO collapse this onto get
func (c *staticPodOperatorClient) CurrentStatus() (operatorv1.OperatorStatus, error) {
	instance, err := c.informers.Kubeapiserver().V1alpha1().KubeAPIServerOperatorConfigs().Lister().Get("instance")
	if err != nil {
		return operatorv1.OperatorStatus{}, err
	}

	return instance.Status.OperatorStatus, nil
}
