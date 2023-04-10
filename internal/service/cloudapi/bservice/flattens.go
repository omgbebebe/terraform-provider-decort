package bservice

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func flattenService(d *schema.ResourceData, bs *BasicServiceExtend) {
	d.Set("account_id", bs.AccountId)
	d.Set("account_name", bs.AccountName)
	d.Set("base_domain", bs.BaseDomain)
	d.Set("computes", flattenBasicServiceComputes(bs.Computes))
	d.Set("cpu_total", bs.CPUTotal)
	d.Set("created_by", bs.CreatedBy)
	d.Set("created_time", bs.CreatedTime)
	d.Set("deleted_by", bs.DeletedBy)
	d.Set("deleted_time", bs.DeletedTime)
	d.Set("disk_total", bs.DiskTotal)
	d.Set("gid", bs.GID)
	d.Set("groups", bs.Groups)
	d.Set("groups_name", bs.GroupsName)
	d.Set("guid", bs.GUID)
	d.Set("milestones", bs.Milestones)
	d.Set("service_name", bs.Name)
	d.Set("service_id", bs.ID)
	d.Set("parent_srv_id", bs.ParentSrvId)
	d.Set("parent_srv_type", bs.ParentSrvType)
	d.Set("ram_total", bs.RamTotal)
	d.Set("rg_id", bs.RGID)
	d.Set("rg_name", bs.RGName)
	d.Set("snapshots", flattenBasicServiceSnapshots(bs.Snapshots))
	d.Set("ssh_key", bs.SSHKey)
	d.Set("ssh_user", bs.SSHUser)
	d.Set("status", bs.Status)
	d.Set("tech_status", bs.TechStatus)
	d.Set("updated_by", bs.UpdatedBy)
	d.Set("updated_time", bs.UpdatedTime)
	d.Set("user_managed", bs.UserManaged)
}

func flattenBasicServiceComputes(bscs BasicServiceComputes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bsc := range bscs {
		temp := map[string]interface{}{
			"compgroup_id":   bsc.CompGroupId,
			"compgroup_name": bsc.CompGroupName,
			"compgroup_role": bsc.CompGroupRole,
			"id":             bsc.ID,
			"name":           bsc.Name,
		}
		res = append(res, temp)
	}

	return res
}

func flattenBasicServiceSnapshots(bsrvss BasicServiceSnapshots) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bsrvs := range bsrvss {
		temp := map[string]interface{}{
			"guid":      bsrvs.GUID,
			"label":     bsrvs.Label,
			"timestamp": bsrvs.Timestamp,
			"valid":     bsrvs.Valid,
		}
		res = append(res, temp)
	}
	return res
}
