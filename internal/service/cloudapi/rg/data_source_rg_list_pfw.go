package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgListPfwRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	listPfw, err := utilityRgListPfwCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	d.Set("items", flattenRgListPfw(listPfw))
	return nil
}

func dataSourceRgListPfwSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the RG",
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"public_port_end": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"public_port_start": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vm_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vm_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vins_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vins_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	return res
}

func DataSourceRgListPfw() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgListPfwRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgListPfwSchemaMake(),
	}
}
