package ketch

import (
	"context"
	"github.com/brunoa19/ketch-terraform-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ketch_app":       resourceApp(),
			"ketch_job":       resourceJob(),
			"ketch_framework": resourceFramework(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := client.NewClient()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Ketch client",
			Detail:   err.Error(),
		})

		return nil, diags
	}

	return c, diags
}
