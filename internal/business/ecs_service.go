package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"time"
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

func (u EcsService) GetFlavorByName(ctx context.Context, flavorName, regionId string) (flavor ctecs.EcsFlavorListFlavorListResponse, err error) {
	resp, err := u.meta.Apis.CtEcsApis.EcsFlavorListApi.Do(ctx, u.meta.Credential, &ctecs.EcsFlavorListRequest{
		RegionId:   regionId,
		FlavorName: flavorName,
	})
	if err != nil {
		return
	}
	if len(resp.FlavorList) == 0 {
		err = fmt.Errorf("云主机规格 %s 不存在", flavorName)
		return
	}
	flavor = resp.FlavorList[0]
	// 因为没传azName，所以flavor_id是不准确的
	flavor.FlavorId = "invalid"
	return
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

func (u EcsService) GetEcsStatus(ctx context.Context, id, regionId string) (string, error) {
	instance, err := u.meta.Apis.CtEcsApis.EcsInstanceDetailsApi.Do(ctx, u.meta.Credential, &ctecs.EcsInstanceDetailsRequest{
		RegionId:   regionId,
		InstanceId: id,
	})
	if err != nil {
		// 实例已经被退订的情况
		if err.ErrorCode() == common.EcsInstanceNotFound {
			return "", fmt.Errorf("云主机 %s 不存在", id)
		}
		return "", err
	}
	return instance.InstanceStatus, nil
}

func (u EcsService) CheckEcsStatus(ctx context.Context, id, regionId string) error {
	var executeSuccessFlag bool
	var status string
	var err error
	retryer, _ := NewRetryer(time.Second*10, 10)
	retryer.Start(
		func(currentTime int) bool {
			status, err = u.GetEcsStatus(ctx, id, regionId)
			if err != nil {
				return false
			}
			switch status {
			case EcsStatusRunning, EcsStatusStopped, EcsStatusShelve:
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return err
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云主机当前状态异常：%s", status)
	}
	return nil
}
