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

package vins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// On success this function returns a string, as returned by API vins/get, which could be unmarshalled
// into VinsGetResp structure
func utilityVinsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate ViNS by one of the following algorithms depending
	// on the parameters passed:
	//    - if resource group ID is specified -> it looks for a ViNS at the RG level
	//    - if account ID is specifeid -> it looks for a ViNS at the account level
	//
	// If succeeded, it returns non empty string that contains JSON formatted facts about the
	// ViNS as returned by vins/get API call.
	// Otherwise it returns empty string and a meaningful error.
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for the Terraform resource Exists method.
	//

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	// make it possible to use "read" & "check presence" functions with ViNS ID set so
	// that Import of ViNS resource is possible
	idSet := false
	theId, err := strconv.Atoi(d.Id())
	if err != nil || theId <= 0 {
		vinsId, argSet := d.GetOk("vins_id") // NB: vins_id is NOT present in vinsResource schema!
		if argSet {
			theId = vinsId.(int)
			idSet = true
		}
	} else {
		idSet = true
	}

	if idSet {
		// ViNS ID is specified, try to get compute instance straight by this ID
		log.Debugf("utilityVinsCheckPresence: locating ViNS by its ID %d", theId)
		urlValues.Add("vinsId", fmt.Sprintf("%d", theId))
		vinsFacts, err := c.DecortAPICall(ctx, "POST", VinsGetAPI, urlValues)
		if err != nil {
			return "", err
		}
		return vinsFacts, nil
	}

	// ID was not set in the schema upon entering this function - work through ViNS name
	// and Account / RG ID

	vinsName, argSet := d.GetOk("name")
	if !argSet {
		// if ViNS name is not set. then we cannot locate ViNS
		return "", fmt.Errorf("Cannot check ViNS presence if ViNS name is empty")
	}
	urlValues.Add("name", vinsName.(string))
	urlValues.Add("show_all", "false")
	log.Debugf("utilityVinsCheckPresence: preparing to locate ViNS name %s", vinsName.(string))

	rgId, rgSet := d.GetOk("rg_id")
	if rgSet {
		log.Debugf("utilityVinsCheckPresence: limiting ViNS search to RG ID %d", rgId.(int))
		urlValues.Add("rgId", fmt.Sprintf("%d", rgId.(int)))
	}

	accountId, accountSet := d.GetOk("account_id")
	if accountSet {
		log.Debugf("utilityVinsCheckPresence: limiting ViNS search to Account ID %d", accountId.(int))
		urlValues.Add("accountId", fmt.Sprintf("%d", accountId.(int)))
	}

	apiResp, err := c.DecortAPICall(ctx, "POST", VinsSearchAPI, urlValues)
	if err != nil {
		return "", err
	}

	// log.Debugf("%s", apiResp)
	// log.Debugf("utilityResgroupCheckPresence: ready to decode response body from %s", VinsSearchAPI)
	model := VinsSearchResp{}
	err = json.Unmarshal([]byte(apiResp), &model)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityVinsCheckPresence: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		if item.Name == vinsName.(string) {
			if (accountSet && item.AccountID != accountId.(int)) ||
				(rgSet && item.RgID != rgId.(int)) {
				// double check that account ID and Rg ID match, if set in the schema
				continue
			}

			log.Debugf("utilityVinsCheckPresence: match ViNS name %s / ID %d, account ID %d, RG ID %d at index %d",
				item.Name, item.ID, item.AccountID, item.RgID, index)

			// element returned by API vins/search does not contain all information we may need to
			// manage ViNS, so we have to get detailed info by calling API vins/get
			rqValues := &url.Values{}
			rqValues.Add("vinsId", fmt.Sprintf("%d", item.ID))
			vinsGetResp, err := c.DecortAPICall(ctx, "POST", VinsGetAPI, rqValues)
			if err != nil {
				return "", err
			}
			return vinsGetResp, nil
		}
	}

	return "", fmt.Errorf("Cannot find ViNS name %s. Check name and/or RG ID & Account ID and your access rights", vinsName.(string))
}
