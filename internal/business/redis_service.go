package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
)

type RedisService struct {
	meta *common.CtyunMetadata
}

func NewRedisService(meta *common.CtyunMetadata) *RedisService {
	return &RedisService{meta: meta}
}

func (c RedisService) GetRedisByID(ctx context.Context, id, regionID string) (*dcs2.Dcs2DescribeInstancesOverviewReturnObjResponse, error) {
	params := &dcs2.Dcs2DescribeInstancesOverviewRequest{
		RegionId:   regionID,
		ProdInstId: id,
	}
	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeInstancesOverviewApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp.Error == common.OpenapiDcs2NotFound || resp.Error == common.OpenapiDcs2Invalid {
		return nil, common.ResourceNotExistError
	} else if resp.StatusCode != common.NormalStatusCode {
		return nil, fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
	} else if resp.ReturnObj == nil {
		return nil, common.InvalidReturnObjError
	}
	return resp.ReturnObj, nil
}
