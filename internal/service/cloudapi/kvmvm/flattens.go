package kvmvm

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/status"
	log "github.com/sirupsen/logrus"
)

func flattenDisks(disks []InfoDisk) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range disks {
		temp := map[string]interface{}{
			"disk_id":  disk.ID,
			"pci_slot": disk.PCISlot,
		}
		res = append(res, temp)
	}
	return res
}
func flattenQOS(qos QOS) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"e_rate":   qos.ERate,
		"guid":     qos.GUID,
		"in_brust": qos.InBurst,
		"in_rate":  qos.InRate,
	}
	res = append(res, temp)
	return res
}
func flattenInterfaces(interfaces ListInterfaces) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, interfaceItem := range interfaces {
		temp := map[string]interface{}{
			"conn_id":       interfaceItem.ConnID,
			"conn_type":     interfaceItem.ConnType,
			"def_gw":        interfaceItem.DefGW,
			"flip_group_id": interfaceItem.FLIPGroupID,
			"guid":          interfaceItem.GUID,
			"ip_address":    interfaceItem.IPAddress,
			"listen_ssh":    interfaceItem.ListenSSH,
			"mac":           interfaceItem.MAC,
			"name":          interfaceItem.Name,
			"net_id":        interfaceItem.NetID,
			"netmask":       interfaceItem.NetMask,
			"net_type":      interfaceItem.NetType,
			"pci_slot":      interfaceItem.PCISlot,
			"qos":           flattenQOS(interfaceItem.QOS),
			"target":        interfaceItem.Target,
			"type":          interfaceItem.Type,
			"vnfs":          interfaceItem.VNFs,
		}
		res = append(res, temp)
	}
	return res
}
func flattenSnapSets(snapSets ListSnapSets) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, snapSet := range snapSets {
		temp := map[string]interface{}{
			"disks":     snapSet.Disks,
			"guid":      snapSet.GUID,
			"label":     snapSet.Label,
			"timestamp": snapSet.Timestamp,
		}
		res = append(res, temp)
	}
	return res
}
func flattenTags(tags map[string]string) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for key, val := range tags {
		temp := map[string]interface{}{
			"key": key,
			"val": val,
		}
		res = append(res, temp)
	}
	return res
}

func flattenListRules(listRules ListRules) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rule := range listRules {
		temp := map[string]interface{}{
			"guid":     rule.GUID,
			"key":      rule.Key,
			"mode":     rule.Mode,
			"policy":   rule.Policy,
			"topology": rule.Topology,
			"value":    rule.Value,
		}
		res = append(res, temp)
	}
	return res
}
func flattenListACL(listAcl ListACL) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acl := range listAcl {
		var explicit interface{}
		switch acl.Explicit.(type) { //Платформенный хак
		case bool:
			explicit = acl.Explicit.(bool)
		case string:
			explicit, _ = strconv.ParseBool(acl.Explicit.(string))
		}
		temp := map[string]interface{}{
			"explicit":      explicit,
			"guid":          acl.GUID,
			"right":         acl.Right,
			"status":        acl.Status,
			"type":          acl.Type,
			"user_group_id": acl.UserGroupID,
		}
		res = append(res, temp)
	}
	return res
}
func flattenComputeList(computes ListComputes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, compute := range computes {
		customFields, _ := json.Marshal(compute.CustomFields)
		devices, _ := json.Marshal(compute.Devices)
		temp := map[string]interface{}{
			"acl":                 flattenListACL(compute.ACL),
			"account_id":          compute.AccountID,
			"account_name":        compute.AccountName,
			"affinity_label":      compute.AffinityLabel,
			"affinity_rules":      flattenListRules(compute.AffinityRules),
			"affinity_weight":     compute.AffinityWeight,
			"anti_affinity_rules": flattenListRules(compute.AntiAffinityRules),
			"arch":                compute.Architecture,
			"boot_order":          compute.BootOrder,
			"bootdisk_size":       compute.BootDiskSize,
			"clone_reference":     compute.CloneReference,
			"clones":              compute.Clones,
			"computeci_id":        compute.ComputeCIID,
			"cpus":                compute.CPU,
			"created_by":          compute.CreatedBy,
			"created_time":        compute.CreatedTime,
			"custom_fields":       string(customFields),
			"deleted_by":          compute.DeletedBy,
			"deleted_time":        compute.DeletedTime,
			"desc":                compute.Description,
			"devices":             string(devices),
			"disks":               flattenDisks(compute.Disks),
			"driver":              compute.Driver,
			"gid":                 compute.GID,
			"guid":                compute.GUID,
			"compute_id":          compute.ID,
			"image_id":            compute.ImageID,
			"interfaces":          flattenInterfaces(compute.Interfaces),
			"lock_status":         compute.LockStatus,
			"manager_id":          compute.ManagerID,
			"manager_type":        compute.ManagerType,
			"migrationjob":        compute.MigrationJob,
			"milestones":          compute.Milestones,
			"name":                compute.Name,
			"pinned":              compute.Pinned,
			"ram":                 compute.RAM,
			"reference_id":        compute.ReferenceID,
			"registered":          compute.Registered,
			"res_name":            compute.ResName,
			"rg_id":               compute.RGID,
			"rg_name":             compute.RGName,
			"snap_sets":           flattenSnapSets(compute.SnapSets),
			"stateless_sep_id":    compute.StatelessSepID,
			"stateless_sep_type":  compute.StatelessSepType,
			"status":              compute.Status,
			"tags":                flattenTags(compute.Tags),
			"tech_status":         compute.TechStatus,
			"total_disk_size":     compute.TotalDiskSize,
			"updated_by":          compute.UpdatedBy,
			"updated_time":        compute.UpdatedTime,
			"user_managed":        compute.UserManaged,
			"vgpus":               compute.VGPUs,
			"vins_connected":      compute.VINSConnected,
			"virtual_image_id":    compute.VirtualImageID,
		}
		res = append(res, temp)
	}

	return res
}

