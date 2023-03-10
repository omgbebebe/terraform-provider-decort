package kvmvm

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityComputeGetConsoleUrlCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	computeConsoleUrlRaw, err := c.DecortAPICall(ctx, "POST", ComputeGetConsoleUrlAPI, urlValues)
	if err != nil {
		return "", err
	}

	return string(computeConsoleUrlRaw), nil
}
