package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctvpc2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	"github.com/google/uuid"
)

type VpcService struct {
	meta *common.CtyunMetadata
}

func NewVpcService(meta *common.CtyunMetadata) *VpcService {
	return &VpcService{meta: meta}
}

func (v VpcService) MustExist(ctx context.Context, vpcId, regionId string) error {
	_, err := v.meta.Apis.CtVpcApis.VpcQueryApi.Do(ctx, v.meta.Credential, &ctvpc.VpcQueryRequest{
		RegionId:    regionId,
		ClientToken: uuid.NewString(),
		VpcId:       vpcId,
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiVpcNotFound {
			return fmt.Errorf("%s 在 %s 不存在", vpcId, regionId)
		}
		return err
	}
	return nil
}

func (v VpcService) GetVpcDetail(ctx context.Context, vpcId, regionId string) (*ctvpc.VpcQueryResponse, error) {
	resp, err := v.meta.Apis.CtVpcApis.VpcQueryApi.Do(ctx, v.meta.Credential, &ctvpc.VpcQueryRequest{
		RegionId:    regionId,
		ClientToken: uuid.NewString(),
		VpcId:       vpcId,
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiVpcNotFound {
			return nil, fmt.Errorf("vpc %s 不存在", vpcId)
		}
		return nil, err
	}
	return resp, nil
}

func (v VpcService) GetVpcID(ctx context.Context, vpcName, regionId, projectId string) (string, error) {
	params := &ctvpc2.CtvpcNewVpcListRequest{
		RegionID: regionId,
		VpcName:  &vpcName,
		PageNo:   1,
		PageSize: 1000,
	}
	if projectId != "" {
		params.ProjectID = &projectId
	}
	// 调用API
	resp, err := v.meta.Apis.SdkCtVpcApis.CtvpcNewVpcListApi.Do(ctx, v.meta.SdkCredential, params)
	if err != nil {
		return "", err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return "", err
	} else if resp.ReturnObj == nil || resp.ReturnObj.Vpcs == nil || len(resp.ReturnObj.Vpcs) == 0 || len(resp.ReturnObj.Vpcs) > 1 {
		err = common.InvalidReturnObjError
		return "", err
	}
	return *resp.ReturnObj.Vpcs[0].VpcID, nil
}

func (v VpcService) GetVpcSubnet(ctx context.Context, vpcId, regionId string) (map[string]ctvpc.SubnetListSubnetsResponse, error) {
	params := &ctvpc.VpcQueryRequest{
		RegionId:    regionId,
		ClientToken: uuid.NewString(),
		VpcId:       vpcId,
	}
	resp, err := v.meta.Apis.CtVpcApis.VpcQueryApi.Do(ctx, v.meta.Credential, params)
	if err != nil {
		if err.ErrorCode() == common.OpenapiVpcNotFound {
			return nil, fmt.Errorf("vpc %s 不存在", vpcId)
		}
		return nil, err
	}
	r, err := v.meta.Apis.CtVpcApis.SubnetListApi.Do(ctx, v.meta.Credential, &ctvpc.SubnetListRequest{
		RegionId:   regionId,
		VpcId:      vpcId,
		SubnetIds:  resp.SubnetIds,
		PageNumber: 1,
		PageSize:   100,
	})
	if err != nil {
		return nil, err
	}
	result := make(map[string]ctvpc.SubnetListSubnetsResponse)
	for _, subnet := range r.Subnets {
		result[subnet.SubnetId] = subnet
	}
	return result, nil
}
