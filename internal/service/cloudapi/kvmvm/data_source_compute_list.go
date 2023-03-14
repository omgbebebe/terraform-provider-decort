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

Source code: https://repos.digitalenergy.online/BASIS/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://repos.digitalenergy.online/BASIS/terraform-provider-decort/wiki
*/

package kvmvm

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/constants"
)

func dataSourceComputeListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeList, err := utilityDataComputeListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenComputeList(computeList))
	return nil
}

func computeDisksSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pci_slot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}
func itemComputeSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListACLSchemaMake(),
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
		"affinity_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"affinity_rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListRulesSchemaMake(),
			},
		},
		"affinity_weight": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"anti_affinity_rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListRulesSchemaMake(),
			},
		},
		"arch": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"boot_order": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"bootdisk_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"clone_reference": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"clones": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"computeci_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cpus": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"custom_fields": { //NEED
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
		"devices": { //NEED
			Type:     schema.TypeString,
			Computed: true,
		},
		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeDisksSchemaMake(),
			},
		},
		"driver": {
			Type:     schema.TypeString,
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
		"compute_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"interfaces": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeInterfacesSchemaMake(),
			},
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manager_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"manager_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"migrationjob": {
			Type:     schema.TypeInt,
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
		"pinned": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"registered": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"res_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rg_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"snap_sets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeSnapSetsSchemaMake(),
			},
		},
		"stateless_sep_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"stateless_sep_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"val": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"total_disk_size": {
			Type:     schema.TypeInt,
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
		"user_managed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"vgpus": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"vins_connected": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"virtual_image_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceCompputeListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"includedeleted": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"page": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"size": {
			Type:     schema.TypeInt,
			Optional: true,
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: itemComputeSchemaMake(),
			},
		},
	}
	return res
}

func DataSourceComputeList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceCompputeListSchemaMake(),
	}
}
