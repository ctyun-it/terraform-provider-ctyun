package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
)

type IPv6BandwidthService struct {
	meta *common.CtyunMetadata
}

func NewIPv6BandwidthService(meta *common.CtyunMetadata) *IPv6BandwidthService {
	return &IPv6BandwidthService{meta: meta}
}

func (u IPv6BandwidthService) MustExist(ctx context.Context, id, regionId string) error {
	resp, err := u.meta.Apis.SdkCtVpcApis.CtvpcShowIPv6BandwidthApi.Do(ctx, u.meta.SdkCredential, &ctvpc.CtvpcShowIPv6BandwidthRequest{
		RegionID:    regionId,
		BandwidthID: id,
	})

	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	} else if utils.SecString(resp.ErrorCode) == common.OpenapiIPv6BandwidthNotFound {
		err = common.ResourceNotExistError
		return err
	}

	return nil
}
