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

package account

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/status"
)

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceAccountCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("name", d.Get("account_name").(string))
	urlValues.Add("username", d.Get("username").(string))

	if emailaddress, ok := d.GetOk("emailaddress"); ok {
		urlValues.Add("emailaddress", emailaddress.(string))
	}
	if sendAccessEmails, ok := d.GetOk("send_access_emails"); ok {
		urlValues.Add("sendAccessEmails", strconv.FormatBool(sendAccessEmails.(bool)))
	}
	if resLimits, ok := d.GetOk("resource_limits"); ok {
		resLimit := resLimits.([]interface{})[0]
		resLimitConv := resLimit.(map[string]interface{})
		if resLimitConv["cu_m"] != nil {
			maxMemCap := int(resLimitConv["cu_m"].(float64))
			if maxMemCap == 0 {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(maxMemCap))
			}
		}
		if resLimitConv["cu_d"] != nil {
			maxDiskCap := int(resLimitConv["cu_d"].(float64))
			if maxDiskCap == 0 {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(maxDiskCap))
			}
		}
		if resLimitConv["cu_c"] != nil {
			maxCPUCap := int(resLimitConv["cu_c"].(float64))
			if maxCPUCap == 0 {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(maxCPUCap))
			}

		}
		if resLimitConv["cu_i"] != nil {
			maxNumPublicIP := int(resLimitConv["cu_i"].(float64))
			if maxNumPublicIP == 0 {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(maxNumPublicIP))
			}

		}
		if resLimitConv["cu_np"] != nil {
			maxNP := int(resLimitConv["cu_np"].(float64))
			if maxNP == 0 {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(maxNP))
			}

		}
		if resLimitConv["gpu_units"] != nil {
			gpuUnits := int(resLimitConv["gpu_units"].(float64))
			if gpuUnits == 0 {
				urlValues.Add("gpu_units", strconv.Itoa(-1))
			} else {
				urlValues.Add("gpu_units", strconv.Itoa(gpuUnits))
			}
		}
	}

	accountId, err := c.DecortAPICall(ctx, "POST", accountCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	accIdParsed, _ := strconv.Atoi(accountId)

	d.SetId(accountId)
	d.Set("account_id", accIdParsed)

	urlValues = &url.Values{}

	if enable, ok := d.GetOk("enable"); ok {
		api := accountDisableAPI
		enable := enable.(bool)
		if enable {
			api = accountEnableAPI
		}
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if users, ok := d.GetOk("users"); ok {
		addedUsers := users.([]interface{})

		if len(addedUsers) > 0 {
			for _, user := range addedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("accesstype", strings.ToUpper(userConv["access_type"].(string)))
				_, err := c.DecortAPICall(ctx, "POST", accountAddUserAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}
	}

	return resourceAccountRead(ctx, d, m)
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceAccountRead: called for account with ID: %v", d.Id())

	c := m.(*controller.ControllerCfg)

	acc, err := utilityAccountCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChanged := false

	switch acc.Status {
	case status.Destroyed:
		d.SetId("")
		return resourceAccountCreate(ctx, d, m)
	case status.Destroying:
		return diag.Errorf("The account is in progress with status: %s", acc.Status)
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("accountId", d.Id())

		_, err := c.DecortAPICall(ctx, "POST", accountRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Disabled:
		log.Debugf("The account is in status: %s, troubles may occur with update. Please, enable account first.", acc.Status)
	case status.Confirmed:
	}

	if hasChanged {
		acc, err = utilityAccountCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	flattenAccount(d, *acc)

	return nil
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceAccountDelete")

	account, err := utilityAccountCheckPresence(ctx, d, m)
	if account == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))

	_, err = c.DecortAPICall(ctx, "POST", accountDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceAccountEdit")
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	acc, err := utilityAccountCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChanged := false

	switch acc.Status {
	case status.Destroyed:
		d.SetId("")
		return resourceAccountCreate(ctx, d, m)
	case status.Destroying:
		return diag.Errorf("The account is in progress with status: %s", acc.Status)
	case status.Deleted:
		urlVal := &url.Values{}
		urlVal.Add("accountId", d.Id())

		_, err := c.DecortAPICall(ctx, "POST", accountRestoreAPI, urlVal)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Disabled:
		log.Debugf("The account is in status: %s, troubles may occur with update. Please, enable account first.", acc.Status)
	case status.Confirmed:
	}

	if hasChanged {
		acc, err = utilityAccountCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enable") {
		api := accountDisableAPI
		enable := d.Get("enable").(bool)
		if enable {
			api = accountEnableAPI
		}
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("account_name") {
		urlValues.Add("name", d.Get("account_name").(string))
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
		_, err := c.DecortAPICall(ctx, "POST", accountUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}
	if d.HasChange("resource_limits") {
		resLimit := d.Get("resource_limits").([]interface{})[0]
		resLimitConv := resLimit.(map[string]interface{})

		if resLimitConv["cu_m"] != nil {
			maxMemCap := int(resLimitConv["cu_m"].(float64))
			if maxMemCap == 0 {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxMemoryCapacity", strconv.Itoa(maxMemCap))
			}
		}
		if resLimitConv["cu_d"] != nil {
			maxDiskCap := int(resLimitConv["cu_d"].(float64))
			if maxDiskCap == 0 {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxVDiskCapacity", strconv.Itoa(maxDiskCap))
			}
		}
		if resLimitConv["cu_c"] != nil {
			maxCPUCap := int(resLimitConv["cu_c"].(float64))
			if maxCPUCap == 0 {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxCPUCapacity", strconv.Itoa(maxCPUCap))
			}

		}
		if resLimitConv["cu_i"] != nil {
			maxNumPublicIP := int(resLimitConv["cu_i"].(float64))
			if maxNumPublicIP == 0 {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNumPublicIP", strconv.Itoa(maxNumPublicIP))
			}

		}
		if resLimitConv["cu_np"] != nil {
			maxNP := int(resLimitConv["cu_np"].(float64))
			if maxNP == 0 {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(-1))
			} else {
				urlValues.Add("maxNetworkPeerTransfer", strconv.Itoa(maxNP))
			}

		}
		if resLimitConv["gpu_units"] != nil {
			gpuUnits := int(resLimitConv["gpu_units"].(float64))
			if gpuUnits == 0 {
				urlValues.Add("gpu_units", strconv.Itoa(-1))
			} else {
				urlValues.Add("gpu_units", strconv.Itoa(gpuUnits))
			}
		}

		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
		_, err := c.DecortAPICall(ctx, "POST", accountUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("send_access_emails") {
		urlValues.Add("sendAccessEmails", strconv.FormatBool(d.Get("send_access_emails").(bool)))
		urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
		_, err := c.DecortAPICall(ctx, "POST", accountUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("restore") {
		restore := d.Get("restore").(bool)
		if restore {
			urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
			_, err := c.DecortAPICall(ctx, "POST", accountRestoreAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
		}
	}

	if d.HasChange("users") {
		deletedUsers := make([]interface{}, 0)
		addedUsers := make([]interface{}, 0)
		updatedUsers := make([]interface{}, 0)

		old, new := d.GetChange("users")
		oldConv := old.([]interface{})
		newConv := new.([]interface{})
		for _, el := range oldConv {
			if !isContainsUser(newConv, el) {
				deletedUsers = append(deletedUsers, el)
			}
		}
		for _, el := range newConv {
			if !isContainsUser(oldConv, el) {
				addedUsers = append(addedUsers, el)
			} else {
				if isChangedUser(oldConv, el) {
					updatedUsers = append(updatedUsers, el)
				}
			}
		}

		if len(deletedUsers) > 0 {
			for _, user := range deletedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("recursivedelete", strconv.FormatBool(userConv["recursive_delete"].(bool)))
				_, err := c.DecortAPICall(ctx, "POST", accountDeleteUserAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

		if len(addedUsers) > 0 {
			for _, user := range addedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("accesstype", strings.ToUpper(userConv["access_type"].(string)))
				_, err := c.DecortAPICall(ctx, "POST", accountAddUserAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

		if len(updatedUsers) > 0 {
			for _, user := range updatedUsers {
				userConv := user.(map[string]interface{})
				urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
				urlValues.Add("userId", userConv["user_id"].(string))
				urlValues.Add("accesstype", strings.ToUpper(userConv["access_type"].(string)))
				_, err := c.DecortAPICall(ctx, "POST", accountUpdateUserAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

	}

	return resourceAccountRead(ctx, d, m)
}

func isContainsUser(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["user_id"].(string) == elConv["user_id"].(string) {
			return true
		}
	}
	return false
}

func isChangedUser(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["user_id"].(string) == elConv["user_id"].(string) &&
			(!strings.EqualFold(elOldConv["access_type"].(string), elConv["access_type"].(string)) ||
				elOldConv["recursive_delete"].(bool) != elConv["recursive_delete"].(bool)) {
			return true
		}
	}
	return false
}

func resourceAccountSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "account name",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "username of owner the account",
		},
		"emailaddress": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "email",
		},
		"send_access_emails": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "if true send emails when a user is granted access to resources",
		},
		"users": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"access_type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"recursive_delete": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"restore": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "restore a deleted account",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "whether to completely delete the account",
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "enable/disable account",
		},
		"resource_limits": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cu_c": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_d": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_i": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_m": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"cu_np": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
					"gpu_units": {
						Type:     schema.TypeFloat,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dc_location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"resources": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: resourcesSchemaMake(),
			},
		},
		"ckey": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclSchemaMake(),
			},
		},
		"company": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"companyurl": {
			Type:     schema.TypeString,
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
		"deactivation_time": {
			Type:     schema.TypeFloat,
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
		"displayname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"service_account": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vins": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: computesSchemaMake(),
			},
		},
		"machines": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: machinesSchemaMake(),
			},
		},
		"vinses": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func ResourceAccount() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceAccountCreate,
		ReadContext:   resourceAccountRead,
		UpdateContext: resourceAccountUpdate,
		DeleteContext: resourceAccountDelete,

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

		Schema: resourceAccountSchemaMake(),
	}
}
