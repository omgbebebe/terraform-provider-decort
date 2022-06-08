/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
This file is part of Terraform (by Hashicorp) provider for Digital Energy Cloud Orchestration
Technology platfom.

Visit https://github.com/rudecs/terraform-provider-decort for full source code package and updates.
*/

package decort

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAccountReservedUnitsRead(d *schema.ResourceData, m interface{}) error {
	accountReservedUnits, err := utilityAccountReservedUnitsCheckPresence(d, m)
	if err != nil {
		return err
	}

	id := uuid.New()
	d.SetId(id.String())
	d.Set("cu_c", accountReservedUnits.CUC)
	d.Set("cu_d", accountReservedUnits.CUD)
	d.Set("cu_i", accountReservedUnits.CUI)
	d.Set("cu_m", accountReservedUnits.CUM)
	d.Set("cu_np", accountReservedUnits.CUNP)
	d.Set("gpu_units", accountReservedUnits.GpuUnits)

	return nil
}

func dataSourceAccountReservedUnitsSchemaMake() map[string]*schema.Schema {
	res := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the account",
		},
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

func dataSourceAccountReservedUnits() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read: dataSourceAccountReservedUnitsRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: dataSourceAccountReservedUnitsSchemaMake(),
	}
}
