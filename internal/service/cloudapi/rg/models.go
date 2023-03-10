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

type ResourceLimits struct {
	CUC      float64 `json:"CU_C"`
	CUD      float64 `json:"CU_D"`
	CUI      float64 `json:"CU_I"`
	CUM      float64 `json:"CU_M"`
	CUNP     float64 `json:"CU_NP"`
	GpuUnits float64 `json:"gpu_units"`
}

type QuotaRecord struct { // this is how quota is reported by /api/.../rg/get
	Cpu        int     `json:"CU_C"`      // CPU count in pcs
	Ram        float64 `json:"CU_M"`      // RAM volume in MB, it is STILL reported as FLOAT
	Disk       int     `json:"CU_D"`      // Disk capacity in GB
	ExtIPs     int     `json:"CU_I"`      // Ext IPs count
	ExtTraffic int     `json:"CU_NP"`     // Ext network traffic
	GpuUnits   int     `json:"gpu_units"` // GPU count
}

// Main information about audit
type ItemAudit struct {
	// Call
	Call string `json:"call"`

	// Response time
	ResponseTime float64 `json:"responsetime"`

	// Status code
	StatusCode uint64 `json:"statuscode"`

	// Timestamp
	Timestamp float64 `json:"timestamp"`

	// User
	User string `json:"user"`
}

// List of audits
type ListAudits []ItemAudit

// Resources used
type Resource struct {
	// Number of cores
	CPU int64 `json:"cpu"`

	// Disk size
	DiskSize float64 `json:"disksize"`

	// Max disk size
	DiskSizeMax uint64 `json:"disksizemax"`

	// Number of External IPs
	ExtIPs int64 `json:"extips"`

	// External traffic
	ExtTraffic int64 `json:"exttraffic"`

	// Number of grafic cores
	GPU int64 `json:"gpu"`

	// Number of RAM
	RAM int64 `json:"ram"`

	// SEPs
	SEPs map[string]map[string]DiskUsage `json:"seps"`
}

// Disk usage
type DiskUsage struct {
	// Disk size
	DiskSize float64 `json:"disksize"`

	// Disk size max
	DiskSizeMax float64 `json:"disksizemax"`
}

// Information about resources
type Resources struct {
	// Current information about resources
	Current Resource `json:"Current"`

	// Reserved information about resources
	Reserved Resource `json:"Reserved"`
}

// Detailed information about resource group
type RecordResourceGroup struct {
	// Resources
	Resources Resources `json:"Resources"`

	// Account ID
	AccountID uint64 `json:"accountId"`

	// Account name
	AccountName string `json:"accountName"`

	// Access Control List
	ACL ListACL `json:"acl"`

	// Created by
	CreatedBy string `json:"createdBy"`

	// Created time
	CreatedTime uint64 `json:"createdTime"`

	// DefNetID
	DefNetID int64 `json:"def_net_id"`

	// DefNetType
	DefNetType string `json:"def_net_type"`

	// Deleted by
	DeletedBy string `json:"deletedBy"`

	// Deleted time
	DeletedTime uint64 `json:"deletedTime"`

	// Description
	Description string `json:"desc"`

	// Dirty
	Dirty bool `json:"dirty"`

	// Grid ID
	GID uint64 `json:"gid"`

	// GUID
	GUID uint64 `json:"guid"`

	// ID
	ID uint64 `json:"id"`

	// Lock status
	LockStatus string `json:"lockStatus"`

	// Milestones
	Milestones uint64 `json:"milestones"`

	// Name
	Name string `json:"name"`

	// RegisterComputes
	RegisterComputes bool `json:"registerComputes"`

	// Resource limits
	ResourceLimits ResourceLimits `json:"resourceLimits"`

	// Secret
	Secret string `json:"secret"`

	// Status
	Status string `json:"status"`

	// Updated by
	UpdatedBy string `json:"updatedBy"`

	// Updated time
	UpdatedTime uint64 `json:"updatedTime"`

	// List of VINS IDs
	VINS []uint64 `json:"vins"`

	// List of compute IDs
	Computes []uint64 `json:"vms"`

	// List of resource types
	ResTypes []string `json:"resourceTypes"`

	// UniqPools
	UniqPools []string `json:"uniqPools"`
}