func flattenComputeDisksDemo(disksList ListComputeDisks, extraDisks []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(disksList))
	for _, disk := range disksList {
		if disk.Name == "bootdisk" || findInExtraDisks(uint(disk.ID), extraDisks) { //skip main bootdisk and extraDisks
			continue
		}
		temp := map[string]interface{}{
			"disk_name": disk.Name,
			"disk_id":   disk.ID,
			"disk_type": disk.Type,
			"sep_id":    disk.SepID,
			"shareable": disk.Shareable,
			"size_max":  disk.SizeMax,
			"size_used": disk.SizeUsed,
			"pool":      disk.Pool,
			"desc":      disk.Description,
			"image_id":  disk.ImageID,
			"size":      disk.SizeMax,
		}
		res = append(res, temp)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i]["disk_id"].(uint64) < res[j]["disk_id"].(uint64)
	})
	return res
}

func flattenNetwork(interfaces ListInterfaces) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(interfaces))
	//index := 0
	for _, network := range interfaces {
		temp := map[string]interface{}{
			"net_id":     network.NetID,
			"net_type":   network.NetType,
			"ip_address": network.IPAddress,
			"mac":        network.MAC,
		}
		res = append(res, temp)
	}
	return res
}

func findBootDisk(disks ListComputeDisks) *ItemComputeDisk {
	for _, disk := range disks {
		if disk.Name == "bootdisk" {
			return &disk
		}
	}
	return nil
}

