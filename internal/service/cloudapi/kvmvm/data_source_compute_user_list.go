package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceComputeUserListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeUserList, err := utilityComputeUserListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	flattenUserList(d, computeUserList)
	return nil
}

func dataSourceComputeUserListSchemaMake() map[string]*schema.Schema {
	res := computeACLSchemaMake()
	res["compute_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	return res
}

func DataSourceComputeUserList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeUserListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeUserListSchemaMake(),
	}
}
