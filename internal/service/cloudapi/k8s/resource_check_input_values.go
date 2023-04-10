package k8s

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"repository.basistech.ru/BASIS/terraform-provider-decort/internal/controller"
)

func existK8sID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	k8sList := []struct {
		ID int `json:"id"`
	}{}

	k8sListAPI := "/restmachine/cloudapi/k8s/list"

	k8sListRaw, err := c.DecortAPICall(ctx, "POST", k8sListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(k8sListRaw), &k8sList)
	if err != nil {
		return false, err
	}

	haveK8s := false
	k8sID := d.Get("k8s_id").(int)
	for _, k8s := range k8sList {
		if k8s.ID == k8sID {
			haveK8s = true
			break
		}
	}

	return haveK8s, nil
}

func existK8sCIID(ctx context.Context, d *schema.ResourceData, m interface{}) (bool, error) {
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	k8sciList := []struct {
		ID int `json:"id"`
	}{}

	k8sciListAPI := "/restmachine/cloudapi/k8ci/list"

	k8sciListRaw, err := c.DecortAPICall(ctx, "POST", k8sciListAPI, urlValues)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(k8sciListRaw), &k8sciList)
	if err != nil {
		return false, err
	}

	haveK8sCI := false
	k8sciID := d.Get("k8sci_id").(int)
	for _, k8ci := range k8sciList {
		if k8ci.ID == k8sciID {
			haveK8sCI = true
			break
		}
	}

	return haveK8sCI, nil
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
	extNetID := d.Get("extnet_id").(int)

	if extNetID == 0 {
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