func flattenCompute(d *schema.ResourceData, compute RecordCompute) error {
	// This function expects that compFacts string contains response from API compute/get,
	// i.e. detailed information about compute instance.
	//
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceComputeExists(...) method
	log.Debugf("flattenCompute: ID %d, RG ID %d", compute.ID, compute.RGID)

	devices, _ := json.Marshal(compute.Devices)
	userdata, _ := json.Marshal(compute.Userdata)

	//check extraDisks, ipa_type, is,
	d.SetId(strconv.FormatUint(compute.ID, 10))
	d.Set("acl", flattenACL(compute.ACL))
	d.Set("account_id", compute.AccountID)
	d.Set("account_name", compute.AccountName)
	d.Set("affinity_weight", compute.AffinityWeight)
	d.Set("arch", compute.Architecture)
	d.Set("boot_order", compute.BootOrder)
	d.Set("boot_disk_size", compute.BootDiskSize)
	bootDisk := findBootDisk(compute.Disks)
	d.Set("boot_disk_id", bootDisk.ID)
	d.Set("sep_id", bootDisk.SepID)
	d.Set("pool", bootDisk.Pool)
	d.Set("clone_reference", compute.CloneReference)
	d.Set("clones", compute.Clones)
	if string(userdata) != "{}" {
		d.Set("cloud_init", string(userdata))
	}
	d.Set("computeci_id", compute.ComputeCIID)
	d.Set("created_by", compute.CreatedBy)
	d.Set("created_time", compute.CreatedTime)
	d.Set("custom_fields", flattenCustomFields(compute.CustomFields))
	d.Set("deleted_by", compute.DeletedBy)
	d.Set("deleted_time", compute.DeletedTime)
	d.Set("description", compute.Description)
	d.Set("devices", string(devices))
	err := d.Set("disks", flattenComputeDisksDemo(compute.Disks, d.Get("extra_disks").(*schema.Set).List()))
	if err != nil {
		return err
	}
	d.Set("driver", compute.Driver)
	d.Set("cpu", compute.CPU)
	d.Set("gid", compute.GID)
	d.Set("guid", compute.GUID)
	d.Set("compute_id", compute.ID)
	if compute.VirtualImageID != 0 {
		d.Set("image_id", compute.VirtualImageID)
	} else {
		d.Set("image_id", compute.ImageID)
	}
	d.Set("interfaces", flattenInterfaces(compute.Interfaces))
	d.Set("lock_status", compute.LockStatus)
	d.Set("manager_id", compute.ManagerID)
	d.Set("manager_type", compute.ManagerType)
	d.Set("migrationjob", compute.MigrationJob)
	d.Set("milestones", compute.Milestones)
	d.Set("name", compute.Name)
	d.Set("natable_vins_id", compute.NatableVINSID)
	d.Set("natable_vins_ip", compute.NatableVINSIP)
	d.Set("natable_vins_name", compute.NatableVINSName)
	d.Set("natable_vins_network", compute.NatableVINSNetwork)
	d.Set("natable_vins_network_name", compute.NatableVINSNetworkName)
	if err := d.Set("os_users", parseOsUsers(compute.OSUsers)); err != nil {
		return err
	}
	d.Set("pinned", compute.Pinned)
	d.Set("ram", compute.RAM)
	d.Set("reference_id", compute.ReferenceID)
	d.Set("registered", compute.Registered)
	d.Set("res_name", compute.ResName)
	d.Set("rg_id", compute.RGID)
	d.Set("rg_name", compute.RGName)
	d.Set("snap_sets", flattenSnapSets(compute.SnapSets))
	d.Set("stateless_sep_id", compute.StatelessSepID)
	d.Set("stateless_sep_type", compute.StatelessSepType)
	d.Set("status", compute.Status)
	d.Set("tags", flattenTags(compute.Tags))
	d.Set("tech_status", compute.TechStatus)
	d.Set("updated_by", compute.UpdatedBy)
	d.Set("updated_time", compute.UpdatedTime)
	d.Set("user_managed", compute.UserManaged)
	d.Set("vgpus", compute.VGPUs)
	d.Set("virtual_image_id", compute.VirtualImageID)
	d.Set("virtual_image_name", compute.VirtualImageName)

	d.Set("enabled", false)
	if compute.Status == status.Enabled {
		d.Set("enabled", true)
	}

	d.Set("started", false)
	if compute.TechStatus == "STARTED" {
		d.Set("started", true)
	}

	d.Set("network", flattenNetwork(compute.Interfaces))

	//if len(model.Disks) > 0 {
	//log.Debugf("flattenCompute: calling parseComputeDisksToExtraDisks for %d disks", len(model.Disks))
	//if err = d.Set("extra_disks", parseComputeDisksToExtraDisks(model.Disks)); err != nil {
	//return err
	//}
	//}

	return nil
}

func flattenDataComputeDisksDemo(disksList ListComputeDisks, extraDisks []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range disksList {
		if findInExtraDisks(uint(disk.ID), extraDisks) { //skip main bootdisk and extraDisks
			continue
		}
		temp := map[string]interface{}{
			"disk_name": disk.Name,
			"disk_id":   disk.ID,
			"disk_type": disk.Type,
			"sep_id":    disk.SepID,
			"shareable": disk.Shareable,
			"size_max":  disk.SizeMax,
			"size_used": disk.SizeUsed,
			"pool":      disk.Pool,
			"desc":      disk.Description,
			"image_id":  disk.ImageID,
			"size":      disk.SizeMax,
		}
		res = append(res, temp)
	}
	return res
}

func flattenACL(acl RecordACL) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"account_acl": flattenListACL(acl.AccountACL),
		"compute_acl": flattenListACL(acl.ComputeACL),
		"rg_acl":      flattenListACL(acl.RGACL),
	}
	res = append(res, temp)
	return res
}

func flattenAffinityRules(affinityRules ListRules) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, affinityRule := range affinityRules {
		temp := map[string]interface{}{
			"guid":     affinityRule.GUID,
			"key":      affinityRule.Key,
			"mode":     affinityRule.Mode,
			"policy":   affinityRule.Policy,
			"topology": affinityRule.Topology,
			"value":    affinityRule.Value,
		}
		res = append(res, temp)
	}

	return res
}

