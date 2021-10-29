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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func extractJob(d *schema.ResourceData) *client.Job {
	raw := d.Get("")
	var job client.Job
	helper.TerraformToStruct(raw, &job)
	return &job
}

func resourceJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	job := extractJob(d)
	log.Printf("CONVERTED job: %+v\n", job)

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

	err = d.Set("name", job.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("type", job.Type)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("framework", job.Framework)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("version", job.Version)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("description", job.Description)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("parallelism", job.Parallelism) // TODO
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("completions", job.Completions)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("suspend", job.Suspend)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("backoff_limit", job.BackoffLimit)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("containers", helper.StructToTerraform(&job.Containers))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("policy", helper.StructToTerraform(job.Policy))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("") {
		return resourceJobRead(ctx, d, m)
	}

	job := extractJob(d)

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
