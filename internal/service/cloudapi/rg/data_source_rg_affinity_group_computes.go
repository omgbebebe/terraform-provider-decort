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

Source code: https://repository.basistech.ru/BASIS/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://repository.basistech.ru/BASIS/terraform-provider-decort/wiki
*/

package rg

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
)

func dataSourceRgAffinityGroupComputesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rgComputes, err := utilityRgAffinityGroupComputesCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("rg_id").(int)))
	d.Set("items", flattenRgAffinityGroupComputes(rgComputes))
	return nil
}

func dataSourceRgAffinityGroupComputesSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the RG",
		},
		"affinity_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Affinity group label",
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"compute_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"other_node": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"other_node_indirect": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"other_node_indirect_soft": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"other_node_soft": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"same_node": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
					"same_node_soft": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},
				},
			},
		},
	}

	return res
}

func DataSourceRgAffinityGroupComputes() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceRgAffinityGroupComputesRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceRgAffinityGroupComputesSchemaMake(),
	}
}
