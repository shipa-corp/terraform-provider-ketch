package client

import (
	"path/filepath"

	"github.com/brunoa19/ketch-terraform-provider/client/v1beta1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	kube client.Client
}

func NewClient() (*Client, error) {

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// create the config object from kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	schema, err := v1beta1.SchemeBuilder.Build()
	if err != nil {
		return nil, err
	}

	kube, err := client.New(config, client.Options{
		Scheme: schema,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		kube: kube,
	}, nil
}
