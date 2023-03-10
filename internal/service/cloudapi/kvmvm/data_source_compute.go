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

package kvmvm

import (
	"context"
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func findInExtraDisks(DiskId uint, ExtraDisks []interface{}) bool {
	for _, ExtraDisk := range ExtraDisks {
		if DiskId == uint(ExtraDisk.(int)) {
			return true
		}
	}
	return false
}

func dataSourceComputeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	compute, err := utilityDataComputeCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(int(compute.ID)))

	flattenDataCompute(d, compute)
	return nil
}

func computeListRulesSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"policy": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"topology": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	return res
}

func computeListACLSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"explicit": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"right": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user_group_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func computeACLSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListACLSchemaMake(),
			},
		},
		"compute_acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListACLSchemaMake(),
			},
		},
		"rg_acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListACLSchemaMake(),
			},
		},
	}
}

func computeIOTuneSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"read_bytes_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"read_bytes_sec_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"read_iops_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"read_iops_sec_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"size_iops_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"total_bytes_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"total_bytes_sec_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"total_iops_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"total_iops_sec_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"write_bytes_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"write_bytes_sec_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"write_iops_sec": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"write_iops_sec_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func computeSnapshotsSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"res_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"snap_set_guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"snap_set_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"timestamp": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}
func computeListDisksSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"_ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"acl": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"boot_partition": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"destruction_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_path": {
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
		"disk_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"images": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"iotune": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeIOTuneSchemaMake(),
			},
		},
		"iqn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"login": {
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
		"order": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"params": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"parent_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"passwd": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pci_slot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pool": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"present_to": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"purge_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"reality_device_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"res_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"role": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sep_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"shareable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"size_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"size_used": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"snapshots": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeSnapshotsSchemaMake(),
			},
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vmid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func computeQOSSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"e_rate": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"in_brust": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"in_rate": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func computeInterfacesSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"conn_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"conn_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"def_gw": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"flip_group_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"listen_ssh": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"net_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"netmask": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"net_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pci_slot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"qos": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeQOSSchemaMake(),
			},
		},
		"target": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vnfs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
}
func computeOsUsersSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"login": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
func computeSnapSetsSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"guid": {
			Type:     schema.TypeString,
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
	}
}
func dataSourceComputeSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Required: true,
		},

		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeACLSchemaMake(),
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
		"custom_fields": {
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
		"devices": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeListDisksSchemaMake(),
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
		"natable_vins_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"natable_vins_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_network_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"os_users": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeOsUsersSchemaMake(),
			},
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
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tech_status": {
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
		"user_managed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"userdata": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vgpus": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"virtual_image_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"virtual_image_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceCompute() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeSchemaMake(),
	}
}
