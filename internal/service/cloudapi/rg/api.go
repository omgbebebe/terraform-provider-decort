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

const (
	ResgroupCreateAPI          = "/restmachine/cloudapi/rg/create"
	ResgroupUpdateAPI          = "/restmachine/cloudapi/rg/update"
	ResgroupListAPI            = "/restmachine/cloudapi/rg/list"
	ResgroupListDeletedAPI     = "/restmachine/cloudapi/rg/listDeleted"
	ResgroupListPfwAPI         = "/restmachine/cloudapi/rg/listPFW"
	ResgroupGetAPI             = "/restmachine/cloudapi/rg/get"
	ResgroupListVinsAPI        = "/restmachine/cloudapi/rg/listVins"
	ResgroupListLbAPI          = "/restmachine/cloudapi/rg/listLb"
	ResgroupDeleteAPI          = "/restmachine/cloudapi/rg/delete"
	RgListComputesAPI          = "/restmachine/cloudapi/rg/listComputes"
	RgAffinityGroupComputesAPI = "/restmachine/cloudapi/rg/affinityGroupComputes"
	RgAffinityGroupsGetAPI     = "/restmachine/cloudapi/rg/affinityGroupsGet"
	RgAffinityGroupsListAPI    = "/restmachine/cloudapi/rg/affinityGroupsList"
	RgAuditsAPI                = "/restmachine/cloudapi/rg/audits"
	RgEnableAPI                = "/restmachine/cloudapi/rg/enable"
	RgDisableAPI               = "/restmachine/cloudapi/rg/disable"
	ResgroupUsageAPI           = "/restmachine/cloudapi/rg/usage"
	RgAccessGrantAPI           = "/restmachine/cloudapi/rg/accessGrant"
	RgAccessRevokeAPI          = "/restmachine/cloudapi/rg/accessRevoke"
	RgSetDefNetAPI             = "/restmachine/cloudapi/rg/setDefNet"
	RgRestoreAPI               = "/restmachine/cloudapi/rg/restore"
)
