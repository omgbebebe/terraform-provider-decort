package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityComputeGetAuditsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListShortAudits, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	computeAudits := &ListShortAudits{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	computeAuditsRaw, err := c.DecortAPICall(ctx, "POST", ComputeGetAuditsAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(computeAuditsRaw), &computeAudits)
	if err != nil {
		return nil, err
	}
	return *computeAudits, nil
}
