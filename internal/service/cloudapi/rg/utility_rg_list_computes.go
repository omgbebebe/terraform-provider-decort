package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgListComputesCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListComputes, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	listComputes := ListComputes{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	if reason, ok := d.GetOk("reason"); ok {
		urlValues.Add("reason", reason.(string))
	}

	listComputesRaw, err := c.DecortAPICall(ctx, "POST", RgListComputesAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(listComputesRaw), &listComputes)
	if err != nil {
		return nil, err
	}

	return listComputes, nil
}
