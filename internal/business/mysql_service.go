package business

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
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
