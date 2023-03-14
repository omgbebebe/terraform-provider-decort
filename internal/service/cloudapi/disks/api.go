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

package disks

const (
	disksCreateAPI         = "/restmachine/cloudapi/disks/create"
	disksGetAPI            = "/restmachine/cloudapi/disks/get"
	disksListAPI           = "/restmachine/cloudapi/disks/list"
	disksResizeAPI         = "/restmachine/cloudapi/disks/resize2"
	disksRenameAPI         = "/restmachine/cloudapi/disks/rename"
	disksDeleteAPI         = "/restmachine/cloudapi/disks/delete"
	disksIOLimitAPI        = "/restmachine/cloudapi/disks/limitIO"
	disksRestoreAPI        = "/restmachine/cloudapi/disks/restore"
	disksListTypesAPI      = "/restmachine/cloudapi/disks/listTypes"
	disksListDeletedAPI    = "/restmachine/cloudapi/disks/listDeleted"
	disksListUnattachedAPI = "/restmachine/cloudapi/disks/listUnattached"

	disksSnapshotDeleteAPI   = "/restmachine/cloudapi/disks/snapshotDelete"
	disksSnapshotRollbackAPI = "/restmachine/cloudapi/disks/snapshotRollback"
	disksShareAPI            = "/restmachine/cloudapi/disks/share"
	disksUnshareAPI          = "/restmachine/cloudapi/disks/unshare"
)
