package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgAuditsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rgAudits, err := utilityRgAuditsCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	d.Set("items", flattenRgAudits(rgAudits))

	return nil
}

func dataSourceRgAuditsSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:     schema.TypeInt,
			Required: true,
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"call": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"responsetime": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"statuscode": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"timestamp": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"user": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	return res
}

func DataSourceRgAudits() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgAuditsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgAuditsSchemaMake(),
	}
}
