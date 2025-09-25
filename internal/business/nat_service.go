package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
)

type NatService struct {
	meta *common.CtyunMetadata
}

func NewNatService(meta *common.CtyunMetadata) *NatService {
	return &NatService{meta: meta}
}

func (c NatService) GetNatByID(ctx context.Context, id, regionID string) (res *ctvpc.CtvpcShowNatGatewayReturnObjResponse, err error) {
	params := &ctvpc.CtvpcShowNatGatewayRequest{
		RegionID:     regionID,
		NatGatewayID: id,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowNatGatewayApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	res = resp.ReturnObj
	return
}
