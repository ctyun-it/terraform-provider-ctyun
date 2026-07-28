package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
)

type PrivateNatService struct {
	meta *common.CtyunMetadata
}

func NewPrivateNatService(meta *common.CtyunMetadata) *PrivateNatService {
	return &PrivateNatService{meta: meta}
}

func (c *PrivateNatService) GetPrivateNatByID(ctx context.Context, id, regionID string) (res *ctnat.CtnatListPrivatenatReturnObjResponse, err error) {
	params := &ctnat.CtnatListPrivatenatRequest{
		RegionID:     regionID,
		NatGatewayID: id,
		PageNo:       1,
		PageSize:     1,
	}
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatListPrivatenatApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil || len(resp.ReturnObj) == 0 {
		err = common.InvalidReturnObjError
		return
	}
	res = resp.ReturnObj[0]
	return
}
