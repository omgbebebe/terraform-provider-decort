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

package kvmvm

const (
	KvmX86CreateAPI                  = "/restmachine/cloudapi/kvmx86/create"
	KvmPPCCreateAPI                  = "/restmachine/cloudapi/kvmppc/create"
	ComputeGetAPI                    = "/restmachine/cloudapi/compute/get"
	RgListComputesAPI                = "/restmachine/cloudapi/rg/listComputes"
	ComputeNetAttachAPI              = "/restmachine/cloudapi/compute/netAttach"
	ComputeNetDetachAPI              = "/restmachine/cloudapi/compute/netDetach"
	ComputeDiskAttachAPI             = "/restmachine/cloudapi/compute/diskAttach"
	ComputeDiskDetachAPI             = "/restmachine/cloudapi/compute/diskDetach"
	ComputeStartAPI                  = "/restmachine/cloudapi/compute/start"
	ComputeStopAPI                   = "/restmachine/cloudapi/compute/stop"
	ComputeResizeAPI                 = "/restmachine/cloudapi/compute/resize"
	DisksResizeAPI                   = "/restmachine/cloudapi/disks/resize2"
	ComputeDeleteAPI                 = "/restmachine/cloudapi/compute/delete"
	ComputeUpdateAPI                 = "/restmachine/cloudapi/compute/update"
	ComputeDiskAddAPI                = "/restmachine/cloudapi/compute/diskAdd"
	ComputeDiskDeleteAPI             = "/restmachine/cloudapi/compute/diskDel"
	ComputeRestoreAPI                = "/restmachine/cloudapi/compute/restore"
	ComputeEnableAPI                 = "/restmachine/cloudapi/compute/enable"
	ComputeDisableAPI                = "/restmachine/cloudapi/compute/disable"
	ComputeAffinityLabelSetAPI       = "/restmachine/cloudapi/compute/affinityLabelSet"
	ComputeAffinityLabelRemoveAPI    = "/restmachine/cloudapi/compute/affinityLabelRemove"
	ComputeAffinityRuleAddAPI        = "/restmachine/cloudapi/compute/affinityRuleAdd"
	ComputeAffinityRuleRemoveAPI     = "/restmachine/cloudapi/compute/affinityRuleRemove"
	ComputeAffinityRulesClearAPI     = "/restmachine/cloudapi/compute/affinityRulesClear"
	ComputeAntiAffinityRuleAddAPI    = "/restmachine/cloudapi/compute/antiAffinityRuleAdd"
	ComputeAntiAffinityRuleRemoveAPI = "/restmachine/cloudapi/compute/antiAffinityRuleRemove"
	ComputeAntiAffinityRulesClearAPI = "/restmachine/cloudapi/compute/antiAffinityRulesClear"
	ComputeListAPI                   = "/restmachine/cloudapi/compute/list"
	ComputeAuditsAPI                 = "/restmachine/cloudapi/compute/audits"
	ComputeGetAuditsAPI              = "/restmachine/cloudapi/compute/getAudits"
	ComputeGetConsoleUrlAPI          = "/restmachine/cloudapi/compute/getConsoleUrl"
	ComputeGetLogAPI                 = "/restmachine/cloudapi/compute/getLog"
	ComputePfwListAPI                = "/restmachine/cloudapi/compute/pfwList"
	ComputeUserListAPI               = "/restmachine/cloudapi/compute/userList"
	ComputeTagAddAPI                 = "/restmachine/cloudapi/compute/tagAdd"
	ComputeTagRemoveAPI              = "/restmachine/cloudapi/compute/tagRemove"
	ComputePinToStackAPI             = "/restmachine/cloudapi/compute/pinToStack"
	ComputeUnpinFromStackAPI         = "/restmachine/cloudapi/compute/unpinFromStack"
	ComputePfwAddAPI                 = "/restmachine/cloudapi/compute/pfwAdd"
	ComputePfwDelAPI                 = "/restmachine/cloudapi/compute/pfwDel"
	ComputeUserGrantAPI              = "/restmachine/cloudapi/compute/userGrant"
	ComputeUserRevokeAPI             = "/restmachine/cloudapi/compute/userRevoke"
	ComputeSnapshotCreateAPI         = "/restmachine/cloudapi/compute/snapshotCreate"
	ComputeSnapshotDeleteAPI         = "/restmachine/cloudapi/compute/snapshotCreate"
	ComputeSnapshotUsageAPI          = "/restmachine/cloudapi/compute/snapshotUsage"
	ComputeSnapshotRollbackAPI       = "/restmachine/cloudapi/compute/snapshotRollback"
	ComputePauseAPI                  = "/restmachine/cloudapi/compute/pause"
	ComputeResumeAPI                 = "/restmachine/cloudapi/compute/resume"
	ComputeCdInsertAPI               = "/restmachine/cloudapi/compute/cdInsert"
	ComputeCdEjectAPI                = "/restmachine/cloudapi/compute/cdEject"
	ComputeResetAPI                  = "/restmachine/cloudapi/compute/reset"
	ComputeRedeployAPI               = "/restmachine/cloudapi/compute/redeploy"
)
