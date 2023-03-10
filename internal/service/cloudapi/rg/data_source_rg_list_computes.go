package rg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgListComputesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	listComputes, err := utilityRgListComputesCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenRgListComputes(listComputes))
	return nil
}

func rulesSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"policy": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"topology": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return res

}

func dataSourceRgListComputesSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the RG",
		},
		"reason": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "reason for action",
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"affinity_label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"affinity_rules": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: rulesSchemaMake(),
						},
					},
					"affinity_weight": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"antiaffinity_rules": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: rulesSchemaMake(),
						},
					},
					"cpus": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"created_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"deleted_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ram": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"registered": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"rg_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tech_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"total_disks_size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"updated_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"user_managed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vins_connected": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}

	return res
}

func DataSourceRgListComputes() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgListComputesRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgListComputesSchemaMake(),
	}
}
