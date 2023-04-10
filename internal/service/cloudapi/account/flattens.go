package account

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/flattens"
)

func flattenAccount(d *schema.ResourceData, acc AccountWithResources) error {
	d.Set("dc_location", acc.DCLocation)
	d.Set("resources", flattenAccResources(acc.Resources))
	d.Set("ckey", acc.CKey)
	d.Set("meta", flattens.FlattenMeta(acc.Meta))
	d.Set("acl", flattenAccAcl(acc.Acl))
	d.Set("company", acc.Company)
	d.Set("companyurl", acc.CompanyUrl)
	d.Set("created_by", acc.CreatedBy)
	d.Set("created_time", acc.CreatedTime)
	d.Set("deactivation_time", acc.DeactiovationTime)
	d.Set("deleted_by", acc.DeletedBy)
	d.Set("deleted_time", acc.DeletedTime)
	d.Set("displayname", acc.DisplayName)
	d.Set("guid", acc.GUID)
	d.Set("account_id", acc.ID)
	d.Set("account_name", acc.Name)
	d.Set("resource_limits", flattenRgResourceLimits(acc.ResourceLimits))
	d.Set("send_access_emails", acc.SendAccessEmails)
	d.Set("service_account", acc.ServiceAccount)
	d.Set("status", acc.Status)
	d.Set("updated_time", acc.UpdatedTime)
	d.Set("version", acc.Version)
	d.Set("vins", acc.Vins)
	d.Set("vinses", acc.Vinses)
	d.Set("computes", flattenAccComputes(acc.Computes))
	d.Set("machines", flattenAccMachines(acc.Machines))

	if username, ok := d.GetOk("username"); ok {
		d.Set("username", username)
	} else {
		d.Set("username", acc.Acl[0].UgroupID)
	}

	return nil
}

func flattenAccComputes(acs Computes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"started": acs.Started,
		"stopped": acs.Stopped,
	}
	res = append(res, temp)
	return res
}

func flattenAccMachines(ams Machines) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"running": ams.Running,
		"halted":  ams.Halted,
	}
	res = append(res, temp)
	return res
}

func flattenAccAcl(acls []AccountAclRecord) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acls := range acls {
		temp := map[string]interface{}{
			"can_be_deleted": acls.CanBeDeleted,
			"explicit":       acls.IsExplicit,
			"guid":           acls.Guid,
			"right":          acls.Rights,
			"status":         acls.Status,
			"type":           acls.Type,
			"user_group_id":  acls.UgroupID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenRgResourceLimits(rl ResourceLimits) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cu_c":      rl.CUC,
		"cu_d":      rl.CUD,
		"cu_i":      rl.CUI,
		"cu_m":      rl.CUM,
		"cu_np":     rl.CUNP,
		"gpu_units": rl.GpuUnits,
	}
	res = append(res, temp)

	return res

}

func flattenAccResources(r Resources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"current":  flattenAccResource(r.Current),
		"reserved": flattenAccResource(r.Reserved),
	}
	res = append(res, temp)
	return res
}

func flattenAccountSeps(seps map[string]map[string]ResourceSep) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for sepKey, sepVal := range seps {
		for dataKey, dataVal := range sepVal {
			temp := map[string]interface{}{
				"sep_id":        sepKey,
				"data_name":     dataKey,
				"disk_size":     dataVal.DiskSize,
				"disk_size_max": dataVal.DiskSizeMax,
			}
			res = append(res, temp)
		}
	}
	return res
}

func flattenAccResource(r Resource) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":        r.CPU,
		"disksize":   r.Disksize,
		"extips":     r.Extips,
		"exttraffic": r.Exttraffic,
		"gpu":        r.GPU,
		"ram":        r.RAM,
		"seps":       flattenAccountSeps(r.SEPs),
	}
	res = append(res, temp)
	return res
}
