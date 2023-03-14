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

package kvmvm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityComputeExtraDisksConfigure(ctx context.Context, d *schema.ResourceData, m interface{}, do_delta bool) error {
	// d is filled with data according to computeResource schema, so extra disks config is retrieved via "extra_disks" key
	// If do_delta is true, this function will identify changes between new and existing specs for extra disks and try to
	// update compute configuration accordingly
	// Otherwise it will apply whatever is found in the new set of "extra_disks" right away.
	// Primary use of do_delta=false is when calling this function from compute Create handler.

	// Note that this function will not abort on API errors, but will continue to configure (attach / detach) other individual
	// disks via atomic API calls. However, it will not retry failed manipulation on the same disk.
	c := m.(*controller.ControllerCfg)

	log.Debugf("utilityComputeExtraDisksConfigure: called for Compute ID %s with do_delta = %t", d.Id(), do_delta)

	// NB: as of rc-1.25 "extra_disks" are TypeSet with the elem of TypeInt
	old_set, new_set := d.GetChange("extra_disks")

	apiErrCount := 0
	var lastSavedError error

	if !do_delta {
		if new_set.(*schema.Set).Len() < 1 {
			return nil
		}

		for _, disk := range new_set.(*schema.Set).List() {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("diskId", fmt.Sprintf("%d", disk.(int)))
			_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAttachAPI, urlValues)
			if err != nil {
				// failed to attach extra disk - partial resource update
				apiErrCount++
				lastSavedError = err
			}
		}

		if apiErrCount > 0 {
			log.Errorf("utilityComputeExtraDisksConfigure: there were %d error(s) when attaching disks to Compute ID %s. Last error was: %s",
				apiErrCount, d.Id(), lastSavedError)
			return lastSavedError
		}

		return nil
	}

	detach_set := old_set.(*schema.Set).Difference(new_set.(*schema.Set))
	log.Debugf("utilityComputeExtraDisksConfigure: detach set has %d items for Compute ID %s", detach_set.Len(), d.Id())

	if detach_set.Len() > 0 {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("force", "false")
		_, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, urlValues)
		if err != nil {
			return err
		}
		for _, diskId := range detach_set.List() {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("diskId", fmt.Sprintf("%d", diskId.(int)))
			_, err := c.DecortAPICall(ctx, "POST", ComputeDiskDetachAPI, urlValues)
			if err != nil {
				// failed to detach disk - there will be partial resource update
				log.Errorf("utilityComputeExtraDisksConfigure: failed to detach disk ID %d from Compute ID %s: %s", diskId.(int), d.Id(), err)
				apiErrCount++
				lastSavedError = err
			}
		}
		urlValues = &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("altBootId", "0")
		_, err = c.DecortAPICall(ctx, "POST", ComputeStartAPI, urlValues)
		if err != nil {
			return err
		}
	}

	attach_set := new_set.(*schema.Set).Difference(old_set.(*schema.Set))
	log.Debugf("utilityComputeExtraDisksConfigure: attach set has %d items for Compute ID %s", attach_set.Len(), d.Id())
	for _, diskId := range attach_set.List() {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("diskId", fmt.Sprintf("%d", diskId.(int)))
		_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAttachAPI, urlValues)
		if err != nil {
			// failed to attach disk - there will be partial resource update
			log.Errorf("utilityComputeExtraDisksConfigure: failed to attach disk ID %d to Compute ID %s: %s", diskId.(int), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	if apiErrCount > 0 {
		log.Errorf("utilityComputeExtraDisksConfigure: there were %d error(s) when managing disks of Compute ID %s. Last error was: %s",
			apiErrCount, d.Id(), lastSavedError)
		return lastSavedError
	}

	return nil
}

func utilityComputeNetworksConfigure(ctx context.Context, d *schema.ResourceData, m interface{}, do_delta bool, skip_zero bool) error {
	// "d" is filled with data according to computeResource schema, so extra networks config is retrieved via "network" key
	// If do_delta is true, this function will identify changes between new and existing specs for network and try to
	// update compute configuration accordingly
	// Otherwise it will apply whatever is found in the new set of "network" right away.
	// Primary use of do_delta=false is when calling this function from compute Create handler.

	c := m.(*controller.ControllerCfg)

	old_set, new_set := d.GetChange("network")

	apiErrCount := 0
	var lastSavedError error

	if !do_delta {
		if new_set.(*schema.Set).Len() < 1 {
			return nil
		}

		for i, runner := range new_set.(*schema.Set).List() {
			if i == 0 && skip_zero {
				continue
			}
			urlValues := &url.Values{}
			net_data := runner.(map[string]interface{})
			urlValues.Add("computeId", d.Id())
			urlValues.Add("netType", net_data["net_type"].(string))
			urlValues.Add("netId", fmt.Sprintf("%d", net_data["net_id"].(int)))
			ipaddr, ipSet := net_data["ip_address"] // "ip_address" key is optional
			if ipSet {
				urlValues.Add("ipAddr", ipaddr.(string))
			}
			_, err := c.DecortAPICall(ctx, "POST", ComputeNetAttachAPI, urlValues)
			if err != nil {
				// failed to attach network - partial resource update
				apiErrCount++
				lastSavedError = err
			}
		}

		if apiErrCount > 0 {
			log.Errorf("utilityComputeNetworksConfigure: there were %d error(s) when managing networks of Compute ID %s. Last error was: %s",
				apiErrCount, d.Id(), lastSavedError)
			return lastSavedError
		}
		return nil
	}

	detach_set := old_set.(*schema.Set).Difference(new_set.(*schema.Set))
	log.Debugf("utilityComputeNetworksConfigure: detach set has %d items for Compute ID %s", detach_set.Len(), d.Id())
	for _, runner := range detach_set.List() {
		urlValues := &url.Values{}
		net_data := runner.(map[string]interface{})
		urlValues.Add("computeId", d.Id())
		urlValues.Add("ipAddr", net_data["ip_address"].(string))
		urlValues.Add("mac", net_data["mac"].(string))
		_, err := c.DecortAPICall(ctx, "POST", ComputeNetDetachAPI, urlValues)
		if err != nil {
			// failed to detach this network - there will be partial resource update
			log.Errorf("utilityComputeNetworksConfigure: failed to detach net ID %d of type %s from Compute ID %s: %s",
				net_data["net_id"].(int), net_data["net_type"].(string), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	attach_set := new_set.(*schema.Set).Difference(old_set.(*schema.Set))
	log.Debugf("utilityComputeNetworksConfigure: attach set has %d items for Compute ID %s", attach_set.Len(), d.Id())
	for _, runner := range attach_set.List() {
		urlValues := &url.Values{}
		net_data := runner.(map[string]interface{})
		urlValues.Add("computeId", d.Id())
		urlValues.Add("netId", fmt.Sprintf("%d", net_data["net_id"].(int)))
		urlValues.Add("netType", net_data["net_type"].(string))
		if net_data["ip_address"].(string) != "" {
			urlValues.Add("ipAddr", net_data["ip_address"].(string))
		}
		_, err := c.DecortAPICall(ctx, "POST", ComputeNetAttachAPI, urlValues)
		if err != nil {
			// failed to attach this network - there will be partial resource update
			log.Errorf("utilityComputeNetworksConfigure: failed to attach net ID %d of type %s to Compute ID %s: %s",
				net_data["net_id"].(int), net_data["net_type"].(string), d.Id(), err)
			apiErrCount++
			lastSavedError = err
		}
	}

	if apiErrCount > 0 {
		log.Errorf("utilityComputeNetworksConfigure: there were %d error(s) when managing networks of Compute ID %s. Last error was: %s",
			apiErrCount, d.Id(), lastSavedError)
		return lastSavedError
	}

	return nil
}

func utilityComputeCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (RecordCompute, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	compute := &RecordCompute{}

	urlValues.Add("computeId", d.Id())
	computeRaw, err := c.DecortAPICall(ctx, "POST", ComputeGetAPI, urlValues)
	if err != nil {
		return *compute, err
	}

	err = json.Unmarshal([]byte(computeRaw), &compute)
	if err != nil {
		return *compute, err
	}
	return *compute, nil
}
