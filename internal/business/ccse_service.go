package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
)

type CcseService struct {
	meta *common.CtyunMetadata
}

func NewCcseService(meta *common.CtyunMetadata) *CcseService {
	return &CcseService{meta: meta}
}

func (c CcseService) GetCcseInfo(ctx context.Context, id, regionID string) (instance *ccse.CcseGetClusterReturnObjResponse, err error) {
	params := &ccse.CcseGetClusterRequest{
		RegionId:  regionID,
		ClusterId: id,
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		if resp.Error == common.OpenapiCCSENotExist {
			err = common.ResourceNotExistError
		} else {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = resp.ReturnObj
	return
}
