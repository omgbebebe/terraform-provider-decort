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

package vins

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityDataVinsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*VINSDetailed, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	vins := &VINSDetailed{}

	urlValues.Add("vinsId", strconv.Itoa(d.Get("vins_id").(int)))
	vinsRaw, err := c.DecortAPICall(ctx, "POST", VinsGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(vinsRaw), vins)
	if err != nil {
		return nil, err
	}
	return vins, nil

}

func utilityVinsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*VINSDetailed, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	vins := &VINSDetailed{}

	urlValues.Add("vinsId", d.Id())
	vinsRaw, err := c.DecortAPICall(ctx, "POST", VinsGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(vinsRaw), vins)
	if err != nil {
		return nil, err
	}
	return vins, nil

}
