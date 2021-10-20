package ketch

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/brunoa19/ketch-terraform-provider/client"
	"github.com/brunoa19/ketch-terraform-provider/helper"
)

var (
	schemaIngressController = &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"class_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service_endpoint": {
					Type:     schema.TypeString,
					Required: true,
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"cluster_issuer": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
)

func resourceFramework() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFrameworkCreate,
		ReadContext:   resourceFrameworkRead,
		UpdateContext: resourceFrameworkUpdate,
		DeleteContext: resourceFrameworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_quota_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ingress_controller": schemaIngressController,
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func extractFramework(d *schema.ResourceData) *client.Framework {
	raw := d.Get("")
	var framework client.Framework
	helper.TerraformToStruct(raw, &framework)
	return &framework
}

func resourceFrameworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	framework := extractFramework(d)
	log.Printf("CONVERTED create framework: %+v %+v\n", framework, framework.IngressController)

	c := m.(*client.Client)
	err := c.CreateFramework(ctx, framework)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(framework.Name)

	resourceFrameworkRead(ctx, d, m)

	return diags
}

func resourceFrameworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	c := m.(*client.Client)
	framework, err := c.GetFramework(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", framework.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("namespace", framework.Namespace)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("app_quota_limit", framework.AppQuotaLimit)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("ingress_controller", helper.StructToTerraform(framework.IngressController))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFrameworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("") {
		return resourceFrameworkRead(ctx, d, m)
	}

	framework := extractFramework(d)

	log.Printf("CONVERTED framework: %+v\n", framework)

	c := m.(*client.Client)
	err := c.UpdateFramework(ctx, framework)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFrameworkRead(ctx, d, m)
}

func resourceFrameworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()
	c := m.(*client.Client)
	err := c.DeleteFramework(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
