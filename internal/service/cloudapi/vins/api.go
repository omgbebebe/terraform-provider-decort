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

Source code: https://repos.digitalenergy.online/BASIS/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://repos.digitalenergy.online/BASIS/terraform-provider-decort/wiki
*/

package vins

const (
	VinsAuditsAPI           = "/restmachine/cloudapi/vins/audits"
	VinsCreateInAccountAPI  = "/restmachine/cloudapi/vins/createInAccount"
	VinsCreateInRgAPI       = "/restmachine/cloudapi/vins/createInRG"
	VinsDeleteAPI           = "/restmachine/cloudapi/vins/delete"
	VinsDisableAPI          = "/restmachine/cloudapi/vins/disable"
	VinsEnableAPI           = "/restmachine/cloudapi/vins/enable"
	VinsExtNetConnectAPI    = "/restmachine/cloudapi/vins/extNetConnect"
	VinsExtNetDisconnectAPI = "/restmachine/cloudapi/vins/extNetDisconnect"
	VinsExtNetListAPI       = "/restmachine/cloudapi/vins/extNetList"
	VinsGetAPI              = "/restmachine/cloudapi/vins/get"
	VinsIpListAPI           = "/restmachine/cloudapi/vins/ipList"
	VinsIpReleaseAPI        = "/restmachine/cloudapi/vins/ipRelease"
	VinsIpReserveAPI        = "/restmachine/cloudapi/vins/ipReserve"
	VinsListAPI             = "/restmachine/cloudapi/vins/list"
	VinsListDeletedAPI      = "/restmachine/cloudapi/vins/listDeleted"
	VinsNatRuleAddAPI       = "/restmachine/cloudapi/vins/natRuleAdd"
	VinsNatRuleDelAPI       = "/restmachine/cloudapi/vins/natRuleDel"
	VinsNatRuleListAPI      = "/restmachine/cloudapi/vins/natRuleList"
	VinsRestoreAPI          = "/restmachine/cloudapi/vins/restore"
	VinsSearchAPI           = "/restmachine/cloudapi/vins/search"
	VinsVnfdevRedeployAPI   = "/restmachine/cloudapi/vins/vnfdevRedeploy"
	VinsVnfdevRestartAPI    = "/restmachine/cloudapi/vins/vnfdevRestart"
)
