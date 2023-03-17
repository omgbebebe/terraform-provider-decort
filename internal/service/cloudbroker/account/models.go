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

package account

type AccountAclRecord struct {
	IsExplicit bool   `json:"explicit"`
	Guid       string `json:"guid"`
	Rights     string `json:"right"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	UgroupID   string `json:"userGroupId"`
}

type ResourceLimits struct {
	CUC      float64 `json:"CU_C"`
	CUD      float64 `json:"CU_D"`
	CUI      float64 `json:"CU_I"`
	CUM      float64 `json:"CU_M"`
	CUNP     float64 `json:"CU_NP"`
	GpuUnits float64 `json:"gpu_units"`
}

type Account struct {
	DCLocation        string             `json:"DCLocation"`
	CKey              string             `jspn:"_ckey"`
	Meta              []interface{}      `json:"_meta"`
	Acl               []AccountAclRecord `json:"acl"`
	Company           string             `json:"company"`
	CompanyUrl        string             `json:"companyurl"`
	CreatedBy         string             `jspn:"createdBy"`
	CreatedTime       int                `json:"createdTime"`
	DeactiovationTime float64            `json:"deactivationTime"`
	DeletedBy         string             `json:"deletedBy"`
	DeletedTime       int                `json:"deletedTime"`
	DisplayName       string             `json:"displayname"`
	GUID              int                `json:"guid"`
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	ResourceLimits    ResourceLimits     `json:"resourceLimits"`
	SendAccessEmails  bool               `json:"sendAccessEmails"`
	ServiceAccount    bool               `json:"serviceAccount"`
	Status            string             `json:"status"`
	UpdatedTime       int                `json:"updatedTime"`
	Version           int                `json:"version"`
	Vins              []int              `json:"vins"`
}

type AccountList []Account

type Resource struct {
	CPU        int `json:"cpu"`
	Disksize   int `json:"disksize"`
	Extips     int `json:"extips"`
	Exttraffic int `json:"exttraffic"`
	GPU        int `json:"gpu"`
	RAM        int `json:"ram"`
}

type Resources struct {
	Current  Resource `json:"Current"`
	Reserved Resource `json:"Reserved"`
}

type AccountWithResources struct {
	Account
	Resources Resources `json:"Resources"`
}

type AccountCompute struct {
	AccountId      int    `json:"accountId"`
	AccountName    string `json:"accountName"`
	CPUs           int    `json:"cpus"`
	CreatedBy      string `json:"createdBy"`
	CreatedTime    int    `json:"createdTime"`
	DeletedBy      string `json:"deletedBy"`
	DeletedTime    int    `json:"deletedTime"`
	ComputeId      int    `json:"id"`
	ComputeName    string `json:"name"`
	RAM            int    `json:"ram"`
	Registered     bool   `json:"registered"`
	RgId           int    `json:"rgId"`
	RgName         string `json:"rgName"`
	Status         string `json:"status"`
	TechStatus     string `json:"techStatus"`
	TotalDisksSize int    `json:"totalDisksSize"`
	UpdatedBy      string `json:"updatedBy"`
	UpdatedTime    int    `json:"updatedTime"`
	UserManaged    bool   `json:"userManaged"`
	VinsConnected  int    `json:"vinsConnected"`
}

type AccountComputesList []AccountCompute

type AccountDisk struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Pool    string `json:"pool"`
	SepId   int    `json:"sepId"`
	SizeMax int    `json:"sizeMax"`
	Type    string `json:"type"`
}

type AccountDisksList []AccountDisk

type AccountVin struct {
	AccountId   int    `json:"accountId"`
	AccountName string `json:"accountName"`
	Computes    int    `json:"computes"`
	CreatedBy   string `json:"createdBy"`
	CreatedTime int    `json:"createdTime"`
	DeletedBy   string `json:"deletedBy"`
	DeletedTime int    `json:"deletedTime"`
	ExternalIP  string `json:"externalIP"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Network     string `json:"network"`
	PriVnfDevId int    `json:"priVnfDevId"`
	RgId        int    `json:"rgId"`
	RgName      string `json:"rgName"`
	Status      string `json:"status"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedTime int    `json:"updatedTime"`
}

type AccountVinsList []AccountVin

type AccountAudit struct {
	Call         string  `json:"call"`
	ResponseTime float64 `json:"responsetime"`
	StatusCode   int     `json:"statuscode"`
	Timestamp    float64 `json:"timestamp"`
	User         string  `json:"user"`
}

type AccountAuditsList []AccountAudit

type AccountRGComputes struct {
	Started int `json:"Started"`
	Stopped int `json:"Stopped"`
}

type AccountRGResources struct {
	Consumed Resource `json:"Consumed"`
	Limits   Resource `json:"Limits"`
	Reserved Resource `json:"Reserved"`
}

type AccountRG struct {
	Computes    AccountRGComputes  `json:"Computes"`
	Resources   AccountRGResources `json:"Resources"`
	CreatedBy   string             `json:"createdBy"`
	CreatedTime int                `json:"createdTime"`
	DeletedBy   string             `json:"deletedBy"`
	DeletedTime int                `json:"deletedTime"`
	RGID        int                `json:"id"`
	Milestones  int                `json:"milestones"`
	RGName      string             `json:"name"`
	Status      string             `json:"status"`
	UpdatedBy   string             `json:"updatedBy"`
	UpdatedTime int                `json:"updatedTime"`
	Vinses      int                `json:"vinses"`
}

type AccountRGList []AccountRG

type AccountFlipGroup struct {
	AccountId   int    `json:"accountId"`
	ClientType  string `json:"clientType"`
	ConnType    string `json:"connType"`
	CreatedBy   string `json:"createdBy"`
	CreatedTime int    `json:"createdTime"`
	DefaultGW   string `json:"defaultGW"`
	DeletedBy   string `json:"deletedBy"`
	DeletedTime int    `json:"deletedTime"`
	Desc        string `json:"desc"`
	GID         int    `json:"gid"`
	GUID        int    `json:"guid"`
	ID          int    `json:"id"`
	IP          string `json:"ip"`
	Milestones  int    `json:"milestones"`
	Name        string `json:"name"`
	NetID       int    `json:"netId"`
	NetType     string `json:"netType"`
	NetMask     int    `json:"netmask"`
	Status      string `json:"status"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedTime int    `json:"updatedTime"`
}

type AccountFlipGroupsList []AccountFlipGroup
