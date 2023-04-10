package lb

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
)

func existLBID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	lbList := []struct {
		ID int `json:"id"`
	}{}

	lbListAPI := "/restmachine/cloudapi/lb/list"

	lbListRaw, err := c.DecortAPICall(ctx, "POST", lbListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(lbListRaw), &lbList)
	if err != nil {
		return false, err
	}

	haveLB := false
	lbId := d.Get("lb_id").(int)

	for _, lb := range lbList {
		if lb.ID == lbId {
			haveLB = true
			break
		}
	}

	return haveLB, nil
}

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
	extNetID := d.Get("extnet_id").(int)
	for _, extNet := range extNetList {
		if extNet.ID == extNetID {
			haveExtNet = true
			break
		}
	}

	return haveExtNet, nil
}

func existViNSID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	vinsList := []struct {
		ID int `json:"id"`
	}{}

	vinsListAPI := "/restmachine/cloudapi/vins/list"

	vinsListRaw, err := c.DecortAPICall(ctx, "POST", vinsListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(vinsListRaw), &vinsList)
	if err != nil {
		return false, err
	}

	haveVins := false
	vinsID := d.Get("vins_id").(int)
	for _, vins := range vinsList {
		if vins.ID == vinsID {
			haveVins = true
			break
		}
	}

	return haveVins, nil
}
