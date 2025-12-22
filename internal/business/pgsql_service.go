package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
)

type PgsqlService struct {
	meta *common.CtyunMetadata
}

func NewPgsqlService(meta *common.CtyunMetadata) *PgsqlService {
	return &PgsqlService{meta: meta}
}

func (u PgsqlService) GetPgsqlFlavorByProdIdAndFlavorName(ctx context.Context, prodID string, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "POSTGRESQL",
		RegionID:     regionID,
		InstanceType: PgsqlInstanceSeriesDict[series],
	}
	headers := &mysql.TeledbMysqlSpecsRequestHeader{}
	resp, err := u.meta.Apis.SdkCtMysqlApis.TeledbMysqlSpecsApi.Do(ctx, u.meta.Credential, params, headers)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	pid := PgsqlProdIDDict[prodID]
	for _, data := range resp.ReturnObj.Data {
		if data.ProdId == pid {
			for _, spec := range data.InstSpecInfoList {
				if spec.SpecName == flavorName {
					flavor = spec
					return
				}
			}
		}
	}
	err = fmt.Errorf("invalid %s for %s", flavorName, prodID)
	return
}

func (u PgsqlService) GetIDByOrder(ctx context.Context, masterOrderID string, projectID string) (id string, err error) {
	params := pgsql.PgsqlGetIDByOrderRequest{
		OrderID: masterOrderID,
	}
	header := pgsql.PgsqlGetIDByOrderRequestHeader{}
	if projectID != "" {
		header.ProjectID = projectID
	}
	resp, err := u.meta.Apis.SdkCtPgsqlApis.PgsqlGetIDByOrderApi.Do(ctx, u.meta.Credential, &params, &header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if len(resp.ReturnObj.Data) > 0 {
		id = resp.ReturnObj.Data[0]
	}
	return
}

func (u PgsqlService) GetDetailByID(ctx context.Context, id string, projectId string, regionId string) (*pgsql.PgsqlDetailResponseReturnObj, error) {

	// 获取pgsql详情
	detailParams := &pgsql.PgsqlDetailRequest{
		ProdInstId: id,
	}
	detailHeaders := &pgsql.PgsqlDetailRequestHeader{
		RegionID: regionId,
	}
	if projectId != "" {
		detailHeaders.ProjectID = &projectId
	}
	resp, err := u.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, u.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj, nil
}

func (u PgsqlService) AddSecurityGroup(ctx context.Context, instanceID string, projectID string, sgID string) error {
	params := &pgsql.PostgresqlAddSecurityGroupRequest{
		SecurityGroupId: sgID,
		InstanceId:      instanceID,
	}
	header := &pgsql.PostgresqlAddSecurityGroupRequestHeader{
		ProjectId: &projectID,
	}
	resp, err := u.meta.Apis.SdkCtPgsqlApis.PostgresqlAddSecurityGroupApi.Do(ctx, u.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("pgsql实例(id=%s)添加安全组(id=%s)失败，接口返回nil，请联系研发确认问题原因！", instanceID, sgID)
		return err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("pgsql实例(id=%s)添加安全组(id=%s)失败，接口返回错误信息：%s", instanceID, sgID, resp.Message)
		return err
	}
	return nil
}

func (u PgsqlService) RemoveSecurityGroup(ctx context.Context, instanceID string, projectID string, sgID string) error {
	params := &pgsql.PgsqlDeleteSecurityGroupRequest{
		SecurityGroupId: sgID,
		InstanceId:      instanceID,
	}
	header := &pgsql.PgsqlDeleteSecurityGroupRequestHeader{
		ProjectID: &projectID,
	}

	resp, err := u.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteSecurityGroupApi.Do(ctx, u.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("pgsql实例(id=%s)删除安全组(id=%s)失败，接口返回nil，请联系研发确认问题原因！", instanceID, sgID)
		return err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("pgsql实例(id=%s)删除安全组(id=%s)失败，接口返回错误信息：%s", instanceID, sgID, resp.Message)
		return err
	}
	return nil
}
