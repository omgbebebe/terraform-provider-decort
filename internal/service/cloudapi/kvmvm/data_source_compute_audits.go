package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceComputeAuditsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeAudits, err := utilityComputeAuditsCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenComputeAudits(computeAudits))
	return nil
}

func dataSourceComputeAuditsSchemaMake() map[string]*schema.Schema {
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
}

func DataSourceComputeAudits() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeAuditsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeAuditsSchemaMake(),
	}
}
