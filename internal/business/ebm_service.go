package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
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
	} else if utils.SecString(resp.ErrorCode) == common.OpenapiEbmNotFound {
		err = common.ResourceNotExistError
		return
	} else if resp.StatusCode != common.NormalStatusCode {
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
	return strings.ToLower(utils.SecString(instance.EbmState)), err
}

// GetEbmImageID 查询物理机镜像id
func (c EbmService) GetEbmImageID(ctx context.Context, id, regionID, azName string) (status string, err error) {
	instance, err := c.GetEbmInfo(ctx, id, regionID, azName)
	if err != nil {
		return
	}
	return utils.SecString(instance.ImageID), err
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

func (c EbmService) GetDeviceTypeStock(ctx context.Context, regionID, azName string) (stocks map[string]int32, err error) {
	stocks = make(map[string]int32)
	params := ctebm.EbmDeviceStockListRequest{
		RegionID: regionID,
		AzName:   azName,
	}
	resp, err := c.meta.Apis.CtEbmApis.EbmDeviceStockListApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj.Results) == 0 {
		return
	}
	for _, d := range resp.ReturnObj.Results[0].Stocks {
		stocks[utils.SecString(d.DeviceType)] = d.Available
	}
	return
}
