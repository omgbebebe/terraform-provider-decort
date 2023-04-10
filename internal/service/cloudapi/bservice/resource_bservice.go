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

package bservice

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/status"
)

func resourceBasicServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceCreate")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	haveRGID, err := existRGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveRGID {
		return diag.Errorf("resourceBasicServiceCreate: can't create basic service because RGID %d is not allowed or does not exist", d.Get("rg_id").(int))
	}

	urlValues.Add("name", d.Get("service_name").(string))
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))

	if sshKey, ok := d.GetOk("ssh_key"); ok {
		urlValues.Add("sshKey", sshKey.(string))
	}
	if sshUser, ok := d.GetOk("ssh_user"); ok {
		urlValues.Add("sshUser", sshUser.(string))
	}

	serviceId, err := c.DecortAPICall(ctx, "POST", bserviceCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceIdParsed, _ := strconv.Atoi(serviceId)

	d.SetId(serviceId)
	d.Set("service_id", serviceIdParsed)

	diagnostics := resourceBasicServiceRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceBasicServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceRead")

	c := m.(*controller.ControllerCfg)

	bs, err := utilityBasicServiceCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChanged := false

	switch bs.Status {
	case status.Modeled:
		return diag.Errorf("The basic service is in status: %s, please, contact support for more information", bs.Status)
	case status.Created:
	case status.Enabled:
	case status.Enabling:
	case status.Disabled:
		log.Debugf("The basic service is in status: %s, troubles can occur with the update. Please, enable bservice first.", bs.Status)
	case status.Disabling:
		log.Debugf("The basic service is in status: %s, troubles can occur with the update.", bs.Status)
	case status.Deleted:
		urlVal := &url.Values{}
		urlVal.Add("serviceId", d.Id())

		_, err := c.DecortAPICall(ctx, "POST", bserviceRestoreAPI, urlVal)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Deleting:
	case status.Destroyed:
		d.SetId("")
		return resourceBasicServiceCreate(ctx, d, m)
	case status.Destroying:
		return diag.Errorf("The basic service is in progress with status: %s", bs.Status)
	case status.Restoring:
	case status.Reconfiguring:
	}

	if hasChanged {
		bs, err = utilityBasicServiceCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	flattenService(d, bs)

	return nil
}

func resourceBasicServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceDelete")

	bs, err := utilityBasicServiceCheckPresence(ctx, d, m)
	if bs == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
	urlValues.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))

	_, err = c.DecortAPICall(ctx, "POST", bserviceDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceBasicServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceBasicServiceEdit")
	c := m.(*controller.ControllerCfg)

	haveRGID, err := existRGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveRGID {
		return diag.Errorf("resourceBasicServiceUpdate: can't create basic service because RGID %d is not allowed or does not exist", d.Get("rg_id").(int))
	}

	bs, err := utilityBasicServiceCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChanged := false

	switch bs.Status {
	case status.Modeled:
		return diag.Errorf("The basic service is in status: %s, please, contact support for more information", bs.Status)
	case status.Created:
	case status.Enabled:
	case status.Enabling:
	case status.Disabled:
		log.Debugf("The basic service is in status: %s, troubles can occur with the update. Please, enable bservice first.", bs.Status)
	case status.Disabling:
		log.Debugf("The basic service is in status: %s, troubles can occur with the update.", bs.Status)
	case status.Deleted:
		urlVal := &url.Values{}
		urlVal.Add("serviceId", d.Id())

		_, err := c.DecortAPICall(ctx, "POST", bserviceRestoreAPI, urlVal)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Deleting:
	case status.Destroyed:
		d.SetId("")
		return resourceBasicServiceCreate(ctx, d, m)
	case status.Destroying:
		return diag.Errorf("The basic service is in progress with status: %s", bs.Status)
	case status.Restoring:
	case status.Reconfiguring:
	}

	if hasChanged {
		bs, err = utilityBasicServiceCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	urlValues := &url.Values{}
	if d.HasChange("enable") {
		api := bserviceDisableAPI
		enable := d.Get("enable").(bool)
		if enable {
			api = bserviceEnableAPI
		}
		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("restore") {
		restore := d.Get("restore").(bool)
		if restore {
			urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
			_, err := c.DecortAPICall(ctx, "POST", bserviceRestoreAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
		}
	}

	if d.HasChange("start") {
		api := bserviceStopAPI
		start := d.Get("start").(bool)
		if start {
			api = bserviceStartAPI
		}
		urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("snapshots") {
		deletedSnapshots := make([]interface{}, 0)
		addedSnapshots := make([]interface{}, 0)
		updatedSnapshots := make([]interface{}, 0)

		old, new := d.GetChange("snapshots")
		oldConv := old.([]interface{})
		newConv := new.([]interface{})
		for _, el := range oldConv {
			if !isContainsSnapshot(newConv, el) {
				deletedSnapshots = append(deletedSnapshots, el)
			}
		}
		for _, el := range newConv {
			if !isContainsSnapshot(oldConv, el) {
				addedSnapshots = append(addedSnapshots, el)
			} else {
				if isRollback(oldConv, el) {
					updatedSnapshots = append(updatedSnapshots, el)
				}
			}
		}

		if len(deletedSnapshots) > 0 {
			for _, snapshot := range deletedSnapshots {
				snapshotConv := snapshot.(map[string]interface{})
				urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
				urlValues.Add("label", snapshotConv["label"].(string))
				_, err := c.DecortAPICall(ctx, "POST", bserviceSnapshotDeleteAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

		if len(addedSnapshots) > 0 {
			for _, snapshot := range addedSnapshots {
				snapshotConv := snapshot.(map[string]interface{})
				urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
				urlValues.Add("label", snapshotConv["label"].(string))
				_, err := c.DecortAPICall(ctx, "POST", bserviceSnapshotCreateAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

		if len(updatedSnapshots) > 0 {
			for _, snapshot := range updatedSnapshots {
				snapshotConv := snapshot.(map[string]interface{})
				urlValues.Add("serviceId", strconv.Itoa(d.Get("service_id").(int)))
				urlValues.Add("label", snapshotConv["label"].(string))
				_, err := c.DecortAPICall(ctx, "POST", bserviceSnapshotRollbackAPI, urlValues)
				if err != nil {
					return diag.FromErr(err)
				}

				urlValues = &url.Values{}
			}
		}

	}

	return nil
}

func isContainsSnapshot(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["guid"].(string) == elConv["guid"].(string) {
			return true
		}
	}
	return false
}

func isRollback(els []interface{}, el interface{}) bool {
	for _, elOld := range els {
		elOldConv := elOld.(map[string]interface{})
		elConv := el.(map[string]interface{})
		if elOldConv["guid"].(string) == elConv["guid"].(string) &&
			elOldConv["rollback"].(bool) != elConv["rollback"].(bool) &&
			elConv["rollback"].(bool) {
			return true
		}
	}
	return false
}

func resourceBasicServiceSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the service",
		},
		"rg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the Resource Group where this service will be placed",
		},
		"ssh_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "SSH key to deploy for the specified user. Same key will be deployed to all computes of the service.",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the user to deploy SSH key for. Pass empty string if no SSH key deployment is required",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "if set to False, Basic service will be deleted to recycle bin. Otherwise destroyed immediately",
		},
		"enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "if set to False, Basic service will be deleted to recycle bin. Otherwise destroyed immediately",
		},
		"restore": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Restores BasicService instance",
		},
		"start": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Start service. Starting a service technically means starting computes from all service groups according to group relations",
		},
		"service_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"account_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"base_domain": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"compgroup_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"compgroup_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"compgroup_role": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"cpu_total": {
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
		"deleted_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disk_total": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"gid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"groups_name": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"parent_srv_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"parent_srv_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ram_total": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"snapshots": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"label": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"rollback": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
					"timestamp": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"valid": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
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
	}
}

func ResourceBasicService() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceBasicServiceCreate,
		ReadContext:   resourceBasicServiceRead,
		UpdateContext: resourceBasicServiceUpdate,
		DeleteContext: resourceBasicServiceDelete,

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

		Schema: resourceBasicServiceSchemaMake(),
	}
}
