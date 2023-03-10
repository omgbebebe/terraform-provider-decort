package kvmvm

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceComputeGetConsoleUrlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeConsoleUrl, err := utilityComputeGetConsoleUrlCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	result := strings.ReplaceAll(string(computeConsoleUrl), "\"", "")
	result = strings.ReplaceAll(string(result), "\\", "")
	d.Set("console_url", result)
	return nil
}

func dataSourceComputeGetConsoleUrlSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"console_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceComputeGetConsoleUrl() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeGetConsoleUrlRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeGetConsoleUrlSchemaMake(),
	}
}
