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

package image

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/status"
)

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageCreate: called for image %s", d.Get("name").(string))

	haveGID, err := existGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveGID {
		return diag.Errorf("resourceImageCreate: can't create Image because GID %d is not allowed or does not exist", d.Get("gid").(int))
	}

	if _, ok := d.GetOk("account_id"); ok {
		haveAccountID, err := existAccountID(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}

		if !haveAccountID {
			return diag.Errorf("resourceImageCreate: can't create Image because AccountID %d is not allowed or does not exist", d.Get("account_id").(int))
		}
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("url", d.Get("url").(string))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))
	urlValues.Add("boottype", d.Get("boot_type").(string))
	urlValues.Add("imagetype", d.Get("type").(string))

	tstr := d.Get("drivers").([]interface{})
	temp := ""
	l := len(tstr)
	for i, str := range tstr {
		s := "\"" + str.(string) + "\""
		if i != (l - 1) {
			s += ","
		}
		temp = temp + s
	}
	temp = "[" + temp + "]"
	urlValues.Add("drivers", temp)

	if hotresize, ok := d.GetOk("hot_resize"); ok {
		urlValues.Add("hotresize", strconv.FormatBool(hotresize.(bool)))
	}
	if username, ok := d.GetOk("username"); ok {
		urlValues.Add("username", username.(string))
	}
	if password, ok := d.GetOk("password"); ok {
		urlValues.Add("password", password.(string))
	}
	if accountId, ok := d.GetOk("account_id"); ok {
		urlValues.Add("accountId", strconv.Itoa(accountId.(int)))
	}
	if usernameDL, ok := d.GetOk("username_dl"); ok {
		urlValues.Add("usernameDL", usernameDL.(string))
	}
	if passwordDL, ok := d.GetOk("password_dl"); ok {
		urlValues.Add("passwordDL", passwordDL.(string))
	}
	if sepId, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sepId", strconv.Itoa(sepId.(int)))
	}
	if poolName, ok := d.GetOk("pool_name"); ok {
		urlValues.Add("poolName", poolName.(string))
	}
	if architecture, ok := d.GetOk("architecture"); ok {
		urlValues.Add("architecture", architecture.(string))
	}

	res, err := c.DecortAPICall(ctx, "POST", imageCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	i := make([]interface{}, 0)
	err = json.Unmarshal([]byte(res), &i)
	if err != nil {
		return diag.FromErr(err)
	}
	imageId := strconv.Itoa(int(i[1].(float64)))
	// end innovation

	d.SetId(imageId)
	d.Set("image_id", imageId)

	_, err = utilityImageCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := resourceImageRead(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func resourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageRead: called for %s id: %s", d.Get("name").(string), d.Id())

	img, err := utilityImageCheckPresence(ctx, d, m)
	if img == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	switch img.Status {
	case status.Modeled:
		return diag.Errorf("The image is in status: %s, please, contact support for more information", img.Status)
	case status.Creating:
	case status.Created:
	case status.Destroyed, status.Purged:
		d.SetId("")
		return resourceImageCreate(ctx, d, m)
	}

	flattenImage(d, img)

	return nil
}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(ctx, d, m)
	if image == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))

	if permanently, ok := d.GetOk("permanently"); ok {
		urlValues.Add("permanently", strconv.FormatBool(permanently.(bool)))
	}

	_, err = c.DecortAPICall(ctx, "POST", imageDeleteAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourceImageEditName(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceImageEditName: called for %s, id: %s", d.Get("name").(string), d.Id())
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))
	urlValues.Add("name", d.Get("name").(string))
	_, err := c.DecortAPICall(ctx, "POST", imageEditNameAPI, urlValues)
	if err != nil {
		return err
	}

	return nil
}

func resourceImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("resourceImageEdit: called for %s, id: %s", d.Get("name").(string), d.Id())

	haveGID, err := existGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !haveGID {
		return diag.Errorf("resourceImageUpdate: can't update Image because GID %d is not allowed or does not exist", d.Get("gid").(int))
	}

	if _, ok := d.GetOk("account_id"); ok {
		haveAccountID, err := existAccountID(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}

		if !haveAccountID {
			return diag.Errorf("resourceImageUpdate: can't update Image because AccountID %d is not allowed or does not exist", d.Get("account_id").(int))
		}
	}

	image, err := utilityImageCheckPresence(ctx, d, m)
	if image == nil {
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	switch image.Status {
	case status.Modeled:
		return diag.Errorf("The image is in status: %s, please, contact support for more information", image.Status)
	case status.Creating:
	case status.Created:
	case status.Destroyed, status.Purged:
		d.SetId("")
		return resourceImageCreate(ctx, d, m)
	}

	if d.HasChange("name") {
		err := resourceImageEditName(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceImageRead(ctx, d, m)
}

func ResourceImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageUpdate,
		DeleteContext: resourceImageDelete,

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

		Schema: resourceImageSchemaMake(dataSourceImageExtendSchemaMake()),
	}
}
