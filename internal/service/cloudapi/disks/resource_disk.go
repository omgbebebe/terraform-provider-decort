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

package disks

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/constants"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/dc"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/status"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDiskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	haveAccount, err := existAccountID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveAccount {
		return diag.Errorf("resourceDiskCreate: can't create Disk because AccountID %d is not allowed or does not exist", d.Get("account_id").(int))
	}

	haveGID, err := existGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveGID {
		return diag.Errorf("resourceDiskCreate: can't create Disk because GID %d is not allowed or does not exist", d.Get("gid").(int))
	}

	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))
	urlValues.Add("gid", strconv.Itoa(d.Get("gid").(int)))
	urlValues.Add("name", d.Get("disk_name").(string))
	urlValues.Add("size", strconv.Itoa(d.Get("size_max").(int)))
	if typeRaw, ok := d.GetOk("type"); ok {
		urlValues.Add("type", strings.ToUpper(typeRaw.(string)))
	} else {
		urlValues.Add("type", "D")
	}

	if sepId, ok := d.GetOk("sep_id"); ok {
		urlValues.Add("sep_id", strconv.Itoa(sepId.(int)))
	}

	if poolName, ok := d.GetOk("pool"); ok {
		urlValues.Add("pool", poolName.(string))
	}

	argVal, argSet := d.GetOk("desc")
	if argSet {
		urlValues.Add("description", argVal.(string))
	}

	diskId, err := c.DecortAPICall(ctx, "POST", disksCreateAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	urlValues = &url.Values{}

	d.SetId(diskId)

	if iotuneRaw, ok := d.GetOk("iotune"); ok {
		iot := iotuneRaw.([]interface{})[0]
		iotune := iot.(map[string]interface{})
		urlValues.Add("diskId", diskId)
		urlValues.Add("iops", strconv.Itoa(iotune["total_iops_sec"].(int)))
		urlValues.Add("read_bytes_sec", strconv.Itoa(iotune["read_bytes_sec"].(int)))
		urlValues.Add("read_bytes_sec_max", strconv.Itoa(iotune["read_bytes_sec_max"].(int)))
		urlValues.Add("read_iops_sec", strconv.Itoa(iotune["read_iops_sec"].(int)))
		urlValues.Add("read_iops_sec_max", strconv.Itoa(iotune["read_iops_sec_max"].(int)))
		urlValues.Add("size_iops_sec", strconv.Itoa(iotune["size_iops_sec"].(int)))
		urlValues.Add("total_bytes_sec", strconv.Itoa(iotune["total_bytes_sec"].(int)))
		urlValues.Add("total_bytes_sec_max", strconv.Itoa(iotune["total_bytes_sec_max"].(int)))
		urlValues.Add("total_iops_sec_max", strconv.Itoa(iotune["total_iops_sec_max"].(int)))
		urlValues.Add("write_bytes_sec", strconv.Itoa(iotune["write_bytes_sec"].(int)))
		urlValues.Add("write_bytes_sec_max", strconv.Itoa(iotune["write_bytes_sec_max"].(int)))
		urlValues.Add("write_iops_sec", strconv.Itoa(iotune["write_iops_sec"].(int)))
		urlValues.Add("write_iops_sec_max", strconv.Itoa(iotune["write_iops_sec_max"].(int)))

		_, err := c.DecortAPICall(ctx, "POST", disksIOLimitAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if shareable := d.Get("shareable"); shareable.(bool) == true {
		urlValues.Add("diskId", diskId)
		_, err := c.DecortAPICall(ctx, "POST", disksShareAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}
		urlValues = &url.Values{}
	}

	dgn := resourceDiskRead(ctx, d, m)
	if dgn != nil {
		return dgn
	}

	return nil
}

func resourceDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlValues := &url.Values{}
	c := m.(*controller.ControllerCfg)
	warnings := dc.Warnings{}

	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hasChangeState := false

	switch disk.Status {
	case status.Destroyed, status.Purged:
		d.Set("disk_id", 0)
		return resourceDiskCreate(ctx, d, m)
	case status.Deleted:
		hasChangeState = true
		urlValues.Add("diskId", d.Id())
		urlValues.Add("reason", d.Get("reason").(string))

		_, err := c.DecortAPICall(ctx, "POST", disksRestoreAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	case status.Assigned:
	case status.Modeled:
		return diag.Errorf("The disk is in status: %s, please, contact support for more information", disk.Status)
	case status.Creating:
	case status.Created:
	case status.Allocated:
	case status.Unallocated:
	}

	if hasChangeState {
		disk, err = utilityDiskCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	flattenDisk(d, *disk)

	return warnings.Get()
}

func resourceDiskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*controller.ControllerCfg)
	warnings := dc.Warnings{}
	urlValues := &url.Values{}

	haveAccount, err := existAccountID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveAccount {
		return diag.Errorf("resourceDiskUpdate: can't update Disk because AccountID %d is not allowed or does not exist", d.Get("account_id").(int))
	}

	haveGID, err := existGID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !haveGID {
		return diag.Errorf("resourceDiskUpdate: can't update Disk because GID %d is not allowed or does not exist", d.Get("gid").(int))
	}

	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	hasChangeState := false

	switch disk.Status {
	case status.Destroyed, status.Purged:
		d.Set("disk_id", 0)
		return resourceDiskCreate(ctx, d, m)
	case status.Deleted:
		hasChangeState = true
		urlValues.Add("diskId", d.Id())
		urlValues.Add("reason", d.Get("reason").(string))

		_, err := c.DecortAPICall(ctx, "POST", disksRestoreAPI, urlValues)
		if err != nil {
			warnings.Add(err)
		}
	case status.Assigned:
	case status.Modeled:
		return diag.Errorf("The disk is in status: %s, please, contact support for more information", disk.Status)
	case status.Creating:
	case status.Created:
	case status.Allocated:
	case status.Unallocated:
	}

	if hasChangeState {
		disk, err = utilityDiskCheckPresence(ctx, d, m)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	if d.HasChange("size_max") {
		oldSize, newSize := d.GetChange("size_max")
		if oldSize.(int) < newSize.(int) {
			log.Debugf("resourceDiskUpdate: resizing disk ID %s - %d GB -> %d GB",
				d.Id(), oldSize.(int), newSize.(int))
			urlValues.Add("diskId", d.Id())
			urlValues.Add("size", strconv.Itoa(newSize.(int)))
			_, err := c.DecortAPICall(ctx, "POST", disksResizeAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
			d.Set("size_max", newSize)
		} else if oldSize.(int) > newSize.(int) {
			return diag.FromErr(fmt.Errorf("resourceDiskUpdate: Disk ID %s - reducing disk size is not allowed", d.Id()))
		}
		urlValues = &url.Values{}
	}

	if d.HasChange("disk_name") {
		urlValues.Add("diskId", d.Id())
		urlValues.Add("name", d.Get("disk_name").(string))
		_, err := c.DecortAPICall(ctx, "POST", disksRenameAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("iotune") {
		iot := d.Get("iotune").([]interface{})[0]
		iotune := iot.(map[string]interface{})
		urlValues.Add("diskId", d.Id())
		urlValues.Add("iops", strconv.Itoa(iotune["total_iops_sec"].(int)))
		urlValues.Add("read_bytes_sec", strconv.Itoa(iotune["read_bytes_sec"].(int)))
		urlValues.Add("read_bytes_sec_max", strconv.Itoa(iotune["read_bytes_sec_max"].(int)))
		urlValues.Add("read_iops_sec", strconv.Itoa(iotune["read_iops_sec"].(int)))
		urlValues.Add("read_iops_sec_max", strconv.Itoa(iotune["read_iops_sec_max"].(int)))
		urlValues.Add("size_iops_sec", strconv.Itoa(iotune["size_iops_sec"].(int)))
		urlValues.Add("total_bytes_sec", strconv.Itoa(iotune["total_bytes_sec"].(int)))
		urlValues.Add("total_bytes_sec_max", strconv.Itoa(iotune["total_bytes_sec_max"].(int)))
		urlValues.Add("total_iops_sec_max", strconv.Itoa(iotune["total_iops_sec_max"].(int)))
		urlValues.Add("write_bytes_sec", strconv.Itoa(iotune["write_bytes_sec"].(int)))
		urlValues.Add("write_bytes_sec_max", strconv.Itoa(iotune["write_bytes_sec_max"].(int)))
		urlValues.Add("write_iops_sec", strconv.Itoa(iotune["write_iops_sec"].(int)))
		urlValues.Add("write_iops_sec_max", strconv.Itoa(iotune["write_iops_sec_max"].(int)))

		_, err := c.DecortAPICall(ctx, "POST", disksIOLimitAPI, urlValues)
		if err != nil {
			return diag.FromErr(err)
		}

		urlValues = &url.Values{}
	}

	if d.HasChange("shareable") {
		oldShare, newShare := d.GetChange("shareable")
		urlValues = &url.Values{}
		urlValues.Add("diskId", d.Id())
		if oldShare.(bool) == false && newShare.(bool) == true {
			_, err := c.DecortAPICall(ctx, "POST", disksShareAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if oldShare.(bool) == true && newShare.(bool) == false {
			_, err := c.DecortAPICall(ctx, "POST", disksUnshareAPI, urlValues)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceDiskRead(ctx, d, m)
}

func resourceDiskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	disk, err := utilityDiskCheckPresence(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if disk.Status == status.Destroyed || disk.Status == status.Purged {
		return nil
	}
	params := &url.Values{}
	params.Add("diskId", d.Id())
	params.Add("detach", strconv.FormatBool(d.Get("detach").(bool)))
	params.Add("permanently", strconv.FormatBool(d.Get("permanently").(bool)))
	params.Add("reason", d.Get("reason").(string))

	c := m.(*controller.ControllerCfg)
	_, err = c.DecortAPICall(ctx, "POST", disksDeleteAPI, params)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDiskSchemaMake() map[string]*schema.Schema {
	rets := map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The unique ID of the subscriber-owner of the disk",
		},
		"disk_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of disk",
		},
		"size_max": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Size in GB",
		},
		"gid": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the grid (platform)",
		},
		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Pool for disk location",
		},
		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Storage endpoint provider ID to create disk",
		},
		"desc": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Description of disk",
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"D", "B", "T"}, false),
			Description:  "The type of disk in terms of its role in compute: 'B=Boot, D=Data, T=Temp'",
		},
		"detach": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Detaching the disk from compute",
		},
		"permanently": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to completely delete the disk, works only with non attached disks",
		},
		"reason": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Reason for deletion",
		},
		"shareable": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"iotune": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"read_bytes_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Number of bytes to read per second",
					},
					"read_bytes_sec_max": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Maximum number of bytes to read",
					},
					"read_iops_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Number of io read operations per second",
					},
					"read_iops_sec_max": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Maximum number of io read operations",
					},
					"size_iops_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Size of io operations",
					},
					"total_bytes_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Total size bytes per second",
					},
					"total_bytes_sec_max": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Maximum total size of bytes per second",
					},
					"total_iops_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Total number of io operations per second",
					},
					"total_iops_sec_max": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Maximum total number of io operations per second",
					},
					"write_bytes_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Number of bytes to write per second",
					},
					"write_bytes_sec_max": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Maximum number of bytes to write per second",
					},
					"write_iops_sec": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Number of write operations per second",
					},
					"write_iops_sec_max": {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "Maximum number of write operations per second",
					},
				},
			},
		},
		"present_to": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"disk_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk ID. Duplicates the value of the ID parameter",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the subscriber '(account') to whom this disk belongs",
		},
		"acl": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"boot_partition": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of disk partitions",
		},
		"computes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"compute_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"compute_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"created_time": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Created time",
		},
		"deleted_time": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Deleted time",
		},
		"destruction_time": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Time of final deletion",
		},
		"devicename": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the device",
		},
		"disk_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Disk path",
		},
		"guid": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk ID on the storage side",
		},
		"image_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Image ID",
		},
		"images": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "IDs of images using the disk",
		},
		"iqn": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Disk IQN",
		},
		"login": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Login to access the disk",
		},
		"milestones": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Milestones",
		},

		"order": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk order",
		},
		"params": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Disk params",
		},
		"parent_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the parent disk",
		},
		"passwd": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Password to access the disk",
		},
		"pci_slot": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the pci slot to which the disk is connected",
		},
		"purge_attempts": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of deletion attempts",
		},
		"purge_time": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Time of the last deletion attempt",
		},
		"reality_device_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Reality device number",
		},
		"reference_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the reference to the disk",
		},
		"res_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Resource ID",
		},
		"res_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the resource",
		},
		"role": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Disk role",
		},
		"sep_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type SEP. Defines the type of storage system and contains one of the values set in the cloud platform",
		},
		"size_used": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "Number of used space, in GB",
		},
		"snapshots": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the snapshot",
					},
					"label": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the snapshot",
					},
					"res_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Reference to the snapshot",
					},
					"snap_set_guid": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The set snapshot ID",
					},
					"snap_set_time": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The set time of the snapshot",
					},
					"timestamp": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Snapshot time",
					},
				},
			},
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Disk status",
		},
		"tech_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Technical status of the disk",
		},
		"vmid": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Virtual Machine ID (Deprecated)",
		},
	}

	return rets
}

func ResourceDisk() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: resourceDiskCreate,
		ReadContext:   resourceDiskRead,
		UpdateContext: resourceDiskUpdate,
		DeleteContext: resourceDiskDelete,

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

		Schema: resourceDiskSchemaMake(),
	}
}