func flattenIotune(iotune IOTune) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"read_bytes_sec":      iotune.ReadBytesSec,
		"read_bytes_sec_max":  iotune.ReadBytesSecMax,
		"read_iops_sec":       iotune.ReadIOPSSec,
		"read_iops_sec_max":   iotune.ReadIOPSSecMax,
		"size_iops_sec":       iotune.SizeIOPSSec,
		"total_bytes_sec":     iotune.TotalBytesSec,
		"total_bytes_sec_max": iotune.TotalBytesSecMax,
		"total_iops_sec":      iotune.TotalIOPSSec,
		"total_iops_sec_max":  iotune.TotalIOPSSecMax,
		"write_bytes_sec":     iotune.WriteBytesSec,
		"write_bytes_sec_max": iotune.WriteBytesSecMax,
		"write_iops_sec":      iotune.WriteIOPSSec,
		"write_iops_sec_max":  iotune.WriteIOPSSecMax,
	}
	res = append(res, temp)

	return res
}

func flattenSnapshots(snapshots SnapshotExtendList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, snapshot := range snapshots {
		temp := map[string]interface{}{
			"guid":          snapshot.GUID,
			"label":         snapshot.Label,
			"res_id":        snapshot.ResID,
			"snap_set_guid": snapshot.SnapSetGUID,
			"snap_set_time": snapshot.SnapSetTime,
			"timestamp":     snapshot.TimeStamp,
		}
		res = append(res, temp)
	}

	return res
}

func flattenListComputeDisks(disks ListComputeDisks) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range disks {
		acl, _ := json.Marshal(disk.ACL)
		temp := map[string]interface{}{
			"_ckey":                 disk.CKey,
			"acl":                   string(acl),
			"account_id":            disk.AccountID,
			"boot_partition":        disk.BootPartition,
			"created_time":          disk.CreatedTime,
			"deleted_time":          disk.DeletedTime,
			"description":           disk.Description,
			"destruction_time":      disk.DestructionTime,
			"disk_path":             disk.DiskPath,
			"gid":                   disk.GID,
			"guid":                  disk.GUID,
			"disk_id":               disk.ID,
			"image_id":              disk.ImageID,
			"images":                disk.Images,
			"iotune":                flattenIotune(disk.IOTune),
			"iqn":                   disk.IQN,
			"login":                 disk.Login,
			"milestones":            disk.Milestones,
			"name":                  disk.Name,
			"order":                 disk.Order,
			"params":                disk.Params,
			"parent_id":             disk.ParentID,
			"passwd":                disk.Passwd,
			"pci_slot":              disk.PCISlot,
			"pool":                  disk.Pool,
			"present_to":            disk.PresentTo,
			"purge_time":            disk.PurgeTime,
			"reality_device_number": disk.RealityDeviceNumber,
			"res_id":                disk.ResID,
			"role":                  disk.Role,
			"sep_id":                disk.SepID,
			"shareable":             disk.Shareable,
			"size_max":              disk.SizeMax,
			"size_used":             disk.SizeUsed,
			"snapshots":             flattenSnapshots(disk.Snapshots),
			"status":                disk.Status,
			"tech_status":           disk.TechStatus,
			"type":                  disk.Type,
			"vmid":                  disk.VMID,
		}
		res = append(res, temp)
	}

	return res
}

func flattenCustomFields(customFileds map[string]interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for key, val := range customFileds {
		value, _ := json.Marshal(val)
		temp := map[string]interface{}{
			"key": key,
			"val": string(value),
		}
		res = append(res, temp)
	}
	return res
}
func flattenOsUsers(osUsers ListOSUser) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, user := range osUsers {
		temp := map[string]interface{}{
			"guid":       user.GUID,
			"login":      user.Login,
			"password":   user.Password,
			"public_key": user.PubKey,
		}
		res = append(res, temp)
	}
	return res
}

