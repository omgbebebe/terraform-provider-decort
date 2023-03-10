package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func utilityComputeUserListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (RecordACL, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	userList := &RecordACL{}

	urlValues.Add("computeId", strconv.Itoa(d.Get("compute_id").(int)))
	computeUserListRaw, err := c.DecortAPICall(ctx, "POST", ComputeUserListAPI, urlValues)
	if err != nil {
		return *userList, err
	}
	err = json.Unmarshal([]byte(computeUserListRaw), &userList)
	if err != nil {
		return *userList, err
	}
	return *userList, err
}
