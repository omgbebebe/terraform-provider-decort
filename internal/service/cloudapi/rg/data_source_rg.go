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
	"strconv"

	"github.com/rudecs/terraform-provider-decort/internal/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func sepsSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"sep_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"data_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"disk_size": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"disk_size_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	return res
}

func resourcesSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"current": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_size": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"disk_size_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"extips": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"exttraffic": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"gpu": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ram": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"seps": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"sep_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"map": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
		"reserved": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_size": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"disk_size_max": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"extips": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"exttraffic": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"gpu": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ram": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"seps": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"sep_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"map": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return res
}

func aclSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
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

	return res
}

func resourceLimitsSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"cu_c": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"cu_d": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"cu_i": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"cu_m": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"cu_np": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"gpu_units": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
	}

	return res
}

func dataSourceRgSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"reason": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"resources": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: resourcesSchemaMake(),
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
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"res_types": {
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
	}
	return res
}

func dataSourceResgroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rg, err := utilityDataResgroupCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("") // ensure ID is empty in this case
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	flattenRg(d, *rg)
	return nil
}

func DataSourceResgroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceResgroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgSchemaMake(),
	}
}
