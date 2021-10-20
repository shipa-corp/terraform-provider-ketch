package ketch

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/brunoa19/ketch-terraform-provider/client"
)

func TestExtractJob(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceJob().Schema, map[string]interface{}{
		"name":          "testjob",
		"type":          "test",
		"framework":     "testfw",
		"version":       "v1",
		"description":   "a test",
		"parallelism":   3,
		"completions":   2,
		"suspend":       false,
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
	expected := &client.Job{
		Name:         "testjob",
		Type:         "test",
		Framework:    "testfw",
		Version:      "v1",
		Description:  "a test",
		Parallelism:  3,
		Completions:  2,
		Suspend:      false,
		BackoffLimit: 5,
		Containers: []*client.Container{
			{
				Name:    "testcontainer",
				Image:   "gcr.io/test",
				Command: []string{"ls", "-a"},
			},
		},
		Policy: &client.Policy{
			RestartPolicy: "NEVER",
		},
	}
	job := extractJob(d)
	require.Equal(t, expected, job)
}
