package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgAffinityGroupsListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	list, err := utilityRgAffinityGroupsListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	d.Set("affinity_groups", flattenRgListGroups(list))
	return nil
}

func dataSourceRgAffinityGroupsListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the RG",
		},
		"affinity_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ids": {
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

func DataSourceRgAffinityGroupsList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgAffinityGroupsListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgAffinityGroupsListSchemaMake(),
	}
}
