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

	schemaApp = map[string]*schema.Schema{
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
	}
)

func resourceApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppCreate,
		ReadContext:   resourceAppRead,
		UpdateContext: resourceAppUpdate,
		DeleteContext: resourceAppDelete,
		Schema:        schemaApp,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func extractApp(d *schema.ResourceData) *client.App {
	raw := d.Get("")
	var app client.App
	helper.TerraformToStruct(raw, &app)
	return &app
}

func resourceAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	app := extractApp(d)
	log.Printf("CONVERTED app: %+v\n", app)

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

	err = d.Set("name", app.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("image", app.Image)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("framework", app.Framework)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("cnames", app.Cname)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("ports", app.Ports)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("units", app.Units)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("processes", helper.StructToTerraform(&app.Processes))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("routing_settings", helper.StructToTerraform(app.RoutingSettings))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("version", app.Version)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("") {
		return resourceAppRead(ctx, d, m)
	}

	app := extractApp(d)

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
