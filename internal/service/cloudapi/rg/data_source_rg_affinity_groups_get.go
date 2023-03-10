package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgAffinityGroupsGetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computes, err := utilityRgAffinityGroupsGetCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	d.Set("ids", computes)
	return nil
}

func dataSourceRgAffinityGroupsGetSchemaMake() map[string]*schema.Schema {
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

		"ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}

	return res
}

func DataSourceRgAffinityGroupsGet() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgAffinityGroupsGetRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgAffinityGroupsGetSchemaMake(),
	}
}
