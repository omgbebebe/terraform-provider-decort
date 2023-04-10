## Version 3.6.0

### Features
- Added validation for required fields in the following resources:
	- disks
	- lb
	- image
	- bservice
	- vins
	- k8s
	- k8s_wg
	- lb_frontend
	- lb_frontend_bind
	- lb_backed
	- lb_backend_server
- Added status handlers in create/update functions (where present)

### Bug Fixes
- Fixed state inconsistency in the following resources/data sources:
	- data_source_account
	- resource_k8s
	- resource_k8s_wg
	- resource_account
