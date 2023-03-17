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

package kvmvm

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
)

func dataSourceComputeGetConsoleUrlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	computeConsoleUrl, err := utilityComputeGetConsoleUrlCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id := uuid.New()
	d.SetId(id.String())
	result := strings.ReplaceAll(string(computeConsoleUrl), "\"", "")
	result = strings.ReplaceAll(string(result), "\\", "")
	d.Set("console_url", result)
	return nil
}

func dataSourceComputeGetConsoleUrlSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"console_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceComputeGetConsoleUrl() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceComputeGetConsoleUrlRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceComputeGetConsoleUrlSchemaMake(),
	}
}
