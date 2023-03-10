package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityRgAuditsCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (ListAudits, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	rgAudits := ListAudits{}

	urlValues.Add("rgId", strconv.Itoa(d.Get("rg_id").(int)))
	rgAuditsRow, err := c.DecortAPICall(ctx, "POST", RgAuditsAPI, urlValues)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(rgAuditsRow), &rgAudits)
	if err != nil {
		return nil, err
	}

	return rgAudits, nil
}
