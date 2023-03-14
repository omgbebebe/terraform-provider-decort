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

package cloudbroker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/account"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/disks"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/image"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/k8s"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/kvmvm"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/pcidevice"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/pfw"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/rg"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/sep"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/snapshot"
	"repos.digitalenergy.online/BASIS/terraform-provider-decort/internal/service/cloudbroker/vins"
)

func NewRersourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"decort_account":       account.ResourceAccount(),
		"decort_disk":          disks.ResourceDisk(),
		"decort_image":         image.ResourceImage(),
		"decort_virtual_image": image.ResourceVirtualImage(),
		"decort_cdrom_image":   image.ResourceCDROMImage(),
		"decort_delete_images": image.ResourceDeleteImages(),
		"decort_pcidevice":     pcidevice.ResourcePcidevice(),
		"decort_sep":           sep.ResourceSep(),
		"decort_sep_config":    sep.ResourceSepConfig(),
		"decort_resgroup":      rg.ResourceResgroup(),
		"decort_kvmvm":         kvmvm.ResourceCompute(),
		"decort_vins":          vins.ResourceVins(),
		"decort_pfw":           pfw.ResourcePfw(),
		"decort_k8s":           k8s.ResourceK8s(),
		"decort_k8s_wg":        k8s.ResourceK8sWg(),
		"decort_snapshot":      snapshot.ResourceSnapshot(),
	}
}
