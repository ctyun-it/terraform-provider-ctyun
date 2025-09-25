package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"strings"
)

type PortService struct {
	meta *common.CtyunMetadata
}

func NewPortService(meta *common.CtyunMetadata) *PortService {
	return &PortService{meta: meta}
}

func (v PortService) Exist(ctx context.Context, portID, regionID string) (exist bool, err error) {
	params := &ctvpc.CtvpcShowPortRequest{
		RegionID:           regionID,
		NetworkInterfaceID: portID,
	}
	resp, err := v.meta.Apis.SdkCtVpcApis.CtvpcShowPortApi.Do(ctx, v.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		if strings.Contains(*resp.Message, "is not exists") {
			return false, nil
		}
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return true, nil
}
