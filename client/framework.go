package client

import (
	"context"
	"github.com/brunoa19/ketch-terraform-provider/client/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

type Framework struct {
	Name              string                 `json:"name"`
	Namespace         string                 `json:"namespace"`
	AppQuotaLimit     int64                  `json:"app_quota_limit,omitempty"`
	IngressController *IngressControllerSpec `json:"ingress_controller"`
}

// IngressControllerSpec contains configuration for an ingress controller.
type IngressControllerSpec struct {
	ClassName       string `json:"class_name"`
	ServiceEndpoint string `json:"service_endpoint"`
	IngressType     string `json:"type"`
	ClusterIssuer   string `json:"cluster_issuer"`
}

func NewFramework(input *v1beta1.Framework) *Framework {
	var appQuotaLimit int64
	if input.Spec.AppQuotaLimit != nil {
		appQuotaLimit = int64(*input.Spec.AppQuotaLimit)
	}

	return &Framework{
		Name:          input.Spec.Name,
		Namespace:     input.Spec.NamespaceName,
		AppQuotaLimit: appQuotaLimit,
		IngressController: &IngressControllerSpec{
			ClassName:       input.Spec.IngressController.ClassName,
			ServiceEndpoint: input.Spec.IngressController.ServiceEndpoint,
			ClusterIssuer:   input.Spec.IngressController.ClusterIssuer,
			IngressType:     input.Spec.IngressController.IngressType.String(),
		},
	}
}

func (f *Framework) convertToKetchFramework() *v1beta1.Framework {
	namespace := f.Namespace
	if namespace == "" {
		// default value
		namespace = "ketch-" + f.Name
	}

	// default value
	appQuotaLimit := -1
	if f.AppQuotaLimit != 0 {
		appQuotaLimit = int(f.AppQuotaLimit)
	}

	return &v1beta1.Framework{
		ObjectMeta: metav1.ObjectMeta{
			Name: f.Name,
		},
		Spec: v1beta1.FrameworkSpec{
			Name:          f.Name,
			NamespaceName: namespace,
			AppQuotaLimit: &appQuotaLimit,
			IngressController: v1beta1.IngressControllerSpec{
				ClassName:       f.IngressController.ClassName,
				ServiceEndpoint: f.IngressController.ServiceEndpoint,
				ClusterIssuer:   f.IngressController.ClusterIssuer,
				IngressType:     v1beta1.IngressControllerType(f.IngressController.IngressType),
			},
		},
	}
}

func (c *Client) GetFramework(ctx context.Context, name string) (*Framework, error) {
	framework, err := c.getFramework(ctx, name)
	if err != nil {
		return nil, err
	}

	return NewFramework(framework), nil
}

func (c *Client) DeleteFramework(ctx context.Context, name string) error {
	framework, err := c.getFramework(ctx, name)
	if err != nil {
		return err
	}

	err = c.kube.Delete(ctx, framework)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getFramework(ctx context.Context, name string) (*v1beta1.Framework, error) {
	framework := &v1beta1.Framework{}
	err := c.kube.Get(ctx, types.NamespacedName{Name: name}, framework)
	if err != nil {
		log.Println("ERR:", err.Error())
		return nil, err
	}

	return framework, nil
}

func (c *Client) CreateFramework(ctx context.Context, input *Framework) error {
	framework := input.convertToKetchFramework()
	return c.kube.Create(ctx, framework)
}

func (c *Client) UpdateFramework(ctx context.Context, input *Framework) error {
	framework, err := c.getFramework(ctx, input.Name)
	if err != nil {
		return err
	}

	updates := input.convertToKetchFramework()
	framework.Spec = updates.Spec
	return c.kube.Update(ctx, framework)
}
