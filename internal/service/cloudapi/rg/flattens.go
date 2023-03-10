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

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package rg

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

func flattenAccountSeps(seps map[string]map[string]DiskUsage) []map[string]interface{} {
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
		"disksize":   r.DiskSize,
		"extips":     r.ExtIPs,
		"exttraffic": r.ExtTraffic,
		"gpu":        r.GPU,
		"ram":        r.RAM,
		"seps":       flattenAccountSeps(r.SEPs),
	}
	res = append(res, temp)

	return res
}

func flattenRgResources(r Resources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"current":  flattenAccResource(r.Current),
		"reserved": flattenAccResource(r.Reserved),
	}
	res = append(res, temp)
	return res
}

func flattenResgroup(d *schema.ResourceData, details RecordResourceGroup) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceRsgroupExists(...) method
	// log.Debugf("%s", rg_facts)
	//log.Debugf("flattenResgroup: ready to decode response body from API")
	//details := ResgroupGetResp{}
	//err := json.Unmarshal([]byte(rg_facts), &details)
	//if err != nil {
	//return err
	//}

	log.Debugf("flattenResgroup: decoded RG name %q / ID %d, account ID %d",
		details.Name, details.ID, details.AccountID)

	d.SetId(fmt.Sprintf("%d", details.ID))

	d.Set("account_id", details.AccountID)
	d.Set("gid", details.GID)
	d.Set("def_net_type", details.DefNetType)
	d.Set("name", details.Name)

	d.Set("resources", flattenRgResource(details.Resources))
	d.Set("account_name", details.AccountName)
	d.Set("acl", flattenRgAcl(details.ACL))
	d.Set("vms", details.Computes)
	d.Set("created_by", details.CreatedBy)
	d.Set("created_time", details.CreatedTime)
	d.Set("def_net_id", details.DefNetID)
	d.Set("deleted_by", details.DeletedBy)
	d.Set("deleted_time", details.DeletedTime)
	d.Set("description", details.Description)
	d.Set("dirty", details.Dirty)
	d.Set("guid", details.GUID)
	d.Set("rg_id", details.ID)
	d.Set("lock_status", details.LockStatus)
	d.Set("milestones", details.Milestones)
	d.Set("register_computes", details.RegisterComputes)
	d.Set("res_types", details.ResTypes)
	d.Set("secret", details.Secret)
	d.Set("status", details.Status)
	d.Set("updated_by", details.UpdatedBy)
	d.Set("updated_time", details.UpdatedTime)
	d.Set("uniq_pools", details.UniqPools)
	d.Set("vins", details.VINS)

	return nil
}

func flattenRgSeps(seps map[string]map[string]DiskUsage) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for sepKey, sepVal := range seps {
		SepMap := map[string]interface{}{}
		for dataKey, dataVal := range sepVal {
			val, _ := json.Marshal(dataVal)
			SepMap[dataKey] = string(val)
		}
		temp := map[string]interface{}{
			"sep_id": sepKey,
			"map":    SepMap,
		}
		res = append(res, temp)
	}
	return res
}

func flattenResource(resource Resource) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)

	temp := map[string]interface{}{
		"cpu":           resource.CPU,
		"disk_size":     resource.DiskSize,
		"disk_size_max": resource.DiskSizeMax,
		"extips":        resource.ExtIPs,
		"exttraffic":    resource.ExtTraffic,
		"gpu":           resource.GPU,
		"ram":           resource.RAM,
		"seps":          flattenRgSeps(resource.SEPs),
	}

	res = append(res, temp)

	return res
}

func flattenRgResource(itemResource Resources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"current":  flattenResource(itemResource.Current),
		"reserved": flattenResource(itemResource.Reserved),
	}

	res = append(res, temp)
	return res
}

func flattenRg(d *schema.ResourceData, itemRg RecordResourceGroup) {
	d.Set("resources", flattenRgResource(itemRg.Resources))
	d.Set("account_id", itemRg.AccountID)
	d.Set("account_name", itemRg.AccountName)
	d.Set("acl", flattenRgAcl(itemRg.ACL))
	d.Set("computes", itemRg.Computes)
	d.Set("created_by", itemRg.CreatedBy)
	d.Set("created_time", itemRg.CreatedTime)
	d.Set("def_net_id", itemRg.DefNetID)
	d.Set("def_net_type", itemRg.DefNetType)
	d.Set("deleted_by", itemRg.DeletedBy)
	d.Set("deleted_time", itemRg.DeletedTime)
	d.Set("desc", itemRg.Description)
	d.Set("dirty", itemRg.Dirty)
	d.Set("gid", itemRg.GID)
	d.Set("guid", itemRg.GUID)
	d.Set("rg_id", itemRg.ID)
	d.Set("lock_status", itemRg.LockStatus)
	d.Set("milestones", itemRg.Milestones)
	d.Set("name", itemRg.Name)
	d.Set("register_computes", itemRg.RegisterComputes)
	d.Set("res_types", itemRg.ResTypes)
	d.Set("resource_limits", flattenRgResourceLimits(itemRg.ResourceLimits))
	d.Set("secret", itemRg.Secret)
	d.Set("status", itemRg.Status)
	d.Set("updated_by", itemRg.UpdatedBy)
	d.Set("updated_time", itemRg.UpdatedTime)
	d.Set("uniq_pools", itemRg.UniqPools)
	d.Set("vins", itemRg.VINS)
}

