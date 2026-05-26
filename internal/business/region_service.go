package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
)

type RegionService struct {
	meta *common.CtyunMetadata
}

func NewRegionService(meta *common.CtyunMetadata) *RegionService {
	return &RegionService{meta: meta}
}

func (c RegionService) GetZonesByRegionID(ctx context.Context, regionID string) (zones []string, err error) {
	params := &ctecs.CtecsQueryZonesInRegionV41Request{
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQueryZonesInRegionV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	for _, r := range resp.ReturnObj.ZoneList {
		zones = append(zones, r.Name)
	}
	return
}

func (c RegionService) GetZonesCnByRegionID(ctx context.Context, regionID string) (zones map[string]string, err error) {
	params := &ctecs.CtecsQueryZonesInRegionV41Request{
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsQueryZonesInRegionV41Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	zones = map[string]string{}
	for _, r := range resp.ReturnObj.ZoneList {
		zones[r.AzDisplayName] = r.Name
	}
	return
}
