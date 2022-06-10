/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Stanislav Solovev, <spsolovev@digitalenergy.online>

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
This file is part of Terraform (by Hashicorp) provider for Digital Energy Cloud Orchestration
Technology platfom.

Visit https://github.com/rudecs/terraform-provider-decort for full source code package and updates.
*/

package decort

import (
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceCDROMImageCreate(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceCDROMImageCreate: called for image %s", d.Get("name").(string))

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("name", d.Get("name").(string))
	urlValues.Add("url", d.Get("url").(string))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))

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
		urlValues.Add("pool_name", poolName.(string))
	}
	if architecture, ok := d.GetOk("architecture"); ok {
		urlValues.Add("architecture", architecture.(string))
	}

	imageId, err := controller.decortAPICall("POST", imageCreateCDROMAPI, urlValues)
	if err != nil {
		return err
	}

	d.SetId(imageId)
	d.Set("image_id", imageId)

	image, err := utilityImageCheckPresence(d, m)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(image.ImageId))
	d.Set("bootable", image.Bootable)
	//d.Set("image_id", image.ImageId)

	err = resourceImageRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceCDROMImageDelete(d *schema.ResourceData, m interface{}) error {
	log.Debugf("resourceCDROMImageDelete: called for %s, id: %s", d.Get("name").(string), d.Id())

	image, err := utilityImageCheckPresence(d, m)
	if image == nil {
		if err != nil {
			return err
		}
		return nil
	}

	controller := m.(*ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("imageId", strconv.Itoa(d.Get("image_id").(int)))

	if permanently, ok := d.GetOk("permanently"); ok {
		urlValues.Add("permanently", strconv.FormatBool(permanently.(bool)))
	}

	_, err = controller.decortAPICall("POST", imageDeleteCDROMAPI, urlValues)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceCDROMImageSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the rescue disk",
		},
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "URL where to download ISO from",
		},
		"gid": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "grid (platform) ID where this template should be create in",
		},
		"boot_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Boot type of image bios or uefi",
		},
		"image_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image type linux, windows or other",
		},
		"drivers": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of types of compute suitable for image. Example: [ \"KVM_X86\" ]",
		},
		"meta": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "meta",
		},
		"hot_resize": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Does this machine supports hot resize",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional username for the image",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Optional password for the image",
		},
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "AccountId to make the image exclusive",
		},
		"username_dl": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "username for upload binary media",
		},
		"password_dl": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "password for upload binary media",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "storage endpoint provider ID",
		},
		"pool_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "pool for image create",
		},
		"architecture": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "binary architecture of this image, one of X86_64 of PPC64_LE",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "image id",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Whether to completely delete the image",
		},
		"bootable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Does this image boot OS",
		},
		"unc_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "unc path",
		},
		"link_to": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "status",
		},
		"tech_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "tech atatus",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "version",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "image size",
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"computeci_id": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"milestones": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"provider_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"purge_attempts": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"res_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"res_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rescuecd": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"desc": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"shared_with": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"enabled_stacks": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"history": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"timestamp": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func resourceCDROMImage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceCDROMImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageEdit,
		Delete: resourceCDROMImageDelete,
		Exists: resourceImageExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  &Timeout60s,
			Read:    &Timeout30s,
			Update:  &Timeout60s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		CustomizeDiff: customdiff.All(
			customdiff.IfValueChange("enabled", func(old, new, meta interface{}) bool {
				return old.(bool) != new.(bool)
			}, resourceImageChangeEnabled),
			customdiff.IfValueChange("name", func(old, new, meta interface{}) bool {
				return old.(string) != new.(string) && old.(string) != ""
			}, resourceImageEditName),
			customdiff.IfValueChange("shared_with", func(old, new, meta interface{}) bool {
				o := old.([]interface{})
				n := new.([]interface{})

				if len(o) != len(n) {
					return true
				} else if len(o) == 0 {
					return false
				}
				count := 0
				for i, v := range n {
					if v.(int) == o[i].(int) {
						count++
					}
				}
				return count == 0
			}, resourceImageShare),
			customdiff.IfValueChange("computeci_id", func(old, new, meta interface{}) bool {
				return old.(int) != new.(int)
			}, resourceImageChangeComputeci),
			customdiff.IfValueChange("enabled_stacks", func(old, new, meta interface{}) bool {
				o := old.([]interface{})
				n := new.([]interface{})

				if len(o) != len(n) {
					return true
				} else if len(o) == 0 {
					return false
				}
				count := 0
				for i, v := range n {
					if v.(string) == o[i].(string) {
						count++
					}
				}
				return count == 0
			}, resourceImageUpdateNodes),
		),

		Schema: resourceCDROMImageSchemaMake(),
	}
}
