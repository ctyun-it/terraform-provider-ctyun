package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"time"
)

type MysqlService struct {
	meta *common.CtyunMetadata
}

func NewMysqlService(meta *common.CtyunMetadata) *MysqlService {
	return &MysqlService{meta: meta}
}

func (u MysqlService) GetFlavorByProdIdAndFlavorName(ctx context.Context, prodID string, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, err error) {
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

func (u MysqlService) GetDetailByID(ctx context.Context, instID, projectID, regionID string) (instance *mysql.DetailRespReturnObj, err error) {
	detailParams := &mysql.TeledbQueryDetailRequest{
		OuterProdInstId: instID,
	}
	detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
		InstID:   instID,
		RegionID: regionID,
	}
	if projectID != "" {
		detailHeaders.ProjectID = &projectID
	}
	resp, err := u.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, u.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	instance = resp.ReturnObj
	return
}

func (u MysqlService) WaitInstanceStatus(ctx context.Context, instID, projectID, regionID string, runningStatus, orderStatus int32) (err error) {
	retryer, err := NewRetryer(time.Second*20, 180)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			var instance *mysql.DetailRespReturnObj
			instance, err = u.GetDetailByID(ctx, instID, projectID, regionID)
			if err != nil {
				return false
			}
			if instance.ProdRunningStatus == runningStatus && instance.ProdOrderStatus == orderStatus {
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
