package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceComputePfwListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computePfwList, err := utilityComputePfwListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenPfwList(computePfwList))
	return nil
}

func dataSourceComputePfwListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pfw_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"local_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"local_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"protocol": {
						Type:     schema.TypeString,
						Computed: true,
					},
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
				},
			},
		},
	}
}

func DataSourceComputePfwList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputePfwListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputePfwListSchemaMake(),
	}
}
