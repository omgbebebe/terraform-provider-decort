package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceComputeGetAuditsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeAudits, err := utilityComputeGetAuditsCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenComputeGetAudits(computeAudits))
	return nil
}

func dataSourceComputeGetAuditsSchemaMake() map[string]*schema.Schema {
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
					"epoch": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"message": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func DataSourceComputeGetAudits() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeGetAuditsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeGetAuditsSchemaMake(),
	}
}
