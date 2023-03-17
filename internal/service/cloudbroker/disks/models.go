/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>

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

package disks

type Disk struct {
	Acl                 map[string]interface{} `json:"acl"`
	AccountID           int                    `json:"accountId"`
	AccountName         string                 `json:"accountName"`
	BootPartition       int                    `json:"bootPartition"`
	CreatedTime         uint64                 `json:"creationTime"`
	ComputeID           int                    `json:"computeId"`
	ComputeName         string                 `json:"computeName"`
	DeletedTime         uint64                 `json:"deletionTime"`
	DeviceName          string                 `json:"devicename"`
	Desc                string                 `json:"desc"`
	DestructionTime     uint64                 `json:"destructionTime"`
	DiskPath            string                 `json:"diskPath"`
	GridID              int                    `json:"gid"`
	GUID                int                    `json:"guid"`
	ID                  uint                   `json:"id"`
	ImageID             int                    `json:"imageId"`
	Images              []int                  `json:"images"`
	IOTune              IOTune                 `json:"iotune"`
	IQN                 string                 `json:"iqn"`
	Login               string                 `json:"login"`
	Name                string                 `json:"name"`
	MachineId           int                    `json:"machineId"`
	MachineName         string                 `json:"machineName"`
	Milestones          uint64                 `json:"milestones"`
	Order               int                    `json:"order"`
	Params              string                 `json:"params"`
	Passwd              string                 `json:"passwd"`
	ParentId            int                    `json:"parentId"`
	PciSlot             int                    `json:"pciSlot"`
	Pool                string                 `json:"pool"`
	PurgeTime           uint64                 `json:"purgeTime"`
	PurgeAttempts       uint64                 `json:"purgeAttempts"`
	RealityDeviceNumber int                    `json:"realityDeviceNumber"`
	ReferenceId         string                 `json:"referenceId"`
	ResID               string                 `json:"resId"`
	ResName             string                 `json:"resName"`
	Role                string                 `json:"role"`
	SepType             string                 `json:"sepType"`
	SepID               int                    `json:"sepId"` // NOTE: absent from compute/get output
	SizeMax             int                    `json:"sizeMax"`
	SizeUsed            float64                `json:"sizeUsed"` // sum over all snapshots of this disk to report total consumed space
	Snapshots           []Snapshot             `json:"snapshots"`
	Status              string                 `json:"status"`
	TechStatus          string                 `json:"techStatus"`
	Type                string                 `json:"type"`
	UpdateBy            uint64                 `json:"updateBy"`
	VMID                int                    `json:"vmid"`
}

type Snapshot struct {
	Guid        string `json:"guid"`
	Label       string `json:"label"`
	ResId       string `json:"resId"`
	SnapSetGuid string `json:"snapSetGuid"`
	SnapSetTime uint64 `json:"snapSetTime"`
	TimeStamp   uint64 `json:"timestamp"`
}

type SnapshotList []Snapshot

type DisksList []Disk

type IOTune struct {
	ReadBytesSec     int `json:"read_bytes_sec"`
	ReadBytesSecMax  int `json:"read_bytes_sec_max"`
	ReadIopsSec      int `json:"read_iops_sec"`
	ReadIopsSecMax   int `json:"read_iops_sec_max"`
	SizeIopsSec      int `json:"size_iops_sec"`
	TotalBytesSec    int `json:"total_bytes_sec"`
	TotalBytesSecMax int `json:"total_bytes_sec_max"`
	TotalIopsSec     int `json:"total_iops_sec"`
	TotalIopsSecMax  int `json:"total_iops_sec_max"`
	WriteBytesSec    int `json:"write_bytes_sec"`
	WriteBytesSecMax int `json:"write_bytes_sec_max"`
	WriteIopsSec     int `json:"write_iops_sec"`
	WriteIopsSecMax  int `json:"write_iops_sec_max"`
}
