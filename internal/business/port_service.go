package business

import (
	"context"
	"errors"
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
	_, err = v.GetPortDetail(ctx, portID, regionID)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			return false, nil
		}
		return
	}
	return true, nil
}

func (v PortService) GetPortDetail(ctx context.Context, portID, regionID string) (re *ctvpc.CtvpcShowPortReturnObjResponse, err error) {
	params := &ctvpc.CtvpcShowPortRequest{
		RegionID:           regionID,
		NetworkInterfaceID: portID,
	}
	resp, err := v.meta.Apis.SdkCtVpcApis.CtvpcShowPortApi.Do(ctx, v.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		if strings.Contains(*resp.Message, "is not exists") {
			err = common.ResourceNotExistError
		} else {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return resp.ReturnObj, nil
}
