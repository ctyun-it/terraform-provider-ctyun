package business

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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

func (u MongodbService) GetMongodbFlavorByProdIdAndFlavorName(ctx context.Context, prodID string, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, specName string, err error) {
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
			specName = data.ProdSpecName
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

func (u MongodbService) GetMongodbFlavorByInstanceTypeAndFlavorName(ctx context.Context, instanceType int64, flavorName, regionID, series string) (flavor mysql.InstSpecInfo, specType string, err error) {
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
	for _, data := range resp.ReturnObj.Data {
		if data.ProdId == instanceType {
			specType = data.ProdSpecName
			for _, spec := range data.InstSpecInfoList {
				if spec.SpecName == flavorName {
					flavor = spec
					return
				}
			}
		}
	}
	err = fmt.Errorf("invalid flavor_name: %s for instance_type: %d", flavorName, instanceType)
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

func (u MongodbService) GetHostIpByInstID(ctx context.Context, instID string, regionID string) (string, error) {
	detail, err := u.GetMongodbDetail(ctx, instID, regionID)
	if err != nil {
		return "", err
	}
	if detail.Host == "" {
		return "", fmt.Errorf("实例 %s 没有唯一host，暂时不支持绑定eip", instID)
	}
	return detail.Host, nil
}

func (u MongodbService) GetMongodbDetail(ctx context.Context, instID string, regionID string) (*mongodb.DetailRespReturnObj, error) {
	detailParams := &mongodb.MongodbQueryDetailRequest{
		ProdInstId: instID,
	}
	detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
		RegionID: regionID,
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

func (u MongodbService) GetShardAndMongosNum(ctx context.Context, regionID string, instID string) (int32, int32, error) {
	detail, err := u.GetMongodbDetail(ctx, instID, regionID)
	if err != nil {
		return 0, 0, err
	}
	shardNum := int32(0)
	mongosNum := int32(0)
	nodeInfoVos := detail.NodeInfoVOS
	for _, node := range nodeInfoVos {
		if node.Role == "Shard" {
			shardNum += 1
		} else if node.Role == "Mongos" {
			mongosNum += 1
		}
	}
	return shardNum, mongosNum, nil
}

// IsProdIDNumeric 判断 prodID 是否为数字字符串，是返回 true 和解析后的值，否则返回 false 和 0
func (u MongodbService) IsProdIDNumeric(prodID string) (bool, int64) {
	val, err := strconv.ParseInt(prodID, 10, 64)
	if err != nil {
		return false, 0
	}
	return true, val
}

