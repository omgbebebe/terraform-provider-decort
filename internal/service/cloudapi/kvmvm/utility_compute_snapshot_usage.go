package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
)

func utilityComputeSnapshotUasgeCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListUsageSnapshots, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	UsageSnapshotList := &ListUsageSnapshots{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	if label, ok := d.GetOk("label"); ok {
		urlValues.Add("label", label.(string))
	}
	computeSnapshotUsage, err := c.DecortAPICall(ctx, "POST", ComputeSnapshotUsageAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(computeSnapshotUsage), &UsageSnapshotList)
	if err != nil {
		return nil, err
	}
	return *UsageSnapshotList, err

}
