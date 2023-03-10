package kvmvm

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	log "github.com/sirupsen/logrus"
)

func existRgID(ctx context.Context, d *schema.ResourceData, m interface{}) bool {
	log.Debugf("resourceComputeCreate: check access for RG ID: %v", d.Get("rg_id").(int))
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	rgList := []struct {
		ID int `json:"id"`
	}{}

	rgListAPI := "/restmachine/cloudapi/rg/list"
	urlValues.Add("includedeleted", "false")
	rgListRaw, err := c.DecortAPICall(ctx, "POST", rgListAPI, urlValues)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(rgListRaw), &rgList)
	if err != nil {
		return false
	}
	rgId := d.Get("rg_id").(int)
	for _, rg := range rgList {
		if rg.ID == rgId {
			return true
		}
	}
	return false
}

func existImageId(ctx context.Context, d *schema.ResourceData, m interface{}) bool {
	log.Debugf("resourceComputeCreate: check access for image ID: %v", d.Get("image_id").(int))
	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	imageList := []struct {
		ID int `json:"id"`
	}{}
	imageListAPI := "/restmachine/cloudapi/image/list"
	imageListRaw, err := c.DecortAPICall(ctx, "POST", imageListAPI, urlValues)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(imageListRaw), &imageList)
	if err != nil {
		return false
	}
	imageId := d.Get("image_id").(int)
	for _, image := range imageList {
		if image.ID == imageId {
			return true
		}
	}
	return false
}

func existVinsIdInList(vinsId int, vinsList []struct {
	ID int `json:"id"`
}) bool {
	for _, vins := range vinsList {
		if vinsId == vins.ID {
			return true
		}
	}
	return false
}

func existVinsId(ctx context.Context, d *schema.ResourceData, m interface{}) (int, bool) {
	log.Debugf("resourceComputeCreate: check access for vinses IDs")

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}

	vinsListAPI := "/restmachine/cloudapi/vins/list"
	urlValues.Add("includeDeleted", "false")
	vinsList := []struct {
		ID int `json:"id"`
	}{}
	vinsListRaw, err := c.DecortAPICall(ctx, "POST", vinsListAPI, urlValues)
	if err != nil {
		return 0, false
	}
	err = json.Unmarshal([]byte(vinsListRaw), &vinsList)
	if err != nil {
		return 0, false
	}

	networks := d.Get("network").(*schema.Set).List()

	for _, networkInterface := range networks {

		networkItem := networkInterface.(map[string]interface{})
		if !existVinsIdInList(networkItem["net_id"].(int), vinsList) {
			return networkItem["net_id"].(int), false
		}
	}
	return 0, true
}
