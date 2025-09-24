package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctvpc2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	"github.com/google/uuid"
)

type EipService struct {
	meta *common.CtyunMetadata
}

func NewEipService(meta *common.CtyunMetadata) *EipService {
	return &EipService{meta: meta}
}

func (u EipService) MustExist(ctx context.Context, id, regionId string) error {
	_, err := u.meta.Apis.CtVpcApis.EipShowApi.Do(ctx, u.meta.Credential, &ctvpc.EipShowRequest{
		RegionId: regionId,
		EipId:    id,
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiEipNotFound {
			return fmt.Errorf("eip %s 不存在", id)
		}
		return err
	}
	return nil
}

func (u EipService) GetEipByAddress(ctx context.Context, address, regionId string) (*ctvpc2.CtvpcNewEipListReturnObjEipsResponse, error) {
	resp, err := u.meta.Apis.SdkCtVpcApis.CtvpcNewEipListApi.Do(ctx, u.meta.SdkCredential, &ctvpc2.CtvpcNewEipListRequest{
		RegionID:    regionId,
		Ip:          &address,
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	} else if resp.StatusCode == common.ErrorStatusCode {
		return nil, fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		return nil, common.InvalidReturnObjError
	} else if len(resp.ReturnObj.Eips) == 0 {
		return nil, fmt.Errorf("弹性IP地址 %s 不存在", address)
	}
	ip := resp.ReturnObj.Eips[0]
	return ip, nil
}

func (u EipService) GetEipAddressByEipID(ctx context.Context, eipID, regionId string) (*ctvpc2.CtvpcNewEipListReturnObjEipsResponse, error) {
	resp, err := u.meta.Apis.SdkCtVpcApis.CtvpcNewEipListApi.Do(ctx, u.meta.SdkCredential, &ctvpc2.CtvpcNewEipListRequest{
		RegionID:    regionId,
		Ids:         []*string{&eipID},
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	} else if resp.StatusCode == common.ErrorStatusCode {
		return nil, fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil {
		return nil, common.InvalidReturnObjError
	} else if len(resp.ReturnObj.Eips) == 0 {
		return nil, fmt.Errorf("eipID %s 不存在", eipID)
	}
	ip := resp.ReturnObj.Eips[0]
	return ip, nil
}