// Main information about Access Control List
type ItemACL struct {
	// Explicit
	Explicit bool `json:"explicit"`

	// GUID
	GUID string `json:"guid"`

	// Right
	Right string `json:"right"`

	// Status
	Status string `json:"status"`

	// Type
	Type string `json:"type"`

	// User group ID
	UserGroupID string `json:"userGroupId"`
}

// List ACL
type ListACL []ItemACL

type ItemResourceGroup struct {
	//
	AccountACL ItemACL `json:"accountAcl"`

	// Account ID
	AccountID uint64 `json:"accountId"`

	// Account name
	AccountName string `json:"accountName"`

	// Access Control List
	ACL ListACL `json:"acl"`

	// Created by
	CreatedBy string `json:"createdBy"`

	// Created time
	CreatedTime uint64 `json:"createdTime"`

	// DefNetID
	DefNetID int64 `json:"def_net_id"`

	// DefNetType
	DefNetType string `json:"def_net_type"`

	// Deleted by
	DeletedBy string `json:"deletedBy"`

	// Deleted time
	DeletedTime uint64 `json:"deletedTime"`

	// Description
	Description string `json:"desc"`

	// Dirty
	Dirty bool `json:"dirty"`

	// Grid ID
	GID uint64 `json:"gid"`

	// GUID
	GUID uint64 `json:"guid"`

	// ID
	ID uint64 `json:"id"`

	// Lock status
	LockStatus string `json:"lockStatus"`

	// Milestones
	Milestones uint64 `json:"milestones"`

	// Name
	Name string `json:"name"`

	// RegisterComputes
	RegisterComputes bool `json:"registerComputes"`

	// Resource limits
	ResourceLimits ResourceLimits `json:"resourceLimits"`

	// Secret
	Secret string `json:"secret"`

	// Status
	Status string `json:"status"`

	// Updated by
	UpdatedBy string `json:"updatedBy"`

	// Updated time
	UpdatedTime uint64 `json:"updatedTime"`

	// List of VINS IDs
	VINS []uint64 `json:"vins"`

	// List of compute IDs
	Computes []uint64 `json:"vms"`

	// List of resource types
	ResTypes []string `json:"resourceTypes"`

	// UniqPools
	UniqPools []string `json:"uniqPools"`
}

// List of resource groups
type ListResourceGroups []ItemResourceGroup

// Main information about affinity rule
type ItemRule struct {
	// GUID
	GUID string `json:"guid"`

	// Key
	Key string `json:"key"`

	// Mode
	Mode string `json:"mode"`

	// Policy
	Policy string `json:"policy"`

	// Topology
	Topology string `json:"topology"`

	// Value
	Value string `json:"value"`
}

// List rules
type ListRules []ItemRule

// Main information about compute
type ItemCompute struct {
	// Account ID
	AccountID uint64 `json:"accountId"`

	// Account name
	AccountName string `json:"accountName"`

	// Affinity label
	AffinityLabel string `json:"affinityLabel"`

	// List affinity rules
	AffinityRules ListRules `json:"affinityRules"`

	// Affinity weight
	AffinityWeight uint64 `json:"affinityWeight"`

	// Anti affinity rules
	AntiAffinityRules ListRules `json:"antiAffinityRules"`

	// Number of CPU
	CPUs uint64 `json:"cpus"`

	// Created by
	CreatedBy string `json:"createdBy"`

	// Created time
	CreatedTime uint64 `json:"createdTime"`

	// Deleted by
	DeletedBy string `json:"deletedBy"`

	// Deleted time
	DeletedTime uint64 `json:"deletedTime"`

	// ID
	ID uint64 `json:"id"`

	// Name
	Name string `json:"name"`

	// Number of RAM
	RAM uint64 `json:"ram"`

	// Registered
	Registered bool `json:"registered"`

	// Resource group ID
	RGID uint64 `json:"rgId"`

	// Resource group name
	RGName string `json:"rgName"`

	// Status
	Status string `json:"status"`

	// Tech status
	TechStatus string `json:"techStatus"`

	// Total disks size
	TotalDisksSize uint64 `json:"totalDisksSize"`

	// Updated by
	UpdatedBy string `json:"updatedBy"`

	// Updated time
	UpdatedTime uint64 `json:"updatedTime"`

	// User managed
	UserManaged bool `json:"userManaged"`

	// VINS connected
	VINSConnected uint64 `json:"vinsConnected"`
}

