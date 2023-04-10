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
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
)

func existAccountID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	accountList := []struct {
		ID int `json:"id"`
	}{}

	accountListAPI := "/restmachine/cloudapi/account/list"

	accountListRaw, err := c.DecortAPICall(ctx, "POST", accountListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(accountListRaw), &accountList)
	if err != nil {
		return false, err
	}

	haveAccount := false

	myAccount := d.Get("account_id").(int)
	for _, account := range accountList {
		if account.ID == myAccount {
			haveAccount = true
			break
		}
	}
	return haveAccount, nil
}

func existGID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	locationList := []struct {
		GID int `json:"gid"`
	}{}

	locationsListAPI := "/restmachine/cloudapi/locations/list"

	locationListRaw, err := c.DecortAPICall(ctx, "POST", locationsListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(locationListRaw), &locationList)
	if err != nil {
		return false, err
	}

	haveGID := false

	myGID := d.Get("gid").(int)
	for _, location := range locationList {
		if location.GID == myGID {
			haveGID = true
			break
		}
	}

	return haveGID, nil
}

func existExtNetID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))

	listExtNet := []struct {
		ID int `json:"id"`
	}{}

	extNetListAPI := "/restmachine/cloudapi/extnet/list"

	listExtNetRaw, err := c.DecortAPICall(ctx, "POST", extNetListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(listExtNetRaw), &listExtNet)
	if err != nil {
		return false, err
	}

	haveExtNet := false

	myExtNetID := d.Get("ext_net_id").(int)
	for _, extNet := range listExtNet {
		if extNet.ID == myExtNetID {
			haveExtNet = true
			break
		}
	}
	return haveExtNet, nil
}
