package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityDataComputeListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListComputes, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	listComputes := &ListComputes{}

	if includeDeleted, ok := d.GetOk("includedeleted"); ok {
		urlValues.Add("includeDeleted", strconv.FormatBool(includeDeleted.(bool)))
	}
	if page, ok := d.GetOk("page"); ok {
		urlValues.Add("page", strconv.Itoa(page.(int)))
	}
	if size, ok := d.GetOk("size"); ok {
		urlValues.Add("size", strconv.Itoa(size.(int)))
	}

	listComputesRaw, err := c.DecortAPICall(ctx, "POST", ComputeListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(listComputesRaw), &listComputes)
	if err != nil {
		return nil, err
	}
	return *listComputes, nil

}
