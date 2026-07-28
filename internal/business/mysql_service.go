package business

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
)

type MysqlService struct {
	meta *common.CtyunMetadata
}

func NewMysqlService(meta *common.CtyunMetadata) *MysqlService {
	return &MysqlService{meta: meta}
}

func (u MysqlService) GetFlavorByProdIdAndFlavorName(ctx context.Context, prodID string, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, prodId int64, prodSpecName string, version string, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1", // RDS
		ProdCode:     "MYSQL",
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
	pid := MysqlProdIdDict[prodID]
	prodId = pid
	for _, data := range resp.ReturnObj.Data {
		if data.ProdId == pid {
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
	err = fmt.Errorf("invalid flavor_name: %s for prod_id: %s", flavorName, prodID)
	return
}

func (u MysqlService) GetFlavorByInstanceTypeAndFlavorName(ctx context.Context, instanceType int64, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, prodId int64, prodSpecName string, version string, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1", // RDS
		ProdCode:     "MYSQL",
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
	prodId = instanceType
	for _, data := range resp.ReturnObj.Data {
		if data.ProdId == instanceType {
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
	err = fmt.Errorf("invalid flavor_name: %s for prod_id: %d", flavorName, instanceType)
	return
}

func (u MysqlService) GetIDByOrder(ctx context.Context, orderID, projectID string) (id string, err error) {
	params := mysql.TeledbGetIDByOrderRequest{
		OrderID: orderID,
	}
	header := mysql.TeledbGetIDByOrderRequestHeader{}
	if projectID != "" {
		header.ProjectID = projectID
	}
	resp, err := u.meta.Apis.SdkCtMysqlApis.TeledbGetIDByOrderApi.Do(ctx, u.meta.Credential, &params, &header)
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

func (u MysqlService) GetDetailByID(ctx context.Context, instID, regionID string) (instance *mysql.DetailRespReturnObj, err error) {
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: instID,
	}
	detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   instID,
		RegionID: regionID,
	}

	resp, err := u.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, u.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return
	} else if resp.StatusCode != 0 {
		if strings.Contains(resp.Message, "MYSQL_10002") || strings.Contains(resp.Message, "outerProdInstId not exist") {
			err = common.ResourceNotExistError
			return
		}
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = resp.ReturnObj
	return
}

func (u MysqlService) WaitInstanceStatus(ctx context.Context, instID, projectId string, regionID string, runningStatus, orderStatus int32) (err error) {
	retryer, err := NewRetryer(time.Second*20, 180)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			var instance *mysql.DetailRespReturnObj
			instance, err = u.GetDetailByID(ctx, instID, regionID)
			if err != nil {
				return false
			}
			//
			if instance.ProdRunningStatus == runningStatus && instance.ProdOrderStatus == orderStatus && instance.Alive == 0 {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == ReachMaxLoopTime {
		return fmt.Errorf("实例 %s 超过预定时间未达到预期状态", instID)
	}
	return
}

func (u MysqlService) GetVersionByParentProdID(ctx context.Context, regionID string, parentProdID int64, series string) (version string, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "MYSQL",
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
	return "", fmt.Errorf("未查询到 prod_id=%d 的父节点 version", parentProdID)
}

func (u MysqlService) GetReadNodeProdIDByVersion(ctx context.Context, regionID string, version string, series string) (prodID int64, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "1",
		ProdCode:     "MYSQL",
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
		if data.ProdVersion == version && data.ProdSpecName == MysqlProdSpecNameRead {
			return data.ProdId, nil
		}
	}

	return -1, fmt.Errorf("未查询到 version=%s 的只读节点 prod_id", version)
}

// IsProdIDNumeric 判断 prodID 是否为数字字符串，是返回 true，否则返回 false
func (u MysqlService) IsProdIDNumeric(prodID string) (bool, int64) {
	val, err := strconv.ParseInt(prodID, 10, 64)
	if err != nil {
		return false, 0
	}
	return true, val
}
