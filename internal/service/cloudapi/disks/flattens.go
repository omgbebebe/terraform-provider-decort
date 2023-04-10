package disks

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDisk(d *schema.ResourceData, disk Disk) {
	diskAcl, _ := json.Marshal(disk.Acl)

	d.Set("account_id", disk.AccountID)
	d.Set("account_name", disk.AccountName)
	d.Set("acl", string(diskAcl))
	d.Set("boot_partition", disk.BootPartition)
	d.Set("computes", flattenDiskComputes(disk.Computes))
	d.Set("created_time", disk.CreatedTime)
	d.Set("deleted_time", disk.DeletedTime)
	d.Set("desc", disk.Desc)
	d.Set("destruction_time", disk.DestructionTime)
	d.Set("devicename", disk.DeviceName)
	d.Set("disk_path", disk.DiskPath)
	d.Set("gid", disk.GridID)
	d.Set("guid", disk.GUID)
	d.Set("disk_id", disk.ID)
	d.Set("image_id", disk.ImageID)
	d.Set("images", disk.Images)
	d.Set("iotune", flattenIOTune(disk.IOTune))
	d.Set("iqn", disk.IQN)
	d.Set("login", disk.Login)
	d.Set("milestones", disk.Milestones)
	d.Set("disk_name", disk.Name)
	d.Set("order", disk.Order)
	d.Set("params", disk.Params)
	d.Set("parent_id", disk.ParentId)
	d.Set("passwd", disk.Passwd)
	d.Set("pci_slot", disk.PciSlot)
	d.Set("pool", disk.Pool)
	d.Set("present_to", disk.PresentTo)
	d.Set("purge_attempts", disk.PurgeAttempts)
	d.Set("purge_time", disk.PurgeTime)
	d.Set("reality_device_number", disk.RealityDeviceNumber)
	d.Set("reference_id", disk.ReferenceId)
	d.Set("res_id", disk.ResID)
	d.Set("res_name", disk.ResName)
	d.Set("role", disk.Role)
	d.Set("sep_id", disk.SepID)
	d.Set("sep_type", disk.SepType)
	d.Set("size_max", disk.SizeMax)
	d.Set("size_used", disk.SizeUsed)
	d.Set("shareable", disk.Shareable)
	d.Set("snapshots", flattenDiskSnapshotList(disk.Snapshots))
	d.Set("status", disk.Status)
	d.Set("tech_status", disk.TechStatus)
	d.Set("type", disk.Type)
	d.Set("vmid", disk.VMID)
}

func flattenDiskSnapshotList(sl SnapshotList) []interface{} {
	res := make([]interface{}, 0)
	for _, snapshot := range sl {
		temp := map[string]interface{}{
			"guid":          snapshot.Guid,
			"label":         snapshot.Label,
			"res_id":        snapshot.ResId,
			"snap_set_guid": snapshot.SnapSetGuid,
			"snap_set_time": snapshot.SnapSetTime,
			"timestamp":     snapshot.TimeStamp,
		}
		res = append(res, temp)
	}

	return res
}

func flattenDiskList(dl DisksList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range dl {
		diskAcl, _ := json.Marshal(disk.Acl)
		temp := map[string]interface{}{
			"account_id":            disk.AccountID,
			"account_name":          disk.AccountName,
			"acl":                   string(diskAcl),
			"computes":              flattenDiskComputes(disk.Computes),
			"boot_partition":        disk.BootPartition,
			"created_time":          disk.CreatedTime,
			"deleted_time":          disk.DeletedTime,
			"desc":                  disk.Desc,
			"destruction_time":      disk.DestructionTime,
			"devicename":            disk.DeviceName,
			"disk_path":             disk.DiskPath,
			"gid":                   disk.GridID,
			"guid":                  disk.GUID,
			"disk_id":               disk.ID,
			"image_id":              disk.ImageID,
			"images":                disk.Images,
			"iotune":                flattenIOTune(disk.IOTune),
			"iqn":                   disk.IQN,
			"login":                 disk.Login,
			"machine_id":            disk.MachineId,
			"machine_name":          disk.MachineName,
			"milestones":            disk.Milestones,
			"disk_name":             disk.Name,
			"order":                 disk.Order,
			"params":                disk.Params,
			"parent_id":             disk.ParentId,
			"passwd":                disk.Passwd,
			"pci_slot":              disk.PciSlot,
			"pool":                  disk.Pool,
			"present_to":            disk.PresentTo,
			"purge_attempts":        disk.PurgeAttempts,
			"purge_time":            disk.PurgeTime,
			"reality_device_number": disk.RealityDeviceNumber,
			"reference_id":          disk.ReferenceId,
			"res_id":                disk.ResID,
			"res_name":              disk.ResName,
			"role":                  disk.Role,
			"sep_id":                disk.SepID,
			"sep_type":              disk.SepType,
			"shareable":             disk.Shareable,
			"size_max":              disk.SizeMax,
			"size_used":             disk.SizeUsed,
			"snapshots":             flattenDiskSnapshotList(disk.Snapshots),
			"status":                disk.Status,
			"tech_status":           disk.TechStatus,
			"type":                  disk.Type,
			"vmid":                  disk.VMID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenIOTune(iot IOTune) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"read_bytes_sec":      iot.ReadBytesSec,
		"read_bytes_sec_max":  iot.ReadBytesSecMax,
		"read_iops_sec":       iot.ReadIopsSec,
		"read_iops_sec_max":   iot.ReadIopsSecMax,
		"size_iops_sec":       iot.SizeIopsSec,
		"total_bytes_sec":     iot.TotalBytesSec,
		"total_bytes_sec_max": iot.TotalBytesSecMax,
		"total_iops_sec":      iot.TotalIopsSec,
		"total_iops_sec_max":  iot.TotalIopsSecMax,
		"write_bytes_sec":     iot.WriteBytesSec,
		"write_bytes_sec_max": iot.WriteBytesSecMax,
		"write_iops_sec":      iot.WriteIopsSec,
		"write_iops_sec_max":  iot.WriteIopsSecMax,
	}

	res = append(res, temp)
	return res
}

func flattenDiskComputes(computes map[string]string) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for computeKey, computeVal := range computes {
		temp := map[string]interface{}{
			"compute_id":   computeKey,
			"compute_name": computeVal,
		}
		res = append(res, temp)
	}
	return res
}
