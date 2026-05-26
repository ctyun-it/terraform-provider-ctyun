package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
)

type IPv6Service struct {
	meta *common.CtyunMetadata
}

func NewIPv6Service(meta *common.CtyunMetadata) *IPv6Service {
	return &IPv6Service{meta: meta}
}

func (u IPv6Service) MustExist(ctx context.Context, ip, regionId string) error {
	resp, err := u.meta.Apis.SdkCtVpcApis.CtvpcNewIPv6ListApi.Do(ctx, u.meta.SdkCredential, &ctvpc.CtvpcNewIPv6ListRequest{
		RegionID:  regionId,
		IpAddress: &ip,
	})

	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	} else if len(resp.ReturnObj.IPv6s) == 0 {
		err = common.ResourceNotExistError
		return err
	}

	return nil
}
