package ketch

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/brunoa19/ketch-terraform-provider/client"
)

func TestExtractApp(t *testing.T) {
	d := schema.TestResourceDataRaw(t, schemaApp, map[string]interface{}{
		"name":      "testjob",
		"image":     "gcr.io/test",
		"framework": "testfw",
		"cnames":    []interface{}{"cname1", "cname2"},
		"ports":     []interface{}{8080, 8081},
		"units":     4,
		"processes": []interface{}{
			map[string]interface{}{
				"cmd":  []interface{}{"./web"},
				"name": "web",
			},
			map[string]interface{}{
				"cmd":  []interface{}{"./worker", "-v"},
				"name": "worker",
			},
		},
		"routing_settings": []interface{}{
			map[string]interface{}{"weight": 100},
		},
		"version":       2,
		"backoff_limit": 5,
		"policy": []interface{}{
			map[string]interface{}{
				"restart_policy": "NEVER",
			}},
		"containers": []interface{}{
			map[string]interface{}{
				"name":    "testcontainer",
				"image":   "gcr.io/test",
				"command": []interface{}{"ls", "-a"},
			},
		},
	})
	expected := &client.App{
		Name:      "testjob",
		Image:     "gcr.io/test",
		Framework: "testfw",
		Cname:     []string{"cname1", "cname2"},
		Ports:     []int{8080, 8081},
		Units:     4,
		Processes: []*client.ProcessParameters{
			{
				Cmd:  []string{"./web"},
				Name: "web",
			},
			{
				Cmd:  []string{"./worker", "-v"},
				Name: "worker",
			},
		},
		RoutingSettings: &client.RoutingSettings{Weight: 100},
		Version:         2,
	}
	app := extractApp(d)
	require.Equal(t, expected, app)
}
