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

package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/constants"
)

func dataSourceRgUsageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	usage, err := utilityDataRgUsageCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	flattenRgUsageResource(d, *usage)
	return nil
}

func dataSourceRgUsageSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"reason": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_size": {
			Type:     schema.TypeInt,
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
	}

	return res
}

func DataSourceRgUsage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgUsageRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgUsageSchemaMake(),
	}
}
