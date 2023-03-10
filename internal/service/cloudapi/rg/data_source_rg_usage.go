package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgUsageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	usage, err := utilityDataRgUsageCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	flattenRgUsageResource(d, *usage)
	return nil
}

func dataSourceRgUsageSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"reason": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_size_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"extips": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"exttraffic": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"seps": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sep_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"map": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}

	return res
}

func DataSourceRgUsage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgUsageRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgUsageSchemaMake(),
	}
}