// List computes
type ListComputes []ItemCompute

// Main information about port forward
type ItemPFW struct {
	// Public port end
	PublicPortEnd uint64 `json:"Public Port End"`

	// Public port start
	PublicPortStart uint64 `json:"Public Port Start"`

	// Virtual machine ID
	VMID uint64 `json:"VM ID"`

	// Virtual machine IP
	VMIP string `json:"VM IP"`

	// Virtual machine name
	VMName string `json:"VM Name"`

	// Virtual machine port
	VMPort uint64 `json:"VM Port"`

	// VINS ID
	VINSID uint64 `json:"ViNS ID"`

	// VINS name
	VINSName string `json:"ViNS Name"`
}

// List PFWs
type ListPFW []ItemPFW

// Main information about VINS
type ItemVINS struct {
	// Account ID
	AccountID uint64 `json:"accountId"`

	// Account name
	AccountName string `json:"accountName"`

	// Computes
	Computes uint64 `json:"computes"`

	// Created by
	CreatedBy string `json:"createdBy"`

	// Created time
	CreatedTime uint64 `json:"createdTime"`

	// Deleted by
	DeletedBy string `json:"deletedBy"`

	// Deleted time
	DeletedTime uint64 `json:"deletedTime"`

	// External IP
	ExternalIP string `json:"externalIP"`

	// ID
	ID uint64 `json:"id"`

	// Name
	Name string `json:"name"`

	// Network
	Network string `json:"network"`

	// PriVNFDev ID
	PriVNFDevID uint64 `json:"priVnfDevId"`

	// Resource group ID
	RGID uint64 `json:"rgId"`

	// Resource group name
	RGName string `json:"rgName"`

	// Status
	Status string `json:"status"`

	// Updated by
	UpdatedBy string `json:"updatedBy"`

	// Updated time
	UpdatedTime uint64 `json:"updatedTime"`
}

// List VINSes
type ListVINS []ItemVINS

// Server settings
type ServerSettings struct {
	// Inter
	Inter uint64 `json:"inter"`

	// GUID
	GUID string `json:"guid"`

	// Down inter
	DownInter uint64 `json:"downinter"`

	// Rise
	Rise uint64 `json:"rise"`

	// Fall
	Fall uint64 `json:"fall"`

	// Slow start
	SlowStart uint64 `json:"slowstart"`

	// Max connections
	MaxConn uint64 `json:"maxconn"`

	// Max queue
	MaxQueue uint64 `json:"maxqueue"`

	// Weight
	Weight uint64 `json:"weight"`
}

// Main information about server
type ItemServer struct {
	// Address
	Address string `json:"address"`

	// Check
	Check string `json:"check"`

	// GUID
	GUID string `json:"guid"`

	// Name
	Name string `json:"name"`

	// Port
	Port uint64 `json:"port"`

	// Server settings
	ServerSettings ServerSettings `json:"serverSettings"`
}

// List of servers
type ListServers []ItemServer

