package business

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
)

type MongodbService struct {
	meta *common.CtyunMetadata
}

func NewMongodbService(meta *common.CtyunMetadata) *MongodbService {
	return &MongodbService{meta: meta}
}

func (u MongodbService) GetMongodbFlavorByProdIdAndFlavorName(ctx context.Context, prodID string, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, err error) {
	params := &mysql.TeledbMysqlSpecsRequest{
		ProdType:     "2", // RDS
		ProdCode:     "DDS",
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
	pid := MongodbProdIDDict[prodID]
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

func (u MongodbService) GetIDByOrder(ctx context.Context, masterOrderID string, projectID string) (id string, err error) {
	params := mongodb.MongodbGetIDByOrderRequest{
		OrderID: masterOrderID,
	}
	header := mongodb.MongodbGetIDByOrderRequestHeader{}
	if projectID != "" {
		header.ProjectID = projectID
	}
	resp, err := u.meta.Apis.SdkMongodbApis.MongodbGetIDByOrderApi.Do(ctx, u.meta.Credential, &params, &header)
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

func (u MongodbService) GetHostIpByInstID(ctx context.Context, instID string, regionID string, projectID string) (string, error) {
	detail, err := u.GetMongodbDetail(ctx, instID, regionID, projectID)
	if err != nil {
		return "", err
	}
	return detail.Host, nil
}

func (u MongodbService) GetMongodbDetail(ctx context.Context, instID string, regionID string, projectID string) (*mongodb.DetailRespReturnObj, error) {
	detailParams := &mongodb.MongodbQueryDetailRequest{
		ProdInstId: instID,
	}
	detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
		RegionID: regionID,
	}
	if projectID != "" {
		detailHeader.ProjectID = &projectID
	}
	resp, err := u.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, u.meta.Credential, detailParams, detailHeader)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("获取mongodb实例为nil，请稍后再试！")
		return nil, err
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	detail := resp.ReturnObj
	return detail, nil
}
