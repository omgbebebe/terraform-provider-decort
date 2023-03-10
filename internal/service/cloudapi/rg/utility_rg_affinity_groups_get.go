package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgAffinityGroupsGetCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) ([]uint64, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	computes := make([]uint64, 0)

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("affinityGroup", d.Get("affinity_group").(string))

    computesRaw, err := c.DecortAPICall(ctx, "POST", RgAffinityGroupsGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(computesRaw), &computes)
	if err != nil {
		return nil, err
	}

	return computes, nil
}
