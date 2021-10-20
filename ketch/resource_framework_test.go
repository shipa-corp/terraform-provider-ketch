package ketch

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/brunoa19/ketch-terraform-provider/client"
)

func TestExtractFramework(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceFramework().Schema, map[string]interface{}{
		"name":            "testfw",
		"namespace":       "testns",
		"app_quota_limit": 2,
		"ingress_controller": []interface{}{
			map[string]interface{}{
				"class_name":       "traefik",
				"service_endpoint": "10.10.10.10",
				"type":             "traefik",
				"cluster_issuer":   "test_issuer",
			}},
	})
	expected := &client.Framework{
		Name:          "testfw",
		Namespace:     "testns",
		AppQuotaLimit: 2,
		IngressController: &client.IngressControllerSpec{
			ClassName:       "traefik",
			ServiceEndpoint: "10.10.10.10",
			IngressType:     "traefik",
			ClusterIssuer:   "test_issuer",
		},
	}

	framework := extractFramework(d)
	require.Equal(t, expected, framework)
}
