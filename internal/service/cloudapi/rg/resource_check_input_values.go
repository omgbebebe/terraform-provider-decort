package rg

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
)

func existAccountID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	accountList := []struct {
		ID int `json:"id"`
	}{}

	accountListAPI := "/restmachine/cloudapi/account/list"

	accountListRaw, err := c.DecortAPICall(ctx, "POST", accountListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(accountListRaw), &accountList)
	if err != nil {
		return false, err
	}

	haveAccount := false

	myAccount := d.Get("account_id").(int)
	for _, account := range accountList {
		if account.ID == myAccount {
			haveAccount = true
			break
		}
	}
	return haveAccount, nil
}
func existGID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}

	locationList := []struct {
		GID int `json:"gid"`
	}{}

	locationsListAPI := "/restmachine/cloudapi/locations/list"

	locationListRaw, err := c.DecortAPICall(ctx, "POST", locationsListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(locationListRaw), &locationList)
	if err != nil {
		return false, err
	}

	haveGID := false

	myGID := d.Get("gid").(int)
	for _, location := range locationList {
		if location.GID == myGID {
			haveGID = true
			break
		}
	}

	return haveGID, nil
}
func existExtNetID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)

	urlValues := &url.Values{}
	urlValues.Add("accountId", strconv.Itoa(d.Get("account_id").(int)))

	listExtNet := []struct {
		ID int `json:"id"`
	}{}

	extNetListAPI := "/restmachine/cloudapi/extnet/list"

	listExtNetRaw, err := c.DecortAPICall(ctx, "POST", extNetListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(listExtNetRaw), &listExtNet)
	if err != nil {
		return false, err
	}

	haveExtNet := false

	myExtNetID := d.Get("ext_net_id").(int)
	for _, extNet := range listExtNet {
		if extNet.ID == myExtNetID {
			haveExtNet = true
			break
		}
	}
	return haveExtNet, nil
}

