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

package k8s

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/constants"
)

func dataSourceK8sListDeletedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	k8sList, err := utilityK8sListCheckPresence(ctx, d, m, K8sListDeletedAPI)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	flattenK8sList(d, k8sList)

	return nil
}

func dataSourceK8sListDeletedSchemaMake() map[string]*schema.Schema {
	k8sListDeleted := createK8sListSchema()
	delete(k8sListDeleted, "includedeleted")
	return k8sListDeleted
}

func DataSourceK8sListDeleted() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sListDeletedRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sListDeletedSchemaMake(),
	}
}
