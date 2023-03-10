package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgListLbCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListLB, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	listLb := ListLB{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))

	listLbRaw, err := c.DecortAPICall(ctx, "POST", ResgroupListLbAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(listLbRaw), &listLb)
	if err != nil {
		return nil, err
	}

	return listLb, nil
}
