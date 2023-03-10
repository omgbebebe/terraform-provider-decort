package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgAffinityGroupComputesCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListAffinityGroupCompute, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	listGroupComputes := ListAffinityGroupCompute{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	urlValues.Add("affinityGroup", d.Get("affinity_group").(string))

	listGroupComputesRaw, err := c.DecortAPICall(ctx, "POST", RgAffinityGroupComputesAPI, urlValues)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(listGroupComputesRaw), &listGroupComputes)
	if err != nil {
		return nil, err
	}

	return listGroupComputes, nil
}
