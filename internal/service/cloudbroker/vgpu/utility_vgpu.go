/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>

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

package vgpu

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
)

func utilityVGPUCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*VGPU, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("size", "50")

	var vgpuId int
	var err error

	if vId, ok := d.GetOk("vgpu_id"); ok {
		vgpuId = vId.(int)
	} else {
		vgpuId, err = strconv.Atoi(d.Id())
		if err != nil {
			return nil, err
		}
	}

	for page := 1; ; page++ {
		urlValues.Set("page", strconv.Itoa(page))
		resp, err := c.DecortAPICall(ctx, "POST", vgpuListAPI, urlValues)
		if err != nil {
			return nil, err
		}

		if resp == "[]" {
			return nil, nil
		}

		var vgpus []VGPU
		if err := json.Unmarshal([]byte(resp), &vgpus); err != nil {
			return nil, err
		}

		for _, vgpu := range vgpus {
			if vgpu.ID == vgpuId {
				return &vgpu, nil
			}
		}
	}
}
