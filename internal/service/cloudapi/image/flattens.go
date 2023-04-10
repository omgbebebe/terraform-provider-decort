package image

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func flattenHistory(history []History) []map[string]interface{} {
	temp := make([]map[string]interface{}, 0)
	for _, item := range history {
		t := map[string]interface{}{
			"id":        item.Id,
			"guid":      item.Guid,
			"timestamp": item.Timestamp,
		}

		temp = append(temp, t)
	}
	return temp
}

func flattenImage(d *schema.ResourceData, img *ImageExtend) {
	d.Set("unc_path", img.UNCPath)
	d.Set("ckey", img.CKey)
	d.Set("account_id", img.AccountId)
	d.Set("acl", img.Acl)
	d.Set("architecture", img.Architecture)
	d.Set("boot_type", img.BootType)
	d.Set("bootable", img.Bootable)
	d.Set("compute_ci_id", img.ComputeCiId)
	d.Set("deleted_time", img.DeletedTime)
	d.Set("desc", img.Description)
	d.Set("drivers", img.Drivers)
	d.Set("enabled", img.Enabled)
	d.Set("gid", img.GridId)
	d.Set("guid", img.GUID)
	d.Set("history", flattenHistory(img.History))
	d.Set("hot_resize", img.HotResize)
	d.Set("image_id", img.Id)
	d.Set("last_modified", img.LastModified)
	d.Set("link_to", img.LinkTo)
	d.Set("milestones", img.Milestones)
	d.Set("image_name", img.Name)
	d.Set("password", img.Password)
	d.Set("pool_name", img.Pool)
	d.Set("provider_name", img.ProviderName)
	d.Set("purge_attempts", img.PurgeAttempts)
	d.Set("present_to", img.PresentTo)
	d.Set("res_id", img.ResId)
	d.Set("rescuecd", img.RescueCD)
	d.Set("sep_id", img.SepId)
	d.Set("shared_with", img.SharedWith)
	d.Set("size", img.Size)
	d.Set("status", img.Status)
	d.Set("tech_status", img.TechStatus)
	d.Set("type", img.Type)
	d.Set("username", img.Username)
	d.Set("version", img.Version)
}
