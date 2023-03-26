package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
)

func dataSourceComputeSnapshotUsageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeSnapshotUsage, err := utilityComputeSnapshotUasgeCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenSnapshotUsage(computeSnapshotUsage))
	return nil
}

func dataSourceComputeSnapshotUsagSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"label": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"stored": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"timestamp": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func DataSourceComputeSnapshotUsage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeSnapshotUsageRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeSnapshotUsagSchemaMake(),
	}
}
