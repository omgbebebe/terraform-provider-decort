### Version 3.5.0

## Features

#### Resgroup
- Add data source rg_affinity_group_computes
- Add data source rg_affinity_groups_get
- Add data source rg_affinity_groups_list
- Add data source rg_audits
- Add data source rg_list
- Add data source rg_list_computes
- Add data source rg_list_deleted
- Add data source rg_list_lb
- Add data source rg_list_pfw
- Add data source rg_list_vins
- Add data source rg_usage
- Update data source rg
- Update block 'qouta' to change resource limits
- Add block 'access' to access/revoke rights for rg
- Add block 'def_net' to set default network in rg
- Add field 'enable' to disable/enable rg
- Add processing of input parameters (account_id, gid, ext_net_id) when creating and updating a resource

#### Kvmvm
- Update data source decort_kvmvm
- Add data source decort_kvmvm_list
- Add data source decort_kvmvm_audits
- Add data source decort_kvmvm_get_audits
- Add data source decort_kvmvm_get_console_url
- Add data source decort_kvmvm_get_log
- Add data source decort_kvmvm_pfw_list
- Add data source decort_kvmvm_user_list
- Update block 'disks' in the resource decort_kvmvm
- Add block 'tags' to add/delete tags
- Add block 'port_forwarding' to add/delete pfws
- Add block 'user_access' to access/revoke user rights for comptue
- Add block 'snapshot' to create/delete snapshots
- Add block 'rollback' to rollback in snapshot
- Add block 'cd' to insert/Eject cdROM disks
- Add field 'pin_to_stack' to pin compute to stack
- Add field 'pause' to pause/resume compute
- Add field 'reset' to reset compute
- Add the ability to redeploy the compute when changing the image_id
- Add field 'data_disks' to redeploy compute
- Add field 'auto_start' to redeploy compute
- Add field 'force_stop' to redeploy compute
- Add warnings in Create resource decort_kvmvm
- Add processing of input parameters (rg_id, image_id and all vins_id's in blocks 'network') when creating and updating a resource

## Bug Fix

- When deleting the 'quote' block, the limits are not set to the default value
- Block 'disks' in resource decort_kvmvm breaks the state
- Import decort_resgroup resource breaks the state
- Import decort_kvmvm resource breaks the state
- If the boot_disk_size is not specified at creation, further changing it leads to an error