func flattenRgAudits(rgAudits ListAudits) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rgAudit := range rgAudits {
		temp := map[string]interface{}{
			"call":         rgAudit.Call,
			"responsetime": rgAudit.ResponseTime,
			"statuscode":   rgAudit.StatusCode,
			"timestamp":    rgAudit.Timestamp,
			"user":         rgAudit.User,
		}

		res = append(res, temp)
	}

	return res
}

func flattenRgList(rgl ListResourceGroups) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rg := range rgl {
		temp := map[string]interface{}{
			"account_acl":       flattenRgAcl(rg.ACL),
			"account_id":        rg.AccountID,
			"account_name":      rg.AccountName,
			"acl":               flattenRgAcl(rg.ACL),
			"created_by":        rg.CreatedBy,
			"created_time":      rg.CreatedTime,
			"def_net_id":        rg.DefNetID,
			"def_net_type":      rg.DefNetType,
			"deleted_by":        rg.DeletedBy,
			"deleted_time":      rg.DeletedTime,
			"desc":              rg.Description,
			"dirty":             rg.Dirty,
			"gid":               rg.GID,
			"guid":              rg.GUID,
			"rg_id":             rg.ID,
			"lock_status":       rg.LockStatus,
			"milestones":        rg.Milestones,
			"name":              rg.Name,
			"register_computes": rg.RegisterComputes,
			"resource_limits":   flattenRgResourceLimits(rg.ResourceLimits),
			"secret":            rg.Secret,
			"status":            rg.Status,
			"updated_by":        rg.UpdatedBy,
			"updated_time":      rg.UpdatedTime,
			"vins":              rg.VINS,
			"vms":               rg.Computes,
			"resource_types":    rg.ResTypes,
			"uniq_pools":        rg.UniqPools,
		}
		res = append(res, temp)
	}
	return res

}

func flattenRgAcl(rgAcls ListACL) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rgAcl := range rgAcls {
		temp := map[string]interface{}{
			"explicit":      rgAcl.Explicit,
			"guid":          rgAcl.GUID,
			"right":         rgAcl.Right,
			"status":        rgAcl.Status,
			"type":          rgAcl.Type,
			"user_group_id": rgAcl.UserGroupID,
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

func flattenRules(list ListRules) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rule := range list {
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

func flattenRgListComputes(lc ListComputes) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, compute := range lc {
		temp := map[string]interface{}{
			"account_id":         compute.AccountID,
			"account_name":       compute.AccountName,
			"affinity_label":     compute.AffinityLabel,
			"affinity_rules":     flattenRules(compute.AffinityRules),
			"affinity_weight":    compute.AffinityWeight,
			"antiaffinity_rules": flattenRules(compute.AntiAffinityRules),
			"cpus":               compute.CPUs,
			"created_by":         compute.CreatedBy,
			"created_time":       compute.CreatedTime,
			"deleted_by":         compute.DeletedBy,
			"deleted_time":       compute.DeletedTime,
			"id":                 compute.ID,
			"name":               compute.Name,
			"ram":                compute.RAM,
			"registered":         compute.Registered,
			"rg_name":            compute.RGName,
			"status":             compute.Status,
			"tech_status":        compute.TechStatus,
			"total_disks_size":   compute.TotalDisksSize,
			"updated_by":         compute.DeletedBy,
			"updated_time":       compute.DeletedTime,
			"user_managed":       compute.UserManaged,
			"vins_connected":     compute.VINSConnected,
		}

		res = append(res, temp)
	}

	return res
}

func flattenServerSettings(settings ServerSettings) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"inter":      settings.Inter,
		"guid":       settings.GUID,
		"down_inter": settings.DownInter,
		"rise":       settings.Rise,
		"fall":       settings.Fall,
		"slow_start": settings.SlowStart,
		"max_conn":   settings.MaxConn,
		"max_queue":  settings.MaxQueue,
		"weight":     settings.Weight,
	}
	res = append(res, temp)
	return res
}

func flattenListServers(list ListServers) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, serv := range list {
		temp := map[string]interface{}{
			"address":         serv.Address,
			"check":           serv.Check,
			"guid":            serv.GUID,
			"name":            serv.Name,
			"port":            serv.Port,
			"server_settings": flattenServerSettings(serv.ServerSettings),
		}
		res = append(res, temp)
	}

	return res
}

func flattenBackends(b ListBackends) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range b {
		temp := map[string]interface{}{
			"algorithm":               item.Algorithm,
			"guid":                    item.GUID,
			"name":                    item.Name,
			"server_default_settings": flattenServerSettings(item.ServerDefaultSettings),
			"servers":                 flattenListServers(item.Servers),
		}
		res = append(res, temp)
	}
	return res
}

