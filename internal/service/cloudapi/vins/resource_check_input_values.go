package vins

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
)

func existRGID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	rgList := []struct {
		ID int `json:"id"`
	}{}

	rgListAPI := "/restmachine/cloudapi/rg/list"

	rgListRaw, err := c.DecortAPICall(ctx, "POST", rgListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(rgListRaw), &rgList)
	if err != nil {
		return false, err
	}

	haveRG := false
	rgId := d.Get("rg_id").(int)
	for _, rg := range rgList {
		if rg.ID == rgId {
			haveRG = true
			break
		}
	}

	return haveRG, nil
}

func existExtNetID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	extNetID := d.Get("ext_net_id").(int)

	if extNetID == 0 || extNetID == -1 {
		return true, nil
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	extNetList := []struct {
		ID int `json:"id"`
	}{}

	extNetListAPI := "/restmachine/cloudapi/extnet/list"

	extNetListRaw, err := c.DecortAPICall(ctx, "POST", extNetListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(extNetListRaw), &extNetList)
	if err != nil {
		return false, err
	}

	haveExtNet := false
	for _, extNet := range extNetList {
		if extNet.ID == extNetID {
			haveExtNet = true
			break
		}
	}

	return haveExtNet, nil
}

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

	gid := d.Get("gid").(int)
	for _, location := range locationList {
		if location.GID == gid {
			haveGID = true
			break
		}
	}

	return haveGID, nil
}
