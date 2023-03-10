package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgListPfwCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListPFW, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	listPfw := ListPFW{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))

	listPfwRaw, err := c.DecortAPICall(ctx, "POST", ResgroupListPfwAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(listPfwRaw), &listPfw)
	if err != nil {
		return nil, err
	}

	return listPfw, nil
}
