package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceComputeGetLogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeGetLog, err := utilityComputeGetLogCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	d.Set("log", computeGetLog)
	return nil
}

func dataSourceComputeGetLogSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"path": {
			Type:     schema.TypeString,
			Required: true,
		},
		"log": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceComputeGetLog() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeGetLogRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeGetLogSchemaMake(),
	}
}
