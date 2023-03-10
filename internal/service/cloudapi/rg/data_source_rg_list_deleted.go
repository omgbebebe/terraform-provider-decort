package rg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgListDeletedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rgList, err := utilityRgListDeletedCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenRgList(rgList))

	return nil
}

func dataSourceRgListDeletedSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"page": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Page number",
		},
		"size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Page size",
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_acl": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: aclSchemaMake(),
						},
					},
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"acl": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: aclSchemaMake(),
						},
					},
					"created_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"def_net_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"def_net_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"deleted_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"dirty": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"gid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"rg_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"lock_status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"milestones": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"register_computes": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"resource_limits": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: resourceLimitsSchemaMake(),
						},
					},
					"secret": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vins": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"vms": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"resource_types": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"uniq_pools": {
						Type:     schema.TypeList,
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

func DataSourceRgListDeleted() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgListDeletedRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgListDeletedSchemaMake(),
	}
}
