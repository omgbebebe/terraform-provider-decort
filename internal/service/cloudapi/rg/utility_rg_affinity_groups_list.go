package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgAffinityGroupsListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (map[string][]uint64, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	groups := make(map[string][]uint64, 0)

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))

	groupsRaw, err := c.DecortAPICall(ctx, "POST", RgAffinityGroupsListAPI, urlValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(groupsRaw), &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
