package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctecs"
	"terraform-provider-ctyun/internal/common"
)

type EcsService struct {
	meta *common.CtyunMetadata
}

func NewEcsService(meta *common.CtyunMetadata) *EcsService {
	return &EcsService{meta: meta}
}

func (u EcsService) FlavorMustExist(ctx context.Context, flavorId, regionId, azName string) error {
	resp, err := u.meta.Apis.CtEcsApis.EcsFlavorListApi.Do(ctx, u.meta.Credential, &ctecs.EcsFlavorListRequest{
		RegionId: regionId,
		AzName:   azName,
		FlavorId: flavorId,
	})
	if err != nil {
		return err
	}
	if len(resp.FlavorList) == 0 {
		return fmt.Errorf("云主机规格 %s 不存在", flavorId)
	}
	return nil
}

func (u EcsService) MustExist(ctx context.Context, id, regionId string) error {
	_, err := u.meta.Apis.CtEcsApis.EcsInstanceDetailsApi.Do(ctx, u.meta.Credential, &ctecs.EcsInstanceDetailsRequest{
		RegionId:   regionId,
		InstanceId: id,
	})
	if err != nil {
		// 实例已经被退订的情况
		if err.ErrorCode() == common.EcsInstanceNotFound {
			return fmt.Errorf("云主机 %s 不存在", id)
		}
		return err
	}
	return nil
}
