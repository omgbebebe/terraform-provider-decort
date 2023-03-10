package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgAffinityGroupComputesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rgComputes, err := utilityRgAffinityGroupComputesCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	d.Set("items", flattenRgAffinityGroupComputes(rgComputes))
	return nil
}

func dataSourceRgAffinityGroupComputesSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the RG",
		},
		"affinity_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Affinity group label",
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"compute_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"other_node": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"other_node_indirect": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"other_node_indirect_soft": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"other_node_soft": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"same_node": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"same_node_soft": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
				},
			},
		},
	}

	return res
}

func DataSourceRgAffinityGroupComputes() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgAffinityGroupComputesRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgAffinityGroupComputesSchemaMake(),
	}
}
