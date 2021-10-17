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
	schemaJobContainers = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"image": {
					Type:     schema.TypeString,
					Required: true,
				},
				"command": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}

	schemaJobPolicy = &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"restart_policy": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
)

func resourceJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobCreate,
		ReadContext:   resourceJobRead,
		UpdateContext: resourceJobUpdate,
		DeleteContext: resourceJobDelete,
		Schema: map[string]*schema.Schema{
			"job": {
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
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"framework": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parallelism": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"completions": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"suspend": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"backoff_limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"containers": schemaJobContainers,
						"policy":     schemaJobPolicy,
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	raw := d.Get("job").([]interface{})[0].(map[string]interface{})

	job := &client.Job{}
	helper.TerraformToStruct(raw, job)

	log.Printf("RAW job: %+v\n", raw)
	log.Printf("CONVERTED job: %+v\n", *job)

	c := m.(*client.Client)
	err := c.CreateJob(ctx, job)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(job.Name)

	resourceJobRead(ctx, d, m)

	return diags
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()

	c := m.(*client.Client)
	job, err := c.GetJob(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("job", helper.StructToTerraform(job)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("job") {
		return resourceJobRead(ctx, d, m)
	}

	raw := d.Get("job").([]interface{})[0].(map[string]interface{})

	job := &client.Job{}
	helper.TerraformToStruct(raw, job)

	log.Printf("RAW job: %+v\n", raw)

	c := m.(*client.Client)
	err := c.UpdateJob(ctx, job)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceJobRead(ctx, d, m)
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Id()
	c := m.(*client.Client)
	err := c.DeleteJob(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
