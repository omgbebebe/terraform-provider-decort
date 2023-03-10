package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgListDeletedCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListResourceGroups, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	rgList := ListResourceGroups{}

	if size, ok := d.GetOk("size"); ok {
		urlValues.Add("size", strconv.Itoa(size.(int)))
	}
	if page, ok := d.GetOk("page"); ok {
		urlValues.Add("page", strconv.Itoa(page.(int)))
	}

	rgListRaw, err := c.DecortAPICall(ctx, "POST", ResgroupListDeletedAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(rgListRaw), &rgList)
	if err != nil {
		return nil, err
	}

	return rgList, nil
}
