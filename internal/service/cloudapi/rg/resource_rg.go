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

package rg

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/constants"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/controller"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/dc"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/location"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/status"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceResgroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// First validate that we have all parameters required to create the new Resource Group

	// Valid account ID is required to create new resource group
	// obtain Account ID by account name - it should not be zero on success

	rgName, argSet := d.GetOk("name")
	if !argSet {
		return diag.FromErr(fmt.Errorf("Cannot create new RG: missing name."))
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	/* Current version of provider works with default grid id (same is true for disk resources)
	grid_id, arg_set := d.GetOk("grid_id")
	if !arg_set {
		return fmt.Errorf("Cannot create new RG %q in account ID %d: missing Grid ID.",
			rg_name.(string), validated_account_id)
	}
	if grid_id.(int) < 1 {
		grid_id = DefaultGridID
	}
	*/

	// all required parameters are set in the schema - we can continue with RG creation
	log.Debugf("resourceResgroupCreate: called for RG name %s, account ID %d",
		rgName.(string), d.Get("account_id").(int))

	// Check input values
	// AccountID
	haveAccount, err := existAccountID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveAccount {
		return diag.Errorf("resourceResgroupCreate: can't create RG bacause AccountID %d not allowed or does not exist", d.Get("account_id").(int))
	}
	// GID
	haveGID, err := existGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveGID {
		return diag.Errorf("resourceResgroupCreate: can't create RG bacause GID %d not allowed or does not exist", d.Get("gid").(int))
	}
	// ExtNetID
	if _, ok := d.GetOk("ext_net_id"); ok {
		haveExtNet, err := existExtNetID(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
		if !haveExtNet {
			return diag.Errorf("resourceResgroupCreate: can't create RG bacause ExtNetID %d not allowed or does not exist", d.Get("ext_net_id").(int))
		}
	}

	// quota settings are optional
	setQuota := false
	var quotaRecord QuotaRecord
	argValue, argSet := d.GetOk("quota")
	if argSet {
		log.Debugf("resourceResgroupCreate: setting Quota on RG requested")
		quotaRecord = makeQuotaRecord(argValue.([]interface{}))
		setQuota = true
	}

	log.Debugf("resourceResgroupCreate: called by user %q for RG name %s, account ID %d",
		c.GetDecortUsername(),
		rgName.(string), d.Get("account_id").(int))

	urlValues = &url.Values{}
	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("name", rgName.(string))
	urlValues.Add("gid", strconv.Itoa(location.DefaultGridID)) // use default Grid ID, similar to disk resource mgmt convention
	urlValues.Add("owner", c.GetDecortUsername())

	// pass quota values as set
	if setQuota {
		urlValues.Add("maxCPUCapacity", strconv.Itoa(quotaRecord.Cpu))
		urlValues.Add("maxVDiskCapacity", strconv.Itoa(quotaRecord.Disk))
		urlValues.Add("maxMemoryCapacity", fmt.Sprintf("%f", quotaRecord.Ram)) // RAM quota is float; this may change in the future
		urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(quotaRecord.ExtTraffic))
		urlValues.Add("maxNumPublicIP", strconv.Itoa(quotaRecord.ExtIPs))
	}

	// parse and handle network settings
	defNetType, argSet := d.GetOk("def_net_type")
	if argSet {
		urlValues.Add("def_net", defNetType.(string)) // NOTE: in API default network type is set by "def_net" parameter
	} else {
		d.Set("def_net_type", "PRIVATE")
	}

	ipcidr, argSet := d.GetOk("ipcidr")
	if argSet {
		urlValues.Add("ipcidr", ipcidr.(string))
	}

	description, argSet := d.GetOk("description")
	if argSet {
		urlValues.Add("desc", description.(string))
	}

	reason, argSet := d.GetOk("reason")
	if argSet {
		urlValues.Add("reason", reason.(string))
	}

	extNetId, argSet := d.GetOk("ext_net_id")
	if argSet {
		urlValues.Add("extNetId", strconv.Itoa(extNetId.(int)))
	}

	extIp, argSet := d.GetOk("ext_ip")
	if argSet {
		urlValues.Add("extIp", extIp.(string))
	}

	apiResp, err := c.DecortAPICall(ctx, "POST", ResgroupCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(apiResp) // rg/create API returns ID of the newly creted resource group on success

	w := dc.Warnings{}
	if access, ok := d.GetOk("access"); ok {
		urlValues = &url.Values{}
		var user, right string

		if access.(*schema.Set).Len() > 0 {
			accessList := access.(*schema.Set).List()
			for _, accessInterface := range accessList {
				access := accessInterface.(map[string]interface{})
				user = access["user"].(string)
				right = access["right"].(string)

				urlValues.Add("rgId", d.Id())
				urlValues.Add("user", user)
				urlValues.Add("right", right)
				if reason, ok := d.GetOk("reason"); ok {
					urlValues.Add("reason", reason.(string))
				}

				_, err := c.DecortAPICall(ctx, "POST", RgAccessGrantAPI, urlValues)
				if err != nil {
					w.Add(err)
				}
			}
		}

	}

	if defNet, ok := d.GetOk("def_net"); ok {
		urlValues := &url.Values{}

		if defNet.(*schema.Set).Len() > 0 {
			defNetList := defNet.(*schema.Set).List()
			defNetItem := defNetList[0].(map[string]interface{})

			netType := defNetItem["net_type"].(string)

			urlValues.Add("rgId", d.Id())
			urlValues.Add("netType", netType)

			if netID, ok := defNetItem["net_id"]; ok {
				urlValues.Add("netId", strconv.Itoa(netID.(int)))
			}
			if reason, ok := defNetItem["reason"]; ok {
				urlValues.Add("reason", reason.(string))
			}

			_, err := c.DecortAPICall(ctx, "POST", RgSetDefNetAPI, urlValues)
			if err != nil {
				w.Add(err)
			}
			d.Set("def_net_type", netType)
		}

	}

	if enable, ok := d.GetOk("enable"); ok {
		urlValues = &url.Values{}

		api := RgDisableAPI
		enable := enable.(bool)
		if enable {
			api = RgEnableAPI
		}
		urlValues.Add("rgId", d.Id())

		if reason, ok := d.GetOk("reason"); ok {
			urlValues.Add("reason", reason.(string))
		}

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			w.Add(err)
		}
	}

	// re-read newly created RG to make sure schema contains complete and up to date set of specifications
	defer resourceResgroupRead(ctx, d, m)
	return w.Get()
}

func resourceResgroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceResgroupRead: called for RG name %s, account ID %d",
		d.Get("name").(string), d.Get("account_id").(int))

	c := m.(*controller.ControllerCfg)

	rgFacts, err := utilityResgroupCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("") // ensure ID is empty
		return diag.FromErr(err)
	}

	switch rgFacts.Status {
	case status.Modeled:
	case status.Created:
	case status.Enabled:
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("rgId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", RgRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	case status.Deleting:
	case status.Destroyed:
		d.SetId("")
		return resourceResgroupCreate(ctx, d, m)
	case status.Destroying:
	case status.Disabled:
	case status.Disabling:
	case status.Enabled:
	case status.Enabling:
	}

	rgFacts, err = utilityResgroupCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("") // ensure ID is empty
		return diag.FromErr(err)
	}
	return diag.FromErr(flattenResgroup(d, *rgFacts))
}

func resourceResgroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceResgroupUpdate: called for RG name %s, account ID %d",
		d.Get("name").(string), d.Get("account_id").(int))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	// Check input values
	// AccountID
	haveAccount, err := existAccountID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveAccount {
		return diag.Errorf("resourceResgroupUpdate: can't create RG bacause AccountID %d not allowed or does not exist", d.Get("account_id").(int))
	}
	// GID
	haveGID, err := existGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveGID {
		return diag.Errorf("resourceResgroupUpdate: can't create RG bacause GID %d not allowed or does not exist", d.Get("gid").(int))
	}
	// ExtNetID
	if _, ok := d.GetOk("ext_net_id"); ok {
		haveExtNet, err := existExtNetID(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
		if !haveExtNet {
			return diag.Errorf("resourceResgroupUpdate: can't create RG bacause ExtNetID %d not allowed or does not exist", d.Get("ext_net_id").(int))
		}
	}

	rgFacts, err := utilityResgroupCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("") // ensure ID is empty
		return diag.FromErr(err)
	}

	switch rgFacts.Status {
	case status.Modeled:
	case status.Created:
	case status.Enabled:
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("rgId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", RgRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	case status.Deleting:
	case status.Destroyed:
		d.SetId("")
		return resourceResgroupCreate(ctx, d, m)
	case status.Destroying:
	case status.Disabled:
	case status.Disabling:
	case status.Enabled:
	case status.Enabling:
	}
	/* NOTE: we do not allow changing the following attributes of an existing RG via terraform:
	   - def_net_type
	   - ipcidr
	   - ext_net_id
	   - ext_ip

	   The following code fragment checks if any of these have been changed and generates error.
	*/
	if ok := d.HasChange("def_net"); ok {
		_, newDefNet := d.GetChange("def_net")
		if newDefNet.(*schema.Set).Len() == 0 {
			return diag.Errorf("resourceResgroupUpdate: block def_net must not be empty")
		}
	}

	for _, attr := range []string{"def_net_type", "ipcidr", "ext_ip"} {
		attr_new, attr_old := d.GetChange("def_net_type")
		if attr_new.(string) != attr_old.(string) {
			return diag.FromErr(fmt.Errorf("resourceResgroupUpdate: RG ID %s: changing %s for existing RG is not allowed", d.Id(), attr))
		}
	}

	attrNew, attrOld := d.GetChange("ext_net_id")
	if attrNew.(int) != attrOld.(int) {
		return diag.FromErr(fmt.Errorf("resourceResgroupUpdate: RG ID %s: changing ext_net_id for existing RG is not allowed", d.Id()))
	}

	doGeneralUpdate := false // will be true if general RG update is necessary (API rg/update)

	urlValues = &url.Values{}
	urlValues.Add("rgId", d.Id())

	nameNew, nameSet := d.GetOk("name")
	if nameSet {
		log.Debugf("resourceResgroupUpdate: name specified - looking for deltas from the old settings.")
		nameOld, _ := d.GetChange("name")
		if nameOld.(string) != nameNew.(string) {
			doGeneralUpdate = true
			urlValues.Add("name", nameNew.(string))
		}
	}

	quotaValue, quotaSet := d.GetOk("quota")
	if quotaSet {
		log.Debugf("resourceResgroupUpdate: quota specified - looking for deltas from the old quota.")
		quotarecordNew := makeQuotaRecord(quotaValue.([]interface{}))
		quotaValueOld, _ := d.GetChange("quota") // returns old as 1st, new as 2nd return value
		quotarecordOld := makeQuotaRecord(quotaValueOld.([]interface{}))
		log.Debug(quotaValueOld, quotarecordNew)

		if quotarecordNew.Cpu != quotarecordOld.Cpu {
			doGeneralUpdate = true
			log.Debugf("resourceResgroupUpdate: Cpu diff %d <- %d", quotarecordNew.Cpu, quotarecordOld.Cpu)
			urlValues.Add("maxCPUCapacity", strconv.Itoa(quotarecordNew.Cpu))
		}

		if quotarecordNew.Disk != quotarecordOld.Disk {
			doGeneralUpdate = true
			log.Debugf("resourceResgroupUpdate: Disk diff %d <- %d", quotarecordNew.Disk, quotarecordOld.Disk)
			urlValues.Add("maxVDiskCapacity", strconv.Itoa(quotarecordNew.Disk))
		}

		if quotarecordNew.Ram != quotarecordOld.Ram { // NB: quota on RAM is stored as float32, in units of MB
			doGeneralUpdate = true
			log.Debugf("resourceResgroupUpdate: Ram diff %f <- %f", quotarecordNew.Ram, quotarecordOld.Ram)
			urlValues.Add("maxMemoryCapacity", fmt.Sprintf("%f", quotarecordNew.Ram))
		}

		if quotarecordNew.ExtTraffic != quotarecordOld.ExtTraffic {
			doGeneralUpdate = true
			log.Debugf("resourceResgroupUpdate: ExtTraffic diff %d <- %d", quotarecordNew.ExtTraffic, quotarecordOld.ExtTraffic)
			urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(quotarecordNew.ExtTraffic))
		}

		if quotarecordNew.ExtIPs != quotarecordOld.ExtIPs {
			doGeneralUpdate = true
			log.Debugf("resourceResgroupUpdate: ExtIPs diff %d <- %d", quotarecordNew.ExtIPs, quotarecordOld.ExtIPs)
			urlValues.Add("maxNumPublicIP", strconv.Itoa(quotarecordNew.ExtIPs))
		}
	} else {
		doGeneralUpdate = true
		urlValues.Add("maxCPUCapacity", "-1")
		urlValues.Add("maxVDiskCapacity", "-1")
		urlValues.Add("maxMemoryCapacity", "-1")
		urlValues.Add("maxNetworkPeerTransfer", "-1")
		urlValues.Add("maxNumPublicIP", "-1")
	}

	descNew, descSet := d.GetOk("description")
	if descSet {
		log.Debugf("resourceResgroupUpdate: description specified - looking for deltas from the old settings.")
		descOld, _ := d.GetChange("description")
		if descOld.(string) != descNew.(string) {
			doGeneralUpdate = true
			urlValues.Add("desc", descNew.(string))
		}
	}

	if doGeneralUpdate {
		log.Debugf("resourceResgroupUpdate: detected delta between new and old RG specs - updating the RG")
		_, err := c.DecortAPICall(ctx, "POST", ResgroupUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		log.Debugf("resourceResgroupUpdate: no difference between old and new state - no update on the RG will be done")
	}

	urlValues = &url.Values{}
	enableOld, enableNew := d.GetChange("enable")
	if enableOld.(bool) && !enableNew.(bool) {
		urlValues.Add("rgId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", RgDisableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if !enableOld.(bool) && enableNew.(bool) {
		urlValues.Add("rgId", d.Id())
		_, err := c.DecortAPICall(ctx, "POST", RgEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	urlValues = &url.Values{}

	oldSet, newSet := d.GetChange("access")

	deletedAccess := (oldSet.(*schema.Set).Difference(newSet.(*schema.Set))).List()
	for _, deletedInterface := range deletedAccess {
		deletedItem := deletedInterface.(map[string]interface{})

		user := deletedItem["user"].(string)

		urlValues.Add("rgId", d.Id())
		urlValues.Add("user", user)
		if reason, ok := d.GetOk("reason"); ok {
			urlValues.Add("reason", reason.(string))
		}

		_, err := c.DecortAPICall(ctx, "POST", RgAccessRevokeAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		urlValues = &url.Values{}
	}

	addedAccess := (newSet.(*schema.Set).Difference(oldSet.(*schema.Set))).List()
	for _, addedInterface := range addedAccess {
		addedItem := addedInterface.(map[string]interface{})

		user := addedItem["user"].(string)
		right := addedItem["right"].(string)

		urlValues.Add("rgId", d.Id())
		urlValues.Add("user", user)
		urlValues.Add("right", right)
		if reason, ok := d.GetOk("reason"); ok {
			urlValues.Add("reason", reason.(string))
		}

		_, err := c.DecortAPICall(ctx, "POST", RgAccessGrantAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		urlValues = &url.Values{}
	}

	if ok := d.HasChange("def_net"); ok {
		oldDefNet, newDefNet := d.GetChange("def_net")
		if newDefNet.(*schema.Set).Len() > 0 {
			changedDefNet := (newDefNet.(*schema.Set).Difference(oldDefNet.(*schema.Set))).List()
			for _, changedDefNetInterface := range changedDefNet {

				defNetItem := changedDefNetInterface.(map[string]interface{})

				netType := defNetItem["net_type"].(string)

				urlValues.Add("rgId", d.Id())
				urlValues.Add("netType", netType)

				if netID, ok := defNetItem["net_id"]; ok {
					urlValues.Add("netId", strconv.Itoa(netID.(int)))
				}
				if reason, ok := defNetItem["reason"]; ok {
					urlValues.Add("reason", reason.(string))
				}

				_, err := c.DecortAPICall(ctx, "POST", RgSetDefNetAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceResgroupRead(ctx, d, m)
}

func resourceResgroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOTE: this method forcibly destroys target resource group with flag "permanently", so there is no way to
	// restore the destroyed resource group as well all Computes & VINSes that existed in it
	log.Debugf("resourceResgroupDelete: called for RG name %s, account ID %d",
		d.Get("name").(string), d.Get("account_id").(int))

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("rgId", d.Id())
	if force, ok := d.GetOk("force"); ok {
		urlValues.Add("force", strconv.FormatBool(force.(bool)))
	}
	if permanently, ok := d.GetOk("permanently"); ok {
		urlValues.Add("permanently", strconv.FormatBool(permanently.(bool)))
	}
	if reason, ok := d.GetOk("reason"); ok {
		urlValues.Add("reason", reason.(string))
	}

	_, err := c.DecortAPICall(ctx, "POST", ResgroupDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceRgSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Unique ID of the account, which this resource group belongs to.",
		},

		"gid": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true, // change of Grid ID will require new RG
			Description: "Unique ID of the grid, where this resource group is deployed.",
		},

		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of this resource group. Names are case sensitive and unique within the context of a account.",
		},

		"def_net_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC", "NONE"}, false),
			Description:  "Type of the network, which this resource group will use as default for its computes - PRIVATE or PUBLIC or NONE.",
		},

		"def_net_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the default network for this resource group (if any).",
		},

		"ipcidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Address of the netowrk inside the private network segment (aka ViNS) if def_net_type=PRIVATE",
		},

		"ext_net_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "ID of the external network for default ViNS. Pass 0 if def_net_type=PUBLIC or no external connection required for the defult ViNS when def_net_type=PRIVATE",
		},

		"ext_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IP address on the external netowrk to request when def_net_type=PRIVATE and ext_net_id is not 0",
		},

		"quota": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: quotaRgSubresourceSchemaMake(),
			},
			Description: "Quota settings for this resource group.",
		},

		"access": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "User or group name to grant access",
					},
					"right": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Access rights to set, one of 'R', 'RCX' or 'ARCXDU'",
					},
					"reason": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Reason for action",
					},
				},
			},
		},

		"def_net": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"net_type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
						Description:  "Network type to set. Must be on of 'PRIVATE' or 'PUBLIC'.",
					},
					"net_id": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     0,
						Description: "Network segment ID. If netType is PUBLIC and netId is 0 then default external network segment will be selected. If netType is PRIVATE and netId=0, the first ViNS defined for this RG will be selected. Otherwise, netId identifies either existing external network segment or ViNS.",
					},
					"reason": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Reason for action",
					},
				},
			},
		},

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "User-defined text description of this resource group.",
		},
		"force": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set to True if you want force delete non-empty RG",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set to True if you want force delete non-empty RG",
		},
		"reason": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set to True if you want force delete non-empty RG",
		},
		"register_computes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Register computes in registration system",
		},

		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "flag for enable/disable RG",
		},

		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the account, which this resource group belongs to.",
		},
		"resources": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: resourcesSchemaMake(),
			},
		},

		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclSchemaMake(),
			},
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deleted_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dirty": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lock_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"secret": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current status of this resource group.",
		},
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vins": {
			Type:     schema.TypeList, //this is a list of ints
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "List of VINs deployed in this resource group.",
		},

		"vms": {
			Type:     schema.TypeList, //t his is a list of ints
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Description: "List of computes deployed in this resource group.",
		},

		"res_types": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"uniq_pools": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func ResourceResgroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceResgroupCreate,
		ReadContext:   resourceResgroupRead,
		UpdateContext: resourceResgroupUpdate,
		DeleteContext: resourceResgroupDelete,

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

		Schema: ResourceRgSchemaMake(),
		CustomizeDiff: customdiff.All(
			customdiff.IfValueChange("def_net",
				func(ctx context.Context, oldValue, newValue, meta interface{}) bool {
					return true
				},
				func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
					oldValue, newValue := d.GetChange("def_net")

					old := len(oldValue.(*schema.Set).List())
					new_ := len(newValue.(*schema.Set).List())

					if old == 1 && new_ == 0 {
						return fmt.Errorf("CustomizeDiff: block def_net must not be empty")
					}
					return nil
				},
			),
		),
	}
}
