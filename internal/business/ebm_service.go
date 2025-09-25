package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	"strings"
)

type EbmService struct {
	meta *common.CtyunMetadata
}

func NewEbmService(meta *common.CtyunMetadata) *EbmService {
	return &EbmService{meta: meta}
}

func (c EbmService) GetEbmInfo(ctx context.Context, id, regionID, azName string) (instance ctebm.EbmDescribeInstanceV4plusReturnObjResponse, err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDescribeInstanceV4plusApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDescribeInstanceV4plusRequest{
		RegionID:     regionID,
		InstanceUUID: id,
		AzName:       azName,
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = *resp.ReturnObj
	return
}

// GetEbmStatus 查询物理机状态
func (c EbmService) GetEbmStatus(ctx context.Context, id, regionID, azName string) (status string, err error) {
	instance, err := c.GetEbmInfo(ctx, id, regionID, azName)
	if err != nil {
		return
	}
	return strings.ToLower(*instance.EbmState), err
}

func (c EbmService) GetDeviceType(ctx context.Context, deviceType, regionID, azName string) (spec ctebm.EbmDeviceTypeListReturnObjResultsResponse, err error) {
	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceTypeListApi.Do(ctx, c.meta.SdkCredential, &ctebm.EbmDeviceTypeListRequest{
		RegionID:   regionID,
		DeviceType: &deviceType,
		AzName:     azName,
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil || len(resp.ReturnObj.Results) == 0 {
		err = common.InvalidReturnObjError
		return
	}
	spec = *resp.ReturnObj.Results[0]
	return
}
