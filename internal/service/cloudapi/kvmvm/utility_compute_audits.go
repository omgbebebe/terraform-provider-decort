package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityComputeAuditsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListAudits, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	computeAudits := &ListAudits{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	computeAuditsRaw, err := c.DecortAPICall(ctx, "POST", ComputeAuditsAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(computeAuditsRaw), &computeAudits)
	if err != nil {
		return nil, err
	}
	return *computeAudits, nil
}
