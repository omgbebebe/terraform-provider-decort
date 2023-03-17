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

package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
)

func flattenAccountRGList(argl AccountRGList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, arg := range argl {
		temp := map[string]interface{}{
			"computes":     flattenAccRGComputes(arg.Computes),
			"resources":    flattenAccRGResources(arg.Resources),
			"created_by":   arg.CreatedBy,
			"created_time": arg.CreatedTime,
			"deleted_by":   arg.DeletedBy,
			"deleted_time": arg.DeletedTime,
			"rg_id":        arg.RGID,
			"milestones":   arg.Milestones,
			"rg_name":      arg.RGName,
			"status":       arg.Status,
			"updated_by":   arg.UpdatedBy,
			"updated_time": arg.UpdatedTime,
			"vinses":       arg.Vinses,
		}
		res = append(res, temp)
	}
	return res

}

func flattenAccRGComputes(argc AccountRGComputes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"started": argc.Started,
		"stopped": argc.Stopped,
	}
	res = append(res, temp)
	return res
}

func flattenAccResourceHack(r ResourceHack) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":        r.CPU,
		"disksize":   r.Disksize,
		"extips":     r.Extips,
		"exttraffic": r.Exttraffic,
		"gpu":        r.GPU,
		"ram":        r.RAM,
		//"seps":       flattenAccountSeps(r.SEPs),
	}
	res = append(res, temp)
	return res
}

func flattenAccResourceRg(r Resource) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":        r.CPU,
		"disksize":   r.Disksize,
		"extips":     r.Extips,
		"exttraffic": r.Exttraffic,
		"gpu":        r.GPU,
		"ram":        r.RAM,
	}
	res = append(res, temp)
	return res
}

func flattenAccRGResources(argr AccountRGResources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"consumed": flattenAccResourceRg(argr.Consumed),
		"limits":   flattenAccResourceHack(argr.Limits),
		"reserved": flattenAccResourceRg(argr.Reserved),
	}
	res = append(res, temp)
	return res
}

func dataSourceAccountRGListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountRGList, err := utilityAccountRGListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("items", flattenAccountRGList(accountRGList))

	return nil
}

func dataSourceAccountRGListSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the account",
		},
		"items": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Search Result",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"computes": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"started": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"stopped": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"resources": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"consumed": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"cpu": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"disksize": {
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
										},
									},
								},

								"limits": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"cpu": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"disksize": {
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
											"disksize": {
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
										},
									},
								},
							},
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
					"deleted_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"deleted_time": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"rg_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"milestones": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"rg_name": {
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
					"vinses": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
	return res
}

func DataSourceAccountRGList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceAccountRGListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceAccountRGListSchemaMake(),
	}
}
