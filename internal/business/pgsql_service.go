package business

import (
	"context"
	"fmt"
	"strings"

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

func (u PgsqlService) GetPgsqlFlavorByProdIdAndFlavorName(ctx context.Context, prodID int64, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, prodSpecName string, version string, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "POSTGRESQL",
		RegionID:     regionID,
		InstanceType: MysqlInstanceSeriesDict[series],
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
	for _, data := range resp.ReturnObj.Data {
		if data.ProdId == prodID {
			prodSpecName = data.ProdSpecName
			version = data.ProdVersion
			for _, spec := range data.InstSpecInfoList {
				if spec.SpecName == flavorName {
					flavor = spec
					return
				}
			}
		}
	}
	err = fmt.Errorf("invalid %s for %d", flavorName, prodID)
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
		if strings.Contains(resp.Error, "PG_2001") || strings.Contains(resp.Message, "未找到实例") {
			return nil, common.ResourceNotExistError
		}
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

func (u PgsqlService) GetParentByID(ctx context.Context, regionID string, readOnlyInstID string) (*pgsql.PgsqlListResponsePageInfo, error) {
	// 查询该账号下所有的pgsql实例
	instList, err := u.GetList(ctx, regionID, nil, nil)
	if err != nil {
		return nil, err
	}
	// 遍历实例，如果实例有只读节点，遍历并与只读的inst id对比
	for _, inst := range instList {
		if inst.ReadonlyInstnaceIds != nil && len(inst.ReadonlyInstnaceIds) > 0 {
			for _, item := range inst.ReadonlyInstnaceIds {
				if item == readOnlyInstID {
					return &inst, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("未查询到id=%s的只读实例！", readOnlyInstID)
}

func (u PgsqlService) GetList(ctx context.Context, regionID string, instID *string, name *string) ([]pgsql.PgsqlListResponsePageInfo, error) {
	params := &pgsql.PgsqlListRequest{
		PageNum:  1,
		PageSize: 100,
	}
	if name != nil {
		params.ProdInstName = name
	}
	if instID != nil {
		params.ProdInstId = instID
	}
	header := &pgsql.PgsqlListRequestHeader{
		RegionID: regionID,
	}
	resp, err := u.meta.Apis.SdkCtPgsqlApis.PgsqlListApi.Do(ctx, u.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询pgsql的实例列表信息失败，接口返回nil，请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil || resp.ReturnObj.List == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj.List, nil
}

func (u PgsqlService) GetReadNodeProdIDByVersion(ctx context.Context, regionID string, version string, series string) (prodID int64, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "POSTGRESQL",
		RegionID:     regionID,
		InstanceType: MysqlInstanceSeriesDict[series],
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

	for _, data := range resp.ReturnObj.Data {
		if data.ProdVersion == version && data.ProdSpecName == PgsqlProdSpecNameRead {
			return data.ProdId, nil
		}
	}

	return -1, fmt.Errorf("未查询到version=%s的只读节点prod_id", version)

}

func (u PgsqlService) GetVersionByParentProdID(ctx context.Context, regionID string, parentProdID int64, series string) (version string, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "POSTGRESQL",
		RegionID:     regionID,
		InstanceType: MysqlInstanceSeriesDict[series],
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
	for _, data := range resp.ReturnObj.Data {
		if data.ProdId == parentProdID {
			return data.ProdVersion, nil
		}
	}
	return "", fmt.Errorf("未查询到prod_id=%d的父节点version", parentProdID)

}
