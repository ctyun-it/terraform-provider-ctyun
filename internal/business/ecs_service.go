package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"strings"
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

func (u EcsService) GetFlavorIDByName(ctx context.Context, flavorName, regionId, azName string) (flavorID string, err error) {
	resp, err := u.meta.Apis.CtEcsApis.EcsFlavorListApi.Do(ctx, u.meta.Credential, &ctecs.EcsFlavorListRequest{
		RegionId:   regionId,
		FlavorName: flavorName,
		AzName:     azName,
	})
	if err != nil {
		return
	}
	if len(resp.FlavorList) == 0 {
		err = fmt.Errorf("云主机规格 %s 不存在", flavorName)
		return
	}
	flavorID = resp.FlavorList[0].FlavorId
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

func (u EcsService) GetEcsAttachedVolume(ctx context.Context, id, regionId string) ([]string, error) {
	instance, err := u.meta.Apis.CtEcsApis.EcsInstanceDetailsApi.Do(ctx, u.meta.Credential, &ctecs.EcsInstanceDetailsRequest{
		RegionId:   regionId,
		InstanceId: id,
	})
	if err != nil {
		// 实例已经被退订的情况
		if err.ErrorCode() == common.EcsInstanceNotFound {
			return nil, fmt.Errorf("云主机 %s 不存在", id)
		}
		return nil, err
	}
	return instance.AttachedVolume, nil
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

// 遍历规格族列表，判定当前规格族的类型：
// 测试结果：
// S7.SMALL.1: S
// S7.MEDIUM.2: S
// C8A.LARGE.4: C
// S8E.LARGE.4: S
// M7.LARGE.8: M
// KS1.MEDIUM.2: KS1
// C7.LARGE.2: C
// C7.LARGE.4: C
// KC1SE.16XLARGE.4: KC1
// HS1.XLARGE.2: HS1
// HS1.XLARGE.4: HS1
// HM1.LARGE.8: HM1
// HC1.LARGE.4: HC1
// HC1.XLARGE.4: HC1
// KC1.LARGE.4: KC1
// KM1.XLARGE.8: KM1
// KC1.LARGE.2: KC1
// KS1.XLARGE.4: KS1
// KS1.MEDIUM.2: KS1
// KS2NE.MEDIUM.2: KS2NE
// KM2NE.XLARGE.8: KM2NE
// KM2X.LARGE.8:
// HS3.LARGE.2:
// HS3X.XLARGE.4:
// HS3.MEDIUM.4:
func (u EcsService) GetInstanceSeries(_ context.Context, hostType string) string {
	for key, _ := range MysqlInstanceSeriesDict {
		if len(hostType) >= len(key) && strings.HasPrefix(hostType, key) {
			return key
		}
	}
	return ""
}
