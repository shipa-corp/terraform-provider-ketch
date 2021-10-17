package ketch

import (
	"context"
	"log"

	"github.com/brunoa19/ketch-terraform-provider/client"
	"github.com/brunoa19/ketch-terraform-provider/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	processesSchema = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"cmd": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}

	routingSettingsSchema = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"weight": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}

	schemaApp = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// Required
				"name": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"image": {
					Type:     schema.TypeString,
					Required: true,
				},
				"framework": {
					Type:     schema.TypeString,
					Required: true,
				},

				// Optional
				"cnames": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"ports": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
				"units": {
					Type:     schema.TypeInt,
					Optional: true,
				},

				"processes": processesSchema,

				"routing_settings": routingSettingsSchema,

				"version": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
)

func resourceApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppCreate,
		ReadContext:   resourceAppRead,
		UpdateContext: resourceAppUpdate,
		DeleteContext: resourceAppDelete,
		Schema: map[string]*schema.Schema{
			"app": schemaApp,
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	raw := d.Get("app").([]interface{})[0].(map[string]interface{})
	app := &client.App{}

	helper.TerraformToStruct(raw, app)

	c := m.(*client.Client)
	err := c.CreateApp(ctx, app)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.Name)

	resourceAppRead(ctx, d, m)

	return diags
}

func resourceAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	c := m.(*client.Client)
	app, err := c.GetApp(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("app", helper.StructToTerraform(app)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("app") {
		return resourceAppRead(ctx, d, m)
	}

	raw := d.Get("app").([]interface{})[0].(map[string]interface{})
	app := &client.App{}
	helper.TerraformToStruct(raw, app)

	log.Printf(" ### RAW app data: %+v\n", raw)
	log.Printf(" ### CONVERTED app data: %+v\n", *app)

	c := m.(*client.Client)
	err := c.UpdateApp(ctx, app)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAppRead(ctx, d, m)
}

func resourceAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()
	c := m.(*client.Client)
	err := c.DeleteApp(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
