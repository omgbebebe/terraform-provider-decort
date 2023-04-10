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

package kvmvm

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/dc"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/statefuncs"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/status"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceComputeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// we assume all mandatory parameters it takes to create a comptue instance are properly
	// specified - we rely on schema "Required" attributes to let Terraform validate them for us

	log.Debugf("resourceComputeCreate: called for Compute name %q, RG ID %d", d.Get("name").(string), d.Get("rg_id").(int))
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if !existRgID(ctx, d, m) {
		return diag.Errorf("resourceComputeCreate: can't create Compute bacause rgID %d not allowed or does not exist", d.Get("rg_id").(int))
	}

	if !existImageId(ctx, d, m) {
		return diag.Errorf("resourceComputeCreate: can't create Compute bacause imageID %d not allowed or does not exist", d.Get("image_id").(int))
	}

	if _, ok := d.GetOk("network"); ok {
		if vinsId, ok := existVinsId(ctx, d, m); !ok {
			return diag.Errorf("resourceResgroupCreate: can't create RG bacause vins ID %d not allowed or does not exist", vinsId)
		}
	}

	// create basic Compute (i.e. without extra disks and network connections - those will be attached
	// by subsequent individual API calls).
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("cpu", strconv.Itoa(d.Get("cpu").(int)))
	urlValues.Add("ram", strconv.Itoa(d.Get("ram").(int)))
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("netType", "NONE")
	urlValues.Add("start", "0") // at the 1st step create compute in a stopped state

	argVal, ok := d.GetOk("description")
	if ok {
		urlValues.Add("desc", argVal.(string))
	}

	if sepID, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sepId", strconv.Itoa(sepID.(int)))
	}

	if pool, ok := d.GetOk("pool"); ok {
		urlValues.Add("pool", pool.(string))
	}

	if ipaType, ok := d.GetOk("ipa_type"); ok {
		urlValues.Add("ipaType", ipaType.(string))
	}

	if bootSize, ok := d.GetOk("boot_disk_size"); ok {
		urlValues.Add("bootDisk", fmt.Sprintf("%d", bootSize.(int)))
	}

	if IS, ok := d.GetOk("is"); ok {
		urlValues.Add("IS", IS.(string))
	}
	if networks, ok := d.GetOk("network"); ok {
		if networks.(*schema.Set).Len() > 0 {
			ns := networks.(*schema.Set).List()
			defaultNetwork := ns[0].(map[string]interface{})
			urlValues.Set("netType", defaultNetwork["net_type"].(string))
			urlValues.Add("netId", fmt.Sprintf("%d", defaultNetwork["net_id"].(int)))
			ipaddr, ipSet := defaultNetwork["ip_address"] // "ip_address" key is optional
			if ipSet {
				urlValues.Add("ipAddr", ipaddr.(string))
			}

		}
	}

	/*
		sshKeysVal, sshKeysSet := d.GetOk("ssh_keys")
		if sshKeysSet {
			// process SSH Key settings and set API values accordingly
			log.Debugf("resourceComputeCreate: calling makeSshKeysArgString to setup SSH keys for guest login(s)")
			urlValues.Add("userdata", makeSshKeysArgString(sshKeysVal.([]interface{})))
		}
	*/

	computeCreateAPI := KvmX86CreateAPI
	driver := d.Get("driver").(string)
	if driver == "KVM_PPC" {
		computeCreateAPI = KvmPPCCreateAPI
		log.Debugf("resourceComputeCreate: creating Compute of type KVM VM PowerPC")
	} else { // note that we do not validate arch value for explicit "KVM_X86" here
		log.Debugf("resourceComputeCreate: creating Compute of type KVM VM x86")
	}

	argVal, ok = d.GetOk("cloud_init")
	if ok {
		// userdata must not be empty string and must not be a reserved keyword "applied"
		userdata := argVal.(string)
		if userdata != "" && userdata != "applied" {
			urlValues.Add("userdata", userdata)
		}
	}

	apiResp, err := c.DecortAPICall(ctx, "POST", computeCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	urlValues = &url.Values{}
	// Compute create API returns ID of the new Compute instance on success

	d.SetId(apiResp) // update ID of the resource to tell Terraform that the resource exists, albeit partially
	compId, _ := strconv.Atoi(apiResp)

	warnings := dc.Warnings{}

	cleanup := false
	defer func() {
		if cleanup {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("permanently", "1")
			urlValues.Add("detachDisks", "1")

			if _, err := c.DecortAPICall(ctx, "POST", ComputeDeleteAPI, urlValues); err != nil {
				log.Errorf("resourceComputeCreate: could not delete compute after failed creation: %v", err)
			}
			d.SetId("")
			urlValues = &url.Values{}
		}
	}()

	log.Debugf("resourceComputeCreate: new simple Compute ID %d, name %s created", compId, d.Get("name").(string))

	argVal, ok = d.GetOk("extra_disks")
	if ok && argVal.(*schema.Set).Len() > 0 {
		log.Debugf("resourceComputeCreate: calling utilityComputeExtraDisksConfigure to attach %d extra disk(s)", argVal.(*schema.Set).Len())
		err = utilityComputeExtraDisksConfigure(ctx, d, m, false) // do_delta=false, as we are working on a new compute
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching extra disk(s) to a new Compute ID %d: %v", compId, err)
			cleanup = true
			return diag.FromErr(err)
		}
	}
	argVal, ok = d.GetOk("network")
	if ok && argVal.(*schema.Set).Len() > 0 {
		log.Debugf("resourceComputeCreate: calling utilityComputeNetworksConfigure to attach %d network(s)", argVal.(*schema.Set).Len())
		err = utilityComputeNetworksConfigure(ctx, d, m, false, true)
		if err != nil {
			log.Errorf("resourceComputeCreate: error when attaching networks to a new Compute ID %d: %s", compId, err)
			cleanup = true
			return diag.FromErr(err)
		}
	}

	// Note bene: we created compute in a STOPPED state (this is required to properly attach 1st network interface),
	// now we need to start it before we report the sequence complete
	if d.Get("started").(bool) {
		reqValues := &url.Values{}
		reqValues.Add("computeId", fmt.Sprintf("%d", compId))
		log.Debugf("resourceComputeCreate: starting Compute ID %d after completing its resource configuration", compId)
		if _, err := c.DecortAPICall(ctx, "POST", ComputeStartAPI, reqValues); err != nil {
			warnings.Add(err)
		}
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		api := ComputeDisableAPI
		if enabled.(bool) {
			api = ComputeEnableAPI
		}
		urlValues := &url.Values{}
		urlValues.Add("computeId", fmt.Sprintf("%d", compId))
		log.Debugf("resourceComputeCreate: enable=%t Compute ID %d after completing its resource configuration", compId, enabled)
		if _, err := c.DecortAPICall(ctx, "POST", api, urlValues); err != nil {
			warnings.Add(err)
		}

	}

	if !cleanup {
		if affinityLabel, ok := d.GetOk("affinity_label"); ok {
			affinityLabel := affinityLabel.(string)
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("affinityLabel", affinityLabel)
			_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityLabelSetAPI, urlValues)
			if err != nil {
				warnings.Add(err)
			}
			urlValues = &url.Values{}
		}

		if disks, ok := d.GetOk("disks"); ok {
			log.Debugf("resourceComputeCreate: Create disks on ComputeID: %d", compId)
			addedDisks := disks.([]interface{})
			if len(addedDisks) > 0 {
				for _, disk := range addedDisks {
					diskConv := disk.(map[string]interface{})

					urlValues.Add("computeId", d.Id())
					urlValues.Add("diskName", diskConv["disk_name"].(string))
					urlValues.Add("size", strconv.Itoa(diskConv["size"].(int)))
					if diskConv["disk_type"].(string) != "" {
						urlValues.Add("diskType", diskConv["disk_type"].(string))
					}
					if diskConv["sep_id"].(int) != 0 {
						urlValues.Add("sepId", strconv.Itoa(diskConv["sep_id"].(int)))
					}
					if diskConv["pool"].(string) != "" {
						urlValues.Add("pool", diskConv["pool"].(string))
					}
					if diskConv["desc"].(string) != "" {
						urlValues.Add("desc", diskConv["desc"].(string))
					}
					if diskConv["image_id"].(int) != 0 {
						urlValues.Add("imageId", strconv.Itoa(diskConv["image_id"].(int)))
					}
					_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAddAPI, urlValues)
					if err != nil {
						cleanup = true
						return diag.FromErr(err)
					}
					urlValues = &url.Values{}
				}
			}
		}

		if ars, ok := d.GetOk("affinity_rules"); ok {
			log.Debugf("resourceComputeCreate: Create affinity rules on ComputeID: %d", compId)
			addedAR := ars.([]interface{})
			if len(addedAR) > 0 {
				for _, ar := range addedAR {
					arConv := ar.(map[string]interface{})

					urlValues.Add("computeId", d.Id())
					urlValues.Add("topology", arConv["topology"].(string))
					urlValues.Add("policy", arConv["policy"].(string))
					urlValues.Add("mode", arConv["mode"].(string))
					urlValues.Add("key", arConv["key"].(string))
					urlValues.Add("value", arConv["value"].(string))
					_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityRuleAddAPI, urlValues)
					if err != nil {
						warnings.Add(err)
					}
					urlValues = &url.Values{}
				}
			}
		}

		if ars, ok := d.GetOk("anti_affinity_rules"); ok {
			log.Debugf("resourceComputeCreate: Create anti affinity rules on ComputeID: %d", compId)
			addedAR := ars.([]interface{})
			if len(addedAR) > 0 {
				for _, ar := range addedAR {
					arConv := ar.(map[string]interface{})

					urlValues.Add("computeId", d.Id())
					urlValues.Add("topology", arConv["topology"].(string))
					urlValues.Add("policy", arConv["policy"].(string))
					urlValues.Add("mode", arConv["mode"].(string))
					urlValues.Add("key", arConv["key"].(string))
					urlValues.Add("value", arConv["value"].(string))
					_, err := c.DecortAPICall(ctx, "POST", ComputeAntiAffinityRuleAddAPI, urlValues)
					if err != nil {
						warnings.Add(err)
					}
					urlValues = &url.Values{}
				}
			}
		}
	}

	if tags, ok := d.GetOk("tags"); ok {
		log.Debugf("resourceComputeCreate: Create tags on ComputeID: %d", compId)
		addedTags := tags.(*schema.Set).List()
		if len(addedTags) > 0 {
			for _, tagInterface := range addedTags {
				urlValues = &url.Values{}
				tagItem := tagInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("key", tagItem["key"].(string))
				urlValues.Add("value", tagItem["value"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeTagAddAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}

	if pfws, ok := d.GetOk("port_forwarding"); ok {
		log.Debugf("resourceComputeCreate: Create port farwarding on ComputeID: %d", compId)
		addedPfws := pfws.(*schema.Set).List()
		if len(addedPfws) > 0 {
			for _, pfwInterface := range addedPfws {
				urlValues = &url.Values{}
				pfwItem := pfwInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("publicPortStart", strconv.Itoa(pfwItem["public_port_start"].(int)))
				urlValues.Add("publicPortEnd", strconv.Itoa(pfwItem["public_port_end"].(int)))
				urlValues.Add("localBasePort", strconv.Itoa(pfwItem["local_port"].(int)))
				urlValues.Add("proto", pfwItem["proto"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputePfwAddAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}
	if userAcess, ok := d.GetOk("user_access"); ok {
		log.Debugf("resourceComputeCreate: Create user access on ComputeID: %d", compId)
		usersAcess := userAcess.(*schema.Set).List()
		if len(usersAcess) > 0 {
			for _, userAcessInterface := range usersAcess {
				urlValues = &url.Values{}
				userAccessItem := userAcessInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("userName", userAccessItem["username"].(string))
				urlValues.Add("accesstype", userAccessItem["access_type"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeUserGrantAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}

	if snapshotList, ok := d.GetOk("snapshot"); ok {
		log.Debugf("resourceComputeCreate: Create snapshot on ComputeID: %d", compId)
		snapshots := snapshotList.(*schema.Set).List()
		if len(snapshots) > 0 {
			for _, snapshotInterface := range snapshots {
				urlValues = &url.Values{}
				snapshotItem := snapshotInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("userName", snapshotItem["label"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeSnapshotCreateAPI, urlValues)
				if err != nil {
					warnings.Add(err)
				}
			}
		}
	}

	if cdtList, ok := d.GetOk("cd"); ok {
		log.Debugf("resourceComputeCreate: Create cd on ComputeID: %d", compId)
		cds := cdtList.(*schema.Set).List()
		if len(cds) > 0 {
			urlValues = &url.Values{}
			snapshotItem := cds[0].(map[string]interface{})

			urlValues.Add("computeId", d.Id())
			urlValues.Add("cdromId", strconv.Itoa(snapshotItem["cdrom_id"].(int)))
			_, err := c.DecortAPICall(ctx, "POST", ComputeCdInsertAPI, urlValues)
			if err != nil {
				warnings.Add(err)
			}
		}
	}

	if d.Get("pin_to_stack").(bool) == true {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", ComputePinToStackAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}

	if d.Get("pause").(bool) == true {
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", ComputePauseAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	}

	log.Debugf("resourceComputeCreate: new Compute ID %d, name %s creation sequence complete", compId, d.Get("name").(string))

	// We may reuse dataSourceComputeRead here as we maintain similarity
	// between Compute resource and Compute data source schemas
	// Compute read function will also update resource ID on success, so that Terraform
	// will know the resource exists
	defer resourceComputeRead(ctx, d, m)
	return warnings.Get()
}

func resourceComputeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceComputeRead: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	compute, err := utilityComputeCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	hasChanged := false

	switch compute.Status {
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", ComputeRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", ComputeEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Destroyed:
		d.SetId("")
		return resourceComputeCreate(ctx, d, m)
	case status.Disabled:
		log.Debugf("The compute is in status: %s, troubles may occur with update. Please, enable compute first.", compute.Status)
	case status.Redeploying:
	case status.Deleting:
	case status.Destroying:
		return diag.Errorf("The compute is in progress with status: %s", compute.Status)
	case status.Modeled:
		return diag.Errorf("The compute is in status: %s, please, contact support for more information", compute.Status)
	}

	if hasChanged {
		compute, err = utilityComputeCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	if err = flattenCompute(d, compute); err != nil {
		return diag.FromErr(err)
	}

	log.Debugf("resourceComputeRead: after flattenCompute: Compute ID %s, name %q, RG ID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	return nil
}

func resourceComputeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceComputeUpdate: called for Compute ID %s / name %s, RGID %d",
		d.Id(), d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	if !existRgID(ctx, d, m) {
		return diag.Errorf("resourceComputeUpdate: can't update Compute bacause rgID %d not allowed or does not exist", d.Get("rg_id").(int))
	}

	if !existImageId(ctx, d, m) {
		return diag.Errorf("resourceComputeUpdate: can't update Compute bacause imageID %d not allowed or does not exist", d.Get("image_id").(int))
	}

	if _, ok := d.GetOk("network"); ok {
		if vinsId, ok := existVinsId(ctx, d, m); !ok {
			return diag.Errorf("resourceResgroupUpdate: can't update RG bacause vinsID %d not allowed or does not exist", vinsId)
		}
	}

	compute, err := utilityComputeCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled")
		api := ComputeDisableAPI
		if enabled.(bool) {
			api = ComputeEnableAPI
		}
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		log.Debugf("resourceComputeUpdate: enable=%t Compute ID %s after completing its resource configuration", d.Id(), enabled)
		if _, err := c.DecortAPICall(ctx, "POST", api, urlValues); err != nil {
			return diag.FromErr(err)
		}
	}

	// check compute statuses
	switch compute.Status {
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", ComputeRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", ComputeEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	case status.Destroyed:
		d.SetId("")
		return resourceComputeCreate(ctx, d, m)
	case status.Disabled:
		log.Debugf("The compute is in status: %s, may troubles can be occured with update. Please, enable compute first.", compute.Status)
	case status.Redeploying:
	case status.Deleting:
	case status.Destroying:
		return diag.Errorf("The compute is in progress with status: %s", compute.Status)
	case status.Modeled:
		return diag.Errorf("The compute is in status: %s, please, contant the support for more information", compute.Status)
	}

	/*
		1. Resize CPU/RAM
		2. Resize (grow) boot disk
		3. Update extra disks
		4. Update networks
		5. Start/stop
	*/

	// 1. Resize CPU/RAM
	urlValues = &url.Values{}
	doUpdate := false
	urlValues.Add("computeId", d.Id())

	oldCpu, newCpu := d.GetChange("cpu")
	if oldCpu.(int) != newCpu.(int) {
		urlValues.Add("cpu", fmt.Sprintf("%d", newCpu.(int)))
		doUpdate = true
	} else {
		urlValues.Add("cpu", "0") // no change to CPU allocation
	}

	oldRam, newRam := d.GetChange("ram")
	if oldRam.(int) != newRam.(int) {
		urlValues.Add("ram", fmt.Sprintf("%d", newRam.(int)))
		doUpdate = true
	} else {
		urlValues.Add("ram", "0")
	}

	if doUpdate {
		log.Debugf("resourceComputeUpdate: changing CPU %d -> %d and/or RAM %d -> %d",
			oldCpu.(int), newCpu.(int),
			oldRam.(int), newRam.(int))
		urlValues.Add("force", "true")
		_, err := c.DecortAPICall(ctx, "POST", ComputeResizeAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// 2. Resize (grow) Boot disk
	oldSize, newSize := d.GetChange("boot_disk_size")
	if oldSize.(int) < newSize.(int) {
		urlValues := &url.Values{}
		if diskId, ok := d.GetOk("boot_disk_id"); ok {
			urlValues.Add("diskId", strconv.Itoa(diskId.(int)))
		} else {
			bootDisk, err := utilityComputeBootDiskCheckPresence(ctx, d, m)
			if err != nil {
				return diag.FromErr(err)
			}
			urlValues.Add("diskId", strconv.FormatUint(bootDisk.ID, 10))
		}
		urlValues.Add("size", strconv.Itoa(newSize.(int)))
		log.Debugf("resourceComputeUpdate: compute ID %s, boot disk ID %d resize %d -> %d",
			d.Id(), d.Get("boot_disk_id").(int), oldSize.(int), newSize.(int))
		_, err := c.DecortAPICall(ctx, "POST", DisksResizeAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if oldSize.(int) > newSize.(int) {
		log.Warnf("resourceComputeUpdate: compute ID %s - shrinking boot disk is not allowed", d.Id())
	}

	// 3. Calculate and apply changes to data disks
	if d.HasChange("extra_disks") {
		err := utilityComputeExtraDisksConfigure(ctx, d, m, true) // pass do_delta = true to apply changes, if any
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// 4. Calculate and apply changes to network connections
	err = utilityComputeNetworksConfigure(ctx, d, m, true, false) // pass do_delta = true to apply changes, if any
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("description") || d.HasChange("name") {
		updateParams := &url.Values{}
		updateParams.Add("computeId", d.Id())
		updateParams.Add("name", d.Get("name").(string))
		updateParams.Add("desc", d.Get("description").(string))
		if _, err := c.DecortAPICall(ctx, "POST", ComputeUpdateAPI, updateParams); err != nil {
			return diag.FromErr(err)
		}
	}

	urlValues = &url.Values{}
	if d.HasChange("disks") {
		deletedDisks := make([]interface{}, 0)
		addedDisks := make([]interface{}, 0)
		updatedDisks := make([]interface{}, 0)

		oldDisks, newDisks := d.GetChange("disks")
		oldConv := oldDisks.([]interface{})
		newConv := newDisks.([]interface{})

		for _, el := range oldConv {
			if !isContainsDisk(newConv, el) {
				deletedDisks = append(deletedDisks, el)
			}
		}

		for _, el := range newConv {
			if !isContainsDisk(oldConv, el) {
				addedDisks = append(addedDisks, el)
			} else {
				if isChangeDisk(oldConv, el) {
					updatedDisks = append(updatedDisks, el)
				}
			}
		}

		if len(deletedDisks) > 0 {
			urlValues.Add("computeId", d.Id())
			urlValues.Add("force", "false")
			_, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
			urlValues = &url.Values{}

			for _, disk := range deletedDisks {
				diskConv := disk.(map[string]interface{})
				if diskConv["disk_name"].(string) == "bootdisk" {
					continue
				}
				urlValues.Add("computeId", d.Id())
				urlValues.Add("diskId", strconv.Itoa(diskConv["disk_id"].(int)))
				urlValues.Add("permanently", strconv.FormatBool(diskConv["permanently"].(bool)))
				_, err := c.DecortAPICall(ctx, "POST", ComputeDiskDeleteAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
			urlValues.Add("computeId", d.Id())
			urlValues.Add("altBootId", "0")
			_, err = c.DecortAPICall(ctx, "POST", ComputeStartAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
			urlValues = &url.Values{}
		}

		if len(addedDisks) > 0 {
			for _, disk := range addedDisks {
				diskConv := disk.(map[string]interface{})
				if diskConv["disk_name"].(string) == "bootdisk" {
					continue
				}
				urlValues.Add("computeId", d.Id())
				urlValues.Add("diskName", diskConv["disk_name"].(string))
				urlValues.Add("size", strconv.Itoa(diskConv["size"].(int)))
				if diskConv["disk_type"].(string) != "" {
					urlValues.Add("diskType", diskConv["disk_type"].(string))
				}
				if diskConv["sep_id"].(int) != 0 {
					urlValues.Add("sepId", strconv.Itoa(diskConv["sep_id"].(int)))
				}
				if diskConv["pool"].(string) != "" {
					urlValues.Add("pool", diskConv["pool"].(string))
				}
				if diskConv["desc"].(string) != "" {
					urlValues.Add("desc", diskConv["desc"].(string))
				}
				if diskConv["image_id"].(int) != 0 {
					urlValues.Add("imageId", strconv.Itoa(diskConv["image_id"].(int)))
				}
				_, err := c.DecortAPICall(ctx, "POST", ComputeDiskAddAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

		if len(updatedDisks) > 0 {
			for _, disk := range updatedDisks {
				diskConv := disk.(map[string]interface{})
				if diskConv["disk_name"].(string) == "bootdisk" {
					continue
				}
				urlValues = &url.Values{}
				urlValues.Add("diskId", strconv.Itoa(diskConv["disk_id"].(int)))
				urlValues.Add("size", strconv.Itoa(diskConv["size"].(int)))
				_, err := c.DecortAPICall(ctx, "POST", DisksResizeAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("started") {
		params := &url.Values{}
		params.Add("computeId", d.Id())
		if d.Get("started").(bool) {
			if _, err := c.DecortAPICall(ctx, "POST", ComputeStartAPI, params); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if _, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, params); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("affinity_label") {
		affinityLabel := d.Get("affinity_label").(string)
		urlValues.Add("computeId", d.Id())
		if affinityLabel == "" {
			_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityLabelRemoveAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		urlValues.Add("affinityLabel", affinityLabel)
		_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityLabelSetAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		urlValues = &url.Values{}
	}

	if d.HasChange("affinity_rules") {
		deletedAR := make([]interface{}, 0)
		addedAR := make([]interface{}, 0)

		oldAR, newAR := d.GetChange("affinity_rules")
		oldConv := oldAR.([]interface{})
		newConv := newAR.([]interface{})

		if len(newConv) == 0 {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityRulesClearAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			for _, el := range oldConv {
				if !isContainsAR(newConv, el) {
					deletedAR = append(deletedAR, el)
				}
			}
			for _, el := range newConv {
				if !isContainsAR(oldConv, el) {
					addedAR = append(addedAR, el)
				}
			}

			if len(deletedAR) > 0 {
				urlValues := &url.Values{}
				for _, ar := range deletedAR {
					arConv := ar.(map[string]interface{})
					urlValues.Add("computeId", d.Id())
					urlValues.Add("topology", arConv["topology"].(string))
					urlValues.Add("policy", arConv["policy"].(string))
					urlValues.Add("mode", arConv["mode"].(string))
					urlValues.Add("key", arConv["key"].(string))
					urlValues.Add("value", arConv["value"].(string))
					_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityRuleRemoveAPI, urlValues)
					if err != nil {
						return diag.FromErr(err)
					}

					urlValues = &url.Values{}
				}
			}
			if len(addedAR) > 0 {
				for _, ar := range addedAR {
					arConv := ar.(map[string]interface{})
					urlValues.Add("computeId", d.Id())
					urlValues.Add("topology", arConv["topology"].(string))
					urlValues.Add("policy", arConv["policy"].(string))
					urlValues.Add("mode", arConv["mode"].(string))
					urlValues.Add("key", arConv["key"].(string))
					urlValues.Add("value", arConv["value"].(string))
					_, err := c.DecortAPICall(ctx, "POST", ComputeAffinityRuleAddAPI, urlValues)
					if err != nil {
						return diag.FromErr(err)
					}

					urlValues = &url.Values{}
				}
			}
		}

	}

	if d.HasChange("anti_affinity_rules") {
		deletedAR := make([]interface{}, 0)
		addedAR := make([]interface{}, 0)

		oldAR, newAR := d.GetChange("anti_affinity_rules")
		oldConv := oldAR.([]interface{})
		newConv := newAR.([]interface{})

		if len(newConv) == 0 {
			urlValues := &url.Values{}
			urlValues.Add("computeId", d.Id())
			_, err := c.DecortAPICall(ctx, "POST", ComputeAntiAffinityRulesClearAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			for _, el := range oldConv {
				if !isContainsAR(newConv, el) {
					deletedAR = append(deletedAR, el)
				}
			}
			for _, el := range newConv {
				if !isContainsAR(oldConv, el) {
					addedAR = append(addedAR, el)
				}
			}

			if len(deletedAR) > 0 {
				urlValues := &url.Values{}
				for _, ar := range deletedAR {
					arConv := ar.(map[string]interface{})
					urlValues.Add("computeId", d.Id())
					urlValues.Add("topology", arConv["topology"].(string))
					urlValues.Add("policy", arConv["policy"].(string))
					urlValues.Add("mode", arConv["mode"].(string))
					urlValues.Add("key", arConv["key"].(string))
					urlValues.Add("value", arConv["value"].(string))
					_, err := c.DecortAPICall(ctx, "POST", ComputeAntiAffinityRuleRemoveAPI, urlValues)
					if err != nil {
						return diag.FromErr(err)
					}

					urlValues = &url.Values{}
				}
			}
			if len(addedAR) > 0 {
				for _, ar := range addedAR {
					arConv := ar.(map[string]interface{})
					urlValues.Add("computeId", d.Id())
					urlValues.Add("topology", arConv["topology"].(string))
					urlValues.Add("policy", arConv["policy"].(string))
					urlValues.Add("mode", arConv["mode"].(string))
					urlValues.Add("key", arConv["key"].(string))
					urlValues.Add("value", arConv["value"].(string))
					_, err := c.DecortAPICall(ctx, "POST", ComputeAntiAffinityRuleAddAPI, urlValues)
					if err != nil {
						return diag.FromErr(err)
					}

					urlValues = &url.Values{}
				}
			}
		}

	}

	if d.HasChange("tags") {
		oldSet, newSet := d.GetChange("tags")
		deletedTags := (oldSet.(*schema.Set).Difference(newSet.(*schema.Set))).List()
		if len(deletedTags) > 0 {
			for _, tagInterface := range deletedTags {
				urlValues := &url.Values{}
				tagItem := tagInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("key", tagItem["key"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeTagRemoveAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		addedTags := (newSet.(*schema.Set).Difference(oldSet.(*schema.Set))).List()
		if len(addedTags) > 0 {
			for _, tagInterface := range addedTags {
				urlValues := &url.Values{}
				tagItem := tagInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("key", tagItem["key"].(string))
				urlValues.Add("value", tagItem["value"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeTagAddAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("port_forwarding") {
		oldSet, newSet := d.GetChange("port_forwarding")
		deletedPfws := (oldSet.(*schema.Set).Difference(newSet.(*schema.Set))).List()
		if len(deletedPfws) > 0 {
			for _, pfwInterface := range deletedPfws {
				urlValues := &url.Values{}
				pfwItem := pfwInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("publicPortStart", strconv.Itoa(pfwItem["public_port_start"].(int)))
				if pfwItem["public_port_end"].(int) == -1 {
					urlValues.Add("publicPortEnd", strconv.Itoa(pfwItem["public_port_start"].(int)))
				} else {
					urlValues.Add("publicPortEnd", strconv.Itoa(pfwItem["public_port_end"].(int)))
				}
				urlValues.Add("localBasePort", strconv.Itoa(pfwItem["local_port"].(int)))
				urlValues.Add("proto", pfwItem["proto"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputePfwDelAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		addedPfws := (newSet.(*schema.Set).Difference(oldSet.(*schema.Set))).List()
		if len(addedPfws) > 0 {
			for _, pfwInterface := range addedPfws {
				urlValues := &url.Values{}
				pfwItem := pfwInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("publicPortStart", strconv.Itoa(pfwItem["public_port_start"].(int)))
				urlValues.Add("publicPortEnd", strconv.Itoa(pfwItem["public_port_end"].(int)))
				urlValues.Add("localBasePort", strconv.Itoa(pfwItem["local_port"].(int)))
				urlValues.Add("proto", pfwItem["proto"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputePfwAddAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("user_access") {
		oldSet, newSet := d.GetChange("user_access")
		deletedUserAcess := (oldSet.(*schema.Set).Difference(newSet.(*schema.Set))).List()
		if len(deletedUserAcess) > 0 {
			for _, userAcessInterface := range deletedUserAcess {
				urlValues := &url.Values{}
				userAccessItem := userAcessInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("userName", userAccessItem["username"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeUserRevokeAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		addedUserAccess := (newSet.(*schema.Set).Difference(oldSet.(*schema.Set))).List()
		if len(addedUserAccess) > 0 {
			for _, userAccessInterface := range addedUserAccess {
				urlValues := &url.Values{}
				userAccessItem := userAccessInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("userName", userAccessItem["username"].(string))
				urlValues.Add("accesstype", userAccessItem["access_type"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeUserGrantAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("snapshot") {
		oldSet, newSet := d.GetChange("snapshot")
		deletedSnapshots := (oldSet.(*schema.Set).Difference(newSet.(*schema.Set))).List()
		if len(deletedSnapshots) > 0 {
			for _, snapshotInterface := range deletedSnapshots {
				urlValues := &url.Values{}
				snapshotItem := snapshotInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("label", snapshotItem["label"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeSnapshotDeleteAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		addedSnapshots := (newSet.(*schema.Set).Difference(oldSet.(*schema.Set))).List()
		if len(addedSnapshots) > 0 {
			for _, snapshotInterface := range addedSnapshots {
				urlValues := &url.Values{}
				snapshotItem := snapshotInterface.(map[string]interface{})

				urlValues.Add("computeId", d.Id())
				urlValues.Add("label", snapshotItem["label"].(string))
				_, err := c.DecortAPICall(ctx, "POST", ComputeSnapshotCreateAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("rollback") {
		if rollback, ok := d.GetOk("rollback"); ok {
			urlValues := &url.Values{}

			//Compute must be stopped before rollback
			urlValues.Add("computeId", d.Id())
			urlValues.Add("force", "false")
			_, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
			rollbackInterface := rollback.(*schema.Set).List()[0]
			rollbackItem := rollbackInterface.(map[string]interface{})

			urlValues.Add("computeId", d.Id())
			urlValues.Add("label", rollbackItem["label"].(string))
			_, err = c.DecortAPICall(ctx, "POST", ComputeSnapshotRollbackAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("cd") {
		oldSet, newSet := d.GetChange("cd")
		deletedCd := (oldSet.(*schema.Set).Difference(newSet.(*schema.Set))).List()
		if len(deletedCd) > 0 {
			urlValues := &url.Values{}

			urlValues.Add("computeId", d.Id())
			_, err := c.DecortAPICall(ctx, "POST", ComputeCdEjectAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		addedCd := (newSet.(*schema.Set).Difference(oldSet.(*schema.Set))).List()
		if len(addedCd) > 0 {
			urlValues := &url.Values{}
			cdItem := addedCd[0].(map[string]interface{})

			urlValues.Add("computeId", d.Id())
			urlValues.Add("cdromId", strconv.Itoa(cdItem["cdrom_id"].(int)))
			_, err := c.DecortAPICall(ctx, "POST", ComputeCdInsertAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("pin_to_stack") {
		oldPin, newPin := d.GetChange("pin_to_stack")
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		if oldPin.(bool) == true && newPin.(bool) == false {
			_, err := c.DecortAPICall(ctx, "POST", ComputeUnpinFromStackAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if oldPin.(bool) == false && newPin.(bool) == true {
			_, err := c.DecortAPICall(ctx, "POST", ComputePinToStackAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("pause") {
		oldPause, newPause := d.GetChange("pause")
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		if oldPause.(bool) == true && newPause.(bool) == false {
			_, err := c.DecortAPICall(ctx, "POST", ComputeResumeAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if oldPause.(bool) == false && newPause.(bool) == true {
			_, err := c.DecortAPICall(ctx, "POST", ComputePauseAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("reset") {
		oldReset, newReset := d.GetChange("reset")
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		if oldReset.(bool) == false && newReset.(bool) == true {
			_, err := c.DecortAPICall(ctx, "POST", ComputeResetAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	//redeploy
	if d.HasChange("image_id") {
		oldImage, newImage := d.GetChange("image_id")
		urlValues := &url.Values{}
		urlValues.Add("computeId", d.Id())
		urlValues.Add("force", "false")
		_, err := c.DecortAPICall(ctx, "POST", ComputeStopAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		if oldImage.(int) != newImage.(int) {
			urlValues := &url.Values{}

			urlValues.Add("computeId", d.Id())
			urlValues.Add("imageId", strconv.Itoa(newImage.(int)))
			if diskSize, ok := d.GetOk("boot_disk_size"); ok {
				urlValues.Add("diskSize", strconv.Itoa(diskSize.(int)))
			}
			if dataDisks, ok := d.GetOk("data_disks"); ok {
				urlValues.Add("dataDisks", dataDisks.(string))
			}
			if autoStart, ok := d.GetOk("auto_start"); ok {
				urlValues.Add("autoStart", strconv.FormatBool(autoStart.(bool)))
			}
			if forceStop, ok := d.GetOk("force_stop"); ok {
				urlValues.Add("forceStop", strconv.FormatBool(forceStop.(bool)))
			}
			_, err := c.DecortAPICall(ctx, "POST", ComputeRedeployAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// we may reuse dataSourceComputeRead here as we maintain similarity
	// between Compute resource and Compute data source schemas
	return resourceComputeRead(ctx, d, m)
}

func isChangeDisk(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["disk_id"].(int) == elConv["disk_id"].(int) &&
			elOldConv["size"].(int) != elConv["size"].(int) {
			return true
		}
	}
	return false
}

func isContainsDisk(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["disk_name"].(string) == elConv["disk_name"].(string) {
			return true
		}
	}
	return false
}

func isContainsAR(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["key"].(string) == elConv["key"].(string) &&
			elOldConv["value"].(string) == elConv["value"].(string) &&
			elOldConv["mode"].(string) == elConv["mode"].(string) &&
			elOldConv["topology"].(string) == elConv["topology"].(string) &&
			elOldConv["policy"].(string) == elConv["policy"].(string) {
			return true
		}
	}
	return false
}

func resourceComputeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOTE: this function destroys target Compute instance "permanently", so
	// there is no way to restore it.
	// If compute being destroyed has some extra disks attached, they are
	// detached from the compute
	log.Debugf("resourceComputeDelete: called for Compute name %s, RG ID %d",
		d.Get("name").(string), d.Get("rg_id").(int))

	c := m.(*controller.ControllerCfg)

	params := &url.Values{}
	params.Add("computeId", d.Id())
	params.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))
	params.Add("detachDisks", strconv.FormatBool(d.Get("detach_disks").(bool)))

	if _, err := c.DecortAPICall(ctx, "POST", ComputeDeleteAPI, params); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func disksSubresourceSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"disk_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name for disk",
		},
		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Disk size in GiB",
		},
		"disk_type": {
			Type:         schema.TypeString,
			Computed:     true,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"B", "D"}, false),
			Description:  "The type of disk in terms of its role in compute: 'B=Boot, D=Data'",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Storage endpoint provider ID; by default the same with boot disk",
		},
		"pool": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Pool name; by default will be chosen automatically",
		},
		"desc": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Optional description",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Specify image id for create disk from template",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Description: "Disk deletion status",
		},
		"disk_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk ID",
		},
		"shareable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"size_max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"size_used": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
	return rets
}

func tagsSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func portForwardingSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"public_port_start": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"public_port_end": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  -1,
		},
		"local_port": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"proto": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
		},
	}
}

func userAccessSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:     schema.TypeString,
			Required: true,
		},
		"access_type": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func snapshotSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func snapshotRollbackSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func cdSubresourceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cdrom_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
	}
}

func ResourceComputeSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of this compute. Compute names are case sensitive and must be unique in the resource group.",
		},

		"rg_id": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "ID of the resource group where this compute should be deployed.",
		},

		"driver": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			StateFunc:    statefuncs.StateFuncToUpper,
			ValidateFunc: validation.StringInSlice([]string{"KVM_X86", "KVM_PPC"}, false), // observe case while validating
			Description:  "Hardware architecture of this compute instance.",
		},

		"cpu": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, constants.MaxCpusPerCompute),
			Description:  "Number of CPUs to allocate to this compute instance.",
		},

		"ram": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(constants.MinRamPerCompute),
			Description:  "Amount of RAM in MB to allocate to this compute instance.",
		},

		"image_id": {
			Type:     schema.TypeInt,
			Required: true,
			//ForceNew:    true, //REDEPLOY
			Description: "ID of the OS image to base this compute instance on.",
		},

		"boot_disk_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "This compute instance boot disk size in GB. Make sure it is large enough to accomodate selected OS image.",
		},

		"affinity_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set affinity label for compute",
		},

		"affinity_rules": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"topology": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"node", "compute"}, false),
						Description:  "compute or node, for whom rule applies",
					},
					"policy": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"RECOMMENDED", "REQUIRED"}, false),
						Description:  "RECOMMENDED or REQUIRED, the degree of 'strictness' of this rule",
					},
					"mode": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"EQ", "NE", "ANY"}, false),
						Description:  "EQ or NE or ANY - the comparison mode is 'value', recorded by the specified 'key'",
					},
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "key that are taken into account when analyzing this rule will be identified",
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "value that must match the key to be taken into account when analyzing this rule",
					},
				},
			},
		},

		"anti_affinity_rules": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"topology": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"node", "compute"}, false),
						Description:  "compute or node, for whom rule applies",
					},
					"policy": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"RECOMMENDED", "REQUIRED"}, false),
						Description:  "RECOMMENDED or REQUIRED, the degree of 'strictness' of this rule",
					},
					"mode": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"EQ", "NE", "ANY"}, false),
						Description:  "EQ or NE or ANY - the comparison mode is 'value', recorded by the specified 'key'",
					},
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "key that are taken into account when analyzing this rule will be identified",
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "value that must match the key to be taken into account when analyzing this rule",
					},
				},
			},
		},

		"disks": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: disksSubresourceSchemaMake(),
			},
		},

		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "ID of SEP to create bootDisk on. Uses image's sepId if not set.",
		},

		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "Pool to use if sepId is set, can be also empty if needed to be chosen by system.",
		},

		"extra_disks": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: constants.MaxExtraDisksPerCompute,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "Optional list of IDs of extra disks to attach to this compute. You may specify several extra disks.",
		},

		"network": {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			MaxItems: constants.MaxNetworksPerCompute,
			Elem: &schema.Resource{
				Schema: networkSubresourceSchemaMake(),
			},
			Description: "Optional network connection(s) for this compute. You may specify several network blocks, one for each connection.",
		},

		/*
			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: MaxSshKeysPerCompute,
				Elem: &schema.Resource{
					Schema: sshSubresourceSchemaMake(),
				},
				Description: "SSH keys to authorize on this compute instance.",
			},
		*/

		"tags": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: tagsSubresourceSchemaMake(),
			},
		},

		"port_forwarding": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: portForwardingSubresourceSchemaMake(),
			},
		},

		"user_access": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: userAccessSubresourceSchemaMake(),
			},
		},

		"snapshot": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: snapshotSubresourceSchemaMake(),
			},
		},

		"rollback": {
			Type:     schema.TypeSet,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: snapshotRollbackSubresourceSchemaMake(),
			},
		},

		"cd": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: cdSubresourceSchemaMake(),
			},
		},

		"pin_to_stack": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional text description of this compute instance.",
		},

		"cloud_init": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional cloud_init parameters. Applied when creating new compute instance only, ignored in all other cases.",
		},

		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "If true - enable compute, else - disable",
		},

		"pause": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"reset": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"auto_start": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Flag for redeploy compute",
		},
		"force_stop": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Flag for redeploy compute",
		},
		"data_disks": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"KEEP", "DETACH", "DESTROY"}, false),
			Default:      "DETACH",
			Description:  "Flag for redeploy compute",
		},
		"started": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Is compute started.",
		},
		"detach_disks": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"permanently": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"is": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "system name",
		},
		"ipa_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "compute purpose",
		},

		// The rest are Compute properties, which are "computed" once it is created
		"account_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the account this compute instance belongs to.",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the account this compute instance belongs to.",
		},
		"affinity_weight": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"arch": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"boot_order": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"boot_disk_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "This compute instance boot disk ID.",
		},
		"clone_reference": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"clones": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"computeci_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"custom_fields": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"val": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"deleted_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"devices": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"compute_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"interfaces": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeInterfacesSchemaMake(),
			},
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manager_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"manager_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"migrationjob": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"natable_vins_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"natable_vins_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_network_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"os_users": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: osUsersSubresourceSchemaMake(),
			},
			Description: "Guest OS users provisioned on this compute instance.",
		},
		"pinned": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"registered": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"res_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rg_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the resource group where this compute instance is located.",
		},
		"snap_sets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computeSnapSetsSchemaMake(),
			},
		},
		"stateless_sep_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"stateless_sep_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"user_managed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"vgpus": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"virtual_image_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"virtual_image_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	return rets
}

func ResourceCompute() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceComputeCreate,
		ReadContext:   resourceComputeRead,
		UpdateContext: resourceComputeUpdate,
		DeleteContext: resourceComputeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &constants.Timeout600s,
			Read:    &constants.Timeout300s,
			Update:  &constants.Timeout300s,
			Delete:  &constants.Timeout300s,
			Default: &constants.Timeout300s,
		},

		Schema: ResourceComputeSchemaMake(),
	}
}
