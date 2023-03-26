## Version 3.5.2

### Features
- Add new datasource decort_kvmvm_snapshot_usage
- Add the ability to change the size in the 'disks' block in the decort_kvmvm resource. Now when you change 'size' field in the block, the disk size on the platform will also be changed

## Bug Fix
- rule "release" in Makefile don't create the necessary archives
- field "register_computes" in resource decort_resgroup is not used when creating the resource
- removed unused optional fields in datasources