// Main information about backend
type ItemBackend struct {
	// Algorithm
	Algorithm string `json:"algorithm"`

	// GUID
	GUID string `json:"guid"`

	// Name
	Name string `json:"name"`

	// Server settings
	ServerDefaultSettings ServerSettings `json:"serverDefaultSettings"`

	// List of servers
	Servers ListServers `json:"servers"`
}

// List of backends
type ListBackends []ItemBackend

// Main information of binding
type ItemBinding struct {
	// Address
	Address string `json:"address"`

	// GUID
	GUID string `json:"guid"`

	// Name
	Name string `json:"name"`

	// Port
	Port uint64 `json:"port"`
}

// List of bindings
type ListBindings []ItemBinding

// Main information about frontend
type ItemFrontend struct {
	// Backend
	Backend string `json:"backend"`

	// List of bindings
	Bindings ListBindings `json:"bindings"`

	// GUID
	GUID string `json:"guid"`

	// Name
	Name string `json:"name"`
}

// List of frontends
type ListFrontends []ItemFrontend

// Main information about node
type RecordNode struct {
	// Backend IP
	BackendIP string `json:"backendIp"`

	// Compute ID
	ComputeID uint64 `json:"computeId"`

	// Frontend IP
	FrontendIP string `json:"frontendIp"`

	// GUID
	GUID string `json:"guid"`

	// MGMT IP
	MGMTIP string `json:"mgmtIp"`

	// Network ID
	NetworkID uint64 `json:"networkId"`
}

// Main information about load balancer
type ItemLB struct {
	// HAMode
	HAMode bool `json:"HAmode"`

	// List ACL
	ACL ListACL `json:"acl"`

	// List backends
	Backends ListBackends `json:"backends"`

	// Created by
	CreatedBy string `json:"createdBy"`

	// Created time
	CreatedTime uint64 `json:"createdTime"`

	// Deleted by
	DeletedBy string `json:"deletedBy"`

	// Deleted time
	DeletedTime uint64 `json:"deletedTime"`

	// Description
	Description string `json:"desc"`

	// DPAPI user
	DPAPIUser string `json:"dpApiUser"`

	// External network ID
	ExtNetID uint64 `json:"extnetId"`

	// List of frontends
	Frontends ListFrontends `json:"frontends"`

	// Grid ID
	GID uint64 `json:"gid"`

	// GUID
	GUID uint64 `json:"guid"`

	// ID
	ID uint64 `json:"id"`

	// Image ID
	ImageID uint64 `json:"imageId"`

	// Milestones
	Milestones uint64 `json:"milestones"`

	// Name
	Name string `json:"name"`

	// Primary node
	PrimaryNode RecordNode `json:"primaryNode"`

	// Resource group ID
	RGID uint64 `json:"rgId"`

	// Resource group name
	RGName string `json:"rgName"`

	// Secondary node
	SecondaryNode RecordNode `json:"secondaryNode"`

	// Status
	Status string `json:"status"`

	// Tech status
	TechStatus string `json:"techStatus"`

	// Updated by
	UpdatedBy string `json:"updatedBy"`

	// Updated time
	UpdatedTime uint64 `json:"updatedTime"`

	// VINS ID
	VINSID uint64 `json:"vinsId"`
}

// List load balancers
type ListLB []ItemLB

// Main information about affinity group
type ItemAffinityGroupCompute struct {
	// Compute ID
	ComputeID uint64 `json:"computeId"`

	// Other node
	OtherNode []uint64 `json:"otherNode"`

	// Other node indirect
	OtherNodeIndirect []uint64 `json:"otherNodeIndirect"`

	// Other node indirect soft
	OtherNodeIndirectSoft []uint64 `json:"otherNodeIndirectSoft"`

	// Other node soft
	OtherNodeSoft []uint64 `json:"otherNodeSoft"`

	// Same node
	SameNode []uint64 `json:"sameNode"`

	// Same node soft
	SameNodeSoft []uint64 `json:"sameNodeSoft"`
}

// List of affinity groups
type ListAffinityGroupCompute []ItemAffinityGroupCompute
