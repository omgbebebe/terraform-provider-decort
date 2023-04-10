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

package lb

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

func resourceLBCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBCreate")

	haveRGID, err := existRGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveRGID {
		return diag.Errorf("resourceLBCreate: can't create LB because RGID %d is not allowed or does not exist", d.Get("rg_id").(int))
	}

	haveExtNetID, err := existExtNetID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveExtNetID {
		return diag.Errorf("resourceLBCreate: can't create LB because ExtNetID %d is not allowed or does not exist", d.Get("extnet_id").(int))
	}

	haveVins, err := existViNSID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveVins {
		return diag.Errorf("resourceLBCreate: can't create LB because ViNSID %d is not allowed or does not exist", d.Get("vins_id").(int))
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("extnetId", strconv.Itoa(d.Get("extnet_id").(int)))
	urlValues.Add("vinsId", strconv.Itoa(d.Get("vins_id").(int)))
	urlValues.Add("start", strconv.FormatBool((d.Get("start").(bool))))

	if desc, ok := d.GetOk("desc"); ok {
		urlValues.Add("desc", desc.(string))
	}

	lbId, err := c.DecortAPICall(ctx, "POST", lbCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbId)
	d.Set("lb_id", lbId)

	_, err = utilityLBCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceLBRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	urlValues = &url.Values{}

	if enable, ok := d.GetOk("enable"); ok {
		api := lbDisableAPI
		if enable.(bool) {
			api = lbEnableAPI
		}
		urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	return nil
}

func resourceLBRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBRead")

	c := m.(*controller.ControllerCfg)

	lb, err := utilityLBCheckPresence(ctx, d, m)
	if lb == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChanged := false

	switch lb.Status {
	case status.Modeled:
		return diag.Errorf("The LB is in status: %s, please, contact support for more information", lb.Status)
	case status.Creating:
	case status.Created:
	case status.Deleting:
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("lbId", d.Id())

		_, err := c.DecortAPICall(ctx, "POST", lbRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", lbEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Destroying:
		return diag.Errorf("The LB is in progress with status: %s", lb.Status)
	case status.Destroyed:
		d.SetId("")
		return resourceLBCreate(ctx, d, m)
	case status.Enabled:
	case status.Enabling:
	case status.Disabling:
	case status.Disabled:
		log.Debugf("The LB is in status: %s, troubles may occur with update. Please, enable LB first.", lb.Status)
	case status.Restoring:
	}

	if hasChanged {
		lb, err = utilityLBCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	d.Set("ha_mode", lb.HAMode)
	d.Set("backends", flattenLBBackends(lb.Backends))
	d.Set("created_by", lb.CreatedBy)
	d.Set("created_time", lb.CreatedTime)
	d.Set("deleted_by", lb.DeletedBy)
	d.Set("deleted_time", lb.DeletedTime)
	d.Set("desc", lb.Description)
	d.Set("dp_api_user", lb.DPAPIUser)
	d.Set("extnet_id", lb.ExtnetId)
	d.Set("frontends", flattenFrontends(lb.Frontends))
	d.Set("gid", lb.GID)
	d.Set("guid", lb.GUID)
	d.Set("lb_id", lb.ID)
	d.Set("image_id", lb.ImageId)
	d.Set("milestones", lb.Milestones)
	d.Set("name", lb.Name)
	d.Set("primary_node", flattenNode(lb.PrimaryNode))
	d.Set("rg_id", lb.RGID)
	d.Set("rg_name", lb.RGName)
	d.Set("secondary_node", flattenNode(lb.SecondaryNode))
	d.Set("status", lb.Status)
	d.Set("tech_status", lb.TechStatus)
	d.Set("updated_by", lb.UpdatedBy)
	d.Set("updated_time", lb.UpdatedTime)
	d.Set("vins_id", lb.VinsId)

	return nil
}

func resourceLBDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBDelete")

	lb, err := utilityLBCheckPresence(ctx, d, m)
	if lb == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

	if permanently, ok := d.GetOk("permanently"); ok {
		urlValues.Add("permanently", strconv.FormatBool(permanently.(bool)))
	}

	_, err = c.DecortAPICall(ctx, "POST", lbDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceLBUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceLBEdit")
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	haveRGID, err := existRGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveRGID {
		return diag.Errorf("resourceLBUpdate: can't update LB because RGID %d is not allowed or does not exist", d.Get("rg_id").(int))
	}

	haveExtNetID, err := existExtNetID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveExtNetID {
		return diag.Errorf("resourceLBUpdate: can't update LB because ExtNetID %d is not allowed or does not exist", d.Get("extnet_id").(int))
	}

	haveVins, err := existViNSID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveVins {
		return diag.Errorf("resourceLBUpdate: can't update LB because ViNSID %d is not allowed or does not exist", d.Get("vins_id").(int))
	}

	lb, err := utilityLBCheckPresence(ctx, d, m)
	if lb == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChanged := false

	switch lb.Status {
	case status.Modeled:
		return diag.Errorf("The LB is in status: %s, please, contact support for more information", lb.Status)
	case status.Creating:
	case status.Created:
	case status.Deleting:
	case status.Deleted:
		urlValues := &url.Values{}
		urlValues.Add("lbId", d.Id())

		_, err := c.DecortAPICall(ctx, "POST", lbRestoreAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = c.DecortAPICall(ctx, "POST", lbEnableAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		hasChanged = true
	case status.Destroying:
		return diag.Errorf("The LB is in progress with status: %s", lb.Status)
	case status.Destroyed:
		d.SetId("")
		return resourceLBCreate(ctx, d, m)
	case status.Enabled:
	case status.Enabling:
	case status.Disabling:
	case status.Disabled:
		log.Debugf("The LB is in status: %s, troubles may occur with update. Please, enable LB first.", lb.Status)
	case status.Restoring:
	}

	if hasChanged {
		lb, err = utilityLBCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enable") {
		api := lbDisableAPI
		enable := d.Get("enable").(bool)
		if enable {
			api = lbEnableAPI
		}
		urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("start") {
		api := lbStopAPI
		start := d.Get("start").(bool)
		if start {
			api = lbStartAPI
		}
		urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))

		_, err := c.DecortAPICall(ctx, "POST", api, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("desc") {
		urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
		urlValues.Add("desc", d.Get("desc").(string))

		_, err := c.DecortAPICall(ctx, "POST", lbUpdateAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("restart") {
		restart := d.Get("restart").(bool)
		if restart {
			urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
			_, err := c.DecortAPICall(ctx, "POST", lbRestartAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
		}
	}

	if d.HasChange("restore") {
		restore := d.Get("restore").(bool)
		if restore {
			urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
			_, err := c.DecortAPICall(ctx, "POST", lbRestoreAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
		}
	}

	if d.HasChange("config_reset") {
		cfgReset := d.Get("config_reset").(bool)
		if cfgReset {
			urlValues.Add("lbId", strconv.Itoa(d.Get("lb_id").(int)))
			_, err := c.DecortAPICall(ctx, "POST", lbConfigResetAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}

			urlValues = &url.Values{}
		}
	}

	//TODO: перенести backend и frontend из ресурсов сюда

	return resourceLBRead(ctx, d, m)
}

func ResourceLB() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceLBCreate,
		ReadContext:   resourceLBRead,
		UpdateContext: resourceLBUpdate,
		DeleteContext: resourceLBDelete,

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

		Schema: lbResourceSchemaMake(),
	}
}