func flattenDataCompute(d *schema.ResourceData, compute RecordCompute) {
	devices, _ := json.Marshal(compute.Devices)
	userdata, _ := json.Marshal(compute.Userdata)
	d.Set("acl", flattenACL(compute.ACL))
	d.Set("account_id", compute.AccountID)
	d.Set("account_name", compute.AccountName)
	d.Set("affinity_label", compute.AffinityLabel)
	d.Set("affinity_rules", flattenAffinityRules(compute.AffinityRules))
	d.Set("affinity_weight", compute.AffinityWeight)
	d.Set("anti_affinity_rules", flattenListRules(compute.AntiAffinityRules))
	d.Set("arch", compute.Architecture)
	d.Set("boot_order", compute.BootOrder)
	d.Set("bootdisk_size", compute.BootDiskSize)
	d.Set("clone_reference", compute.CloneReference)
	d.Set("clones", compute.Clones)
	d.Set("computeci_id", compute.ComputeCIID)
	d.Set("cpus", compute.CPU)
	d.Set("created_by", compute.CreatedBy)
	d.Set("created_time", compute.CreatedTime)
	d.Set("custom_fields", flattenCustomFields(compute.CustomFields))
	d.Set("deleted_by", compute.DeletedBy)
	d.Set("deleted_time", compute.DeletedTime)
	d.Set("desc", compute.Description)
	d.Set("devices", string(devices))
	d.Set("disks", flattenListComputeDisks(compute.Disks))
	d.Set("driver", compute.Driver)
	d.Set("gid", compute.GID)
	d.Set("guid", compute.GUID)
	d.Set("compute_id", compute.ID)
	d.Set("image_id", compute.ImageID)
	d.Set("interfaces", flattenInterfaces(compute.Interfaces))
	d.Set("lock_status", compute.LockStatus)
	d.Set("manager_id", compute.ManagerID)
	d.Set("manager_type", compute.ManagerType)
	d.Set("migrationjob", compute.MigrationJob)
	d.Set("milestones", compute.Milestones)
	d.Set("name", compute.Name)
	d.Set("natable_vins_id", compute.NatableVINSID)
	d.Set("natable_vins_ip", compute.NatableVINSIP)
	d.Set("natable_vins_name", compute.NatableVINSName)
	d.Set("natable_vins_network", compute.NatableVINSNetwork)
	d.Set("natable_vins_network_name", compute.NatableVINSNetworkName)
	d.Set("os_users", flattenOsUsers(compute.OSUsers))
	d.Set("pinned", compute.Pinned)
	d.Set("ram", compute.RAM)
	d.Set("reference_id", compute.ReferenceID)
	d.Set("registered", compute.Registered)
	d.Set("res_name", compute.ResName)
	d.Set("rg_id", compute.RGID)
	d.Set("rg_name", compute.RGName)
	d.Set("snap_sets", flattenSnapSets(compute.SnapSets))
	d.Set("stateless_sep_id", compute.StatelessSepID)
	d.Set("stateless_sep_type", compute.StatelessSepType)
	d.Set("status", compute.Status)
	d.Set("tags", compute.Tags)
	d.Set("tech_status", compute.TechStatus)
	d.Set("updated_by", compute.UpdatedBy)
	d.Set("updated_time", compute.UpdatedTime)
	d.Set("user_managed", compute.UserManaged)
	d.Set("userdata", string(userdata))
	d.Set("vgpus", compute.VGPUs)
	d.Set("virtual_image_id", compute.VirtualImageID)
	d.Set("virtual_image_name", compute.VirtualImageName)
}

func flattenComputeAudits(computeAudits ListAudits) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, computeAudit := range computeAudits {
		temp := map[string]interface{}{
			"call":         computeAudit.Call,
			"responsetime": computeAudit.ResponseTime,
			"statuscode":   computeAudit.StatusCode,
			"timestamp":    computeAudit.Timestamp,
			"user":         computeAudit.User,
		}
		res = append(res, temp)
	}
	return res
}

func flattenPfwList(computePfws ListPFWs) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, computePfw := range computePfws {
		temp := map[string]interface{}{
			"pfw_id":            computePfw.ID,
			"local_ip":          computePfw.LocalIP,
			"local_port":        computePfw.LocalPort,
			"protocol":          computePfw.Protocol,
			"public_port_end":   computePfw.PublicPortEnd,
			"public_port_start": computePfw.PublicPortStart,
			"vm_id":             computePfw.VMID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenUserList(d *schema.ResourceData, userList RecordACL) {
	d.Set("account_acl", flattenListACL(userList.AccountACL))
	d.Set("compute_acl", flattenListACL(userList.ComputeACL))
	d.Set("rg_acl", flattenListACL(userList.RGACL))
}

func flattenComputeGetAudits(computeAudits ListShortAudits) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, computeAudit := range computeAudits {
		temp := map[string]interface{}{
			"epoch":   computeAudit.Epoch,
			"message": computeAudit.Message,
		}
		res = append(res, temp)
	}
	return res
}
