/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Terraform DECORT provider - manage resources provided by DECORT (Digital Energy Cloud
Orchestration Technology) with Terraform by Hashicorp.

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package rg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceRgListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rgList, err := utilityRgListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenRgList(rgList))

	return nil
}

func dataSourceRgListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"includedeleted": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "included deleted resource groups",
		},
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

func DataSourceRgList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgListSchemaMake(),
	}
}
