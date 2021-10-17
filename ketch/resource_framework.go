package ketch

import (
	"context"
	"github.com/brunoa19/ketch-terraform-provider/client"
	"github.com/brunoa19/ketch-terraform-provider/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
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
			"framework": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
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
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFrameworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	raw := d.Get("framework").([]interface{})[0].(map[string]interface{})

	framework := &client.Framework{}
	helper.TerraformToStruct(raw, framework)

	log.Printf("RAW framework: %+v\n", raw)
	log.Printf("CONVERTED framework: %+v\n", *framework)

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

	if err = d.Set("framework", helper.StructToTerraform(framework)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFrameworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("framework") {
		return resourceFrameworkRead(ctx, d, m)
	}

	raw := d.Get("framework").([]interface{})[0].(map[string]interface{})

	framework := &client.Framework{}
	helper.TerraformToStruct(raw, framework)

	log.Printf("RAW framework: %+v\n", raw)

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