func flattenBindings(list ListBindings) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, bind := range list {
		temp := map[string]interface{}{
			"address": bind.Address,
			"guid":    bind.GUID,
			"name":    bind.Name,
			"port":    bind.Port,
		}
		res = append(res, temp)
	}

	return res
}

func flattenFrontends(list ListFrontends) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, front := range list {
		temp := map[string]interface{}{
			"backend":  front.Backend,
			"bindings": flattenBindings(front.Bindings),
			"guid":     front.GUID,
			"name":     front.Name,
		}
		res = append(res, temp)
	}

	return res
}

func flattenNode(node RecordNode) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"backend_ip":  node.BackendIP,
		"compute_id":  node.ComputeID,
		"frontend_ip": node.FrontendIP,
		"guid":        node.GUID,
		"mgmt_ip":     node.MGMTIP,
		"network_id":  node.NetworkID,
	}
	res = append(res, temp)
	return res
}

func flattenRgListLb(listLb ListLB) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, lb := range listLb {
		temp := map[string]interface{}{
			"ha_mode":        lb.HAMode,
			"acl":            lb.ACL,
			"backends":       flattenBackends(lb.Backends),
			"created_by":     lb.CreatedBy,
			"created_time":   lb.CreatedTime,
			"deleted_by":     lb.DeletedBy,
			"deleted_time":   lb.DeletedTime,
			"desc":           lb.Description,
			"dp_api_user":    lb.DPAPIUser,
			"extnet_id":      lb.ExtNetID,
			"frontends":      flattenFrontends(lb.Frontends),
			"gid":            lb.GID,
			"guid":           lb.GUID,
			"id":             lb.ID,
			"image_id":       lb.ImageID,
			"milestones":     lb.Milestones,
			"name":           lb.Name,
			"primary_node":   flattenNode(lb.PrimaryNode),
			"rg_name":        lb.RGName,
			"secondary_node": flattenNode(lb.SecondaryNode),
			"status":         lb.Status,
			"tech_status":    lb.TechStatus,
			"updated_by":     lb.UpdatedBy,
			"updated_time":   lb.UpdatedTime,
			"vins_id":        lb.VINSID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenRgListPfw(listPfw ListPFW) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, pfw := range listPfw {
		temp := map[string]interface{}{
			"public_port_end":   pfw.PublicPortEnd,
			"public_port_start": pfw.PublicPortStart,
			"vm_id":             pfw.VMID,
			"vm_ip":             pfw.VMIP,
			"vm_name":           pfw.VMName,
			"vm_port":           pfw.VMPort,
			"vins_id":           pfw.VINSID,
			"vins_name":         pfw.VINSName,
		}
		res = append(res, temp)
	}

	return res
}

func flattenRgListVins(lv ListVINS) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, vins := range lv {
		temp := map[string]interface{}{
			"account_id":     vins.AccountID,
			"account_name":   vins.AccountName,
			"computes":       vins.Computes,
			"created_by":     vins.CreatedBy,
			"created_time":   vins.CreatedTime,
			"deleted_by":     vins.DeletedBy,
			"deleted_time":   vins.DeletedTime,
			"external_ip":    vins.ExternalIP,
			"id":             vins.ID,
			"name":           vins.Name,
			"network":        vins.Network,
			"pri_vnf_dev_id": vins.PriVNFDevID,
			"rg_name":        vins.RGName,
			"status":         vins.Status,
			"updated_by":     vins.UpdatedBy,
			"updated_time":   vins.UpdatedTime,
		}

		res = append(res, temp)
	}

	return res
}

func flattenRgAffinityGroupComputes(list ListAffinityGroupCompute) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)

	for _, item := range list {
		temp := map[string]interface{}{
			"compute_id":               item.ComputeID,
			"other_node":               item.OtherNode,
			"other_node_indirect":      item.OtherNodeIndirect,
			"other_node_indirect_soft": item.OtherNodeIndirectSoft,
			"other_node_soft":          item.OtherNodeSoft,
			"same_node":                item.SameNode,
			"same_node_soft":           item.SameNodeSoft,
		}
		res = append(res, temp)
	}

	return res
}

func flattenRgAffinityGroupsGet(list []uint64) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"items": list,
	}
	res = append(res, temp)

	return res
}

func flattenRgListGroups(list map[string][]uint64) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for groupKey, groupVal := range list {
		temp := map[string]interface{}{
			"label": groupKey,
			"ids":   groupVal,
		}
		res = append(res, temp)
	}

	return res
}

func flattenRgUsageResource(d *schema.ResourceData, usage Resource) {
	d.Set("cpu", usage.CPU)
	d.Set("disk_size", usage.DiskSize)
	d.Set("disk_size_max", usage.DiskSizeMax)
	d.Set("extips", usage.ExtIPs)
	d.Set("exttraffic", usage.ExtTraffic)
	d.Set("gpu", usage.GPU)
	d.Set("ram", usage.RAM)
	d.Set("seps", flattenRgSeps(usage.SEPs))
}
