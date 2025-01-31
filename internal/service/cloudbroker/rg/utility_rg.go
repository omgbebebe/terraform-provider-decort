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

package rg

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

// On success this function returns a string, as returned by API rg/get, which could be unmarshalled
// into ResgroupGetResp structure
func utilityResgroupCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	// This function tries to locate resource group by one of the following algorithms depending
	// on the parameters passed:
	//    - if resource group ID is specified -> by RG ID
	//    - if resource group name is specifeid -> by RG name and either account ID or account name
	//
	// If succeeded, it returns non empty string that contains JSON formatted facts about the
	// resource group as returned by rg/get API call.
	// Otherwise it returns empty string and a meaningful error.
	//
	// NOTE: As our provider always deletes RGs permanently, there is no "restore" method and
	// consequently we are not interested in matching RGs in DELETED state. Hence, we call
	// .../rg/list API with includedeleted=false
	//
	// This function does not modify its ResourceData argument, so it is safe to use it as core
	// method for the Terraform resource Exists method.
	//

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	// make it possible to use "read" & "check presence" functions with RG ID set so
	// that Import of RG resource is possible
	idSet := false
	theId, err := strconv.Atoi(d.Id())
	if err != nil || theId <= 0 {
		rgId, argSet := d.GetOk("rg_id")
		if argSet {
			theId = rgId.(int)
			idSet = true
		}
	} else {
		idSet = true
	}

	if idSet {
		// go straight for the RG by its ID
		log.Debugf("utilityResgroupCheckPresence: locating RG by its ID %d", theId)
		urlValues.Add("rgId", fmt.Sprintf("%d", theId))
		rgFacts, err := c.DecortAPICall(ctx, "POST", ResgroupGetAPI, urlValues)
		if err != nil {
			return "", err
		}
		return rgFacts, nil
	}

	rgName, argSet := d.GetOk("name")
	if !argSet {
		// no RG ID and no RG name - we cannot locate resource group in this case
		return "", fmt.Errorf("Cannot check resource group presence if name is empty and no resource group ID specified")
	}

	// Valid account ID is required to locate a resource group
	// obtain Account ID by account name - it should not be zero on success

	urlValues.Add("includedeleted", "false")
	apiResp, err := c.DecortAPICall(ctx, "POST", ResgroupListAPI, urlValues)
	if err != nil {
		return "", err
	}
	// log.Debugf("%s", apiResp)
	log.Debugf("utilityResgroupCheckPresence: ready to decode response body from %s", ResgroupListAPI)
	model := ResgroupListResp{}
	err = json.Unmarshal([]byte(apiResp), &model)
	if err != nil {
		return "", err
	}

	log.Debugf("utilityResgroupCheckPresence: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// match by RG name & account ID
		if item.Name == rgName.(string) && item.AccountID == d.Get("account_id").(int) {
			log.Debugf("utilityResgroupCheckPresence: match RG name %s / ID %d, account ID %d at index %d",
				item.Name, item.ID, item.AccountID, index)

			// not all required information is returned by rg/list API, so we need to initiate one more
			// call to rg/get to obtain extra data to complete Resource population.
			// Namely, we need resource quota settings
			reqValues := &url.Values{}
			reqValues.Add("rgId", fmt.Sprintf("%d", item.ID))
			apiResp, err := c.DecortAPICall(ctx, "POST", ResgroupGetAPI, reqValues)
			if err != nil {
				return "", err
			}

			return apiResp, nil
		}
	}

	return "", fmt.Errorf("Cannot find RG name %s owned by account ID %d", rgName, d.Get("account_id").(int))
}
