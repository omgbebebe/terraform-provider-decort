package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityDataComputeCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (RecordCompute, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	compute := &RecordCompute{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	computeRaw, err := c.DecortAPICall(ctx, "POST", ComputeGetAPI, urlValues)
	if err != nil {
		return *compute, err
	}

	err = json.Unmarshal([]byte(computeRaw), &compute)
	if err != nil {
		return *compute, err
	}
	return *compute, nil
}
