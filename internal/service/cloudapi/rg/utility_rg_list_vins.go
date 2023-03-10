package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgListVinsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListVINS, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	listVins := ListVINS{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))

	if val, ok := d.GetOk("reason"); ok {
		urlValues.Add("reason", val.(string))
	}

	listVinsRaw, err := c.DecortAPICall(ctx, "POST", ResgroupListVinsAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(listVinsRaw), &listVins)
	if err != nil {
		return nil, err
	}

	return listVins, nil
}
