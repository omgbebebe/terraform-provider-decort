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

package sep

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilitySepPoolCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (SepPool, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	sepPool := SepPool{}

	urlValues.Add("sep_id", strconv.Itoa(d.Get("sep_id").(int)))
	urlValues.Add("pool_name", d.Get("pool_name").(string))

	log.Debugf("utilitySepDesPoolCheckPresence: load sep")
	sepPoolRaw, err := c.DecortAPICall(ctx, "POST", sepGetPoolAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(sepPoolRaw), &sepPool)
	if err != nil {
		return nil, err
	}

	return sepPool, nil
}
