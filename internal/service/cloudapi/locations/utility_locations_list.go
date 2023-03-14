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

package locations

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/controller"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityLocationsListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (LocationsList, error) {
	locationsList := LocationsList{}
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if page, ok := d.GetOk("page"); ok {
		urlValues.Add("page", strconv.Itoa(page.(int)))
	}
	if size, ok := d.GetOk("size"); ok {
		urlValues.Add("size", strconv.Itoa(size.(int)))
	}

	log.Debugf("utilityLocationsListCheckPresence: load locations list")
	locationsListRaw, err := c.DecortAPICall(ctx, "POST", locationsListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(locationsListRaw), &locationsList)
	if err != nil {
		return nil, err
	}

	return locationsList, nil
}
