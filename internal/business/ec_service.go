package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
)

type EcService struct {
	meta *common.CtyunMetadata
}

func NewEcService(meta *common.CtyunMetadata) *EcService {
	return &EcService{meta: meta}
}

func (c EcService) GetCgw(ctx context.Context, ecID, cgwID string) (cgw *ec.EcEcListGatewayReturnObjResultsResponse, err error) {
	listReq := &ec.EcEcListGatewayRequest{
		EcID:  ecID,
		CgwID: &cgwID,
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcListGatewayApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		return nil, fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return nil, common.InvalidReturnObjError
	} else if len(resp.ReturnObj.Results) == 0 {
		return nil, common.ResourceNotExistError
	}
	return resp.ReturnObj.Results[0], nil
}
