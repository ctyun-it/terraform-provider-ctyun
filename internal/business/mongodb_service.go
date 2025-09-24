package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
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
