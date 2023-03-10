package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityComputePfwListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListPFWs, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	listPFWs := &ListPFWs{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	computePfwListRaw, err := c.DecortAPICall(ctx, "POST", ComputePfwListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(computePfwListRaw), &listPFWs)
	if err != nil {
		return nil, err
	}
	return *listPFWs, err

}
