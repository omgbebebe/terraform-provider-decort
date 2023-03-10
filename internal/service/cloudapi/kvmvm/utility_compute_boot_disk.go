package kvmvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func utilityComputeBootDiskCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*ItemComputeDisk, error) {
	compute, err := utilityComputeCheckPresence(ctx, d, m)
	if err != nil {
		return nil, err
	}

	bootDisk := &ItemComputeDisk{}
	for _, disk := range compute.Disks {
		if disk.Name == "bootdisk" {
			*bootDisk = disk
			break
		}
	}
	return bootDisk, nil
}
