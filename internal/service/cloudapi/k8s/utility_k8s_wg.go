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

package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/service/cloudapi/kvmvm"
)

func utilityDataK8sWgCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*K8SGroup, []kvmvm.ComputeGetResp, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	k8sId := d.Get("k8s_id").(int)
	wgId := d.Get("wg_id").(int)

	urlValues.Add("k8sId", strconv.Itoa(k8sId))

	k8sRaw, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, nil, err
	}

	k8s := K8SRecord{}
	err = json.Unmarshal([]byte(k8sRaw), &k8s)
	if err != nil {
		return nil, nil, err
	}

	curWg := K8SGroup{}
	for _, wg := range k8s.K8SGroups.Workers {
		if wg.ID == uint64(wgId) {
			curWg = wg
			break
		}
	}
	if curWg.ID == 0 {
		return nil, nil, fmt.Errorf("WG with id %v in k8s cluster %v not found", wgId, k8sId)
	}

	workersComputeList := make([]kvmvm.ComputeGetResp, 0, 0)
	for _, info := range curWg.DetailedInfo {
		compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
		if err != nil {
			return nil, nil, err
		}

		workersComputeList = append(workersComputeList, *compute)
	}

	return &curWg, workersComputeList, nil
}

func utilityK8sWgCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*K8SGroup, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	var wgId int
	var k8sId int
	var err error

	if strings.Contains(d.Id(), "#") {
		wgId, err = strconv.Atoi(strings.Split(d.Id(), "#")[0])
		if err != nil {
			return nil, err
		}
		k8sId, err = strconv.Atoi(strings.Split(d.Id(), "#")[1])
		if err != nil {
			return nil, err
		}
	} else {
		wgId, err = strconv.Atoi(d.Id())
		if err != nil {
			return nil, err
		}
		k8sId = d.Get("k8s_id").(int)
	}

	urlValues.Add("k8sId", strconv.Itoa(k8sId))

	resp, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, err
	}

	var k8s K8SRecord
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	for _, wg := range k8s.K8SGroups.Workers {
		if wg.ID == uint64(wgId) {
			return &wg, nil
		}
	}

	return nil, fmt.Errorf("Not found wg with id: %v in k8s cluster: %v", wgId, k8s.ID)
}
