package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityDataRgUsageCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (*Resource, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := url.Values{}
	usage := Resource{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))

	if reason, ok := d.GetOk("reason"); ok {
		urlValues.Add("reason", reason.(string))
	}

	usageRaw, err := c.DecortAPICall(ctx, "POST", ResgroupUsageAPI, &urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(usageRaw), &usage)
	if err != nil {
		return nil, err
	}

	return &usage, nil
}
