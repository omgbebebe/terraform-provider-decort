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

package k8s

const K8sCreateAPI = "/restmachine/cloudbroker/k8s/create"
const K8sGetAPI = "/restmachine/cloudbroker/k8s/get"
const K8sUpdateAPI = "/restmachine/cloudbroker/k8s/update"
const K8sDeleteAPI = "/restmachine/cloudbroker/k8s/delete"

const K8sWgCreateAPI = "/restmachine/cloudbroker/k8s/workersGroupAdd"
const K8sWgDeleteAPI = "/restmachine/cloudbroker/k8s/workersGroupDelete"

const K8sWorkerAddAPI = "/restmachine/cloudbroker/k8s/workerAdd"
const K8sWorkerDeleteAPI = "/restmachine/cloudbroker/k8s/deleteWorkerFromGroup"

const K8sGetConfigAPI = "/restmachine/cloudbroker/k8s/getConfig"

const LbGetAPI = "/restmachine/cloudbroker/lb/get"

const AsyncTaskGetAPI = "/restmachine/cloudbroker/tasks/get"
