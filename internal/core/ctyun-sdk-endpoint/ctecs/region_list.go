package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// RegionListApi 资源池列表查询
// https://www.ctyun.cn/document/10026730/10040588
type RegionListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewRegionListApi(client *ctyunsdk.CtyunClient) *RegionListApi {
	return &RegionListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/region/list-regions",
		},
	}
}

func (this *RegionListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *RegionListRequest) (*RegionListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionName", req.RegionName)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	result := &regionListRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	var regionList []RegionListRegionListResponse
	for _, region := range result.RegionList {
		regionList = append(regionList, RegionListRegionListResponse{
			RegionParent: region.RegionParent,
			RegionId:     region.RegionID,
			RegionType:   region.RegionType,
			RegionName:   region.RegionName,
			IsMultiZones: region.IsMultiZones,
			ZoneList:     region.ZoneList,
		})
	}
	return &RegionListResponse{
		RegionList: regionList,
	}, nil
}

type regionListRealResponse struct {
	RegionList []struct {
		RegionParent string   `json:"regionParent"`
		RegionID     string   `json:"regionID"`
		RegionType   string   `json:"regionType"`
		RegionName   string   `json:"regionName"`
		IsMultiZones bool     `json:"isMultiZones"`
		ZoneList     []string `json:"zoneList"`
	} `json:"regionList"`
}

type RegionListRequest struct {
	RegionName string
}

type RegionListRegionListResponse struct {
	RegionParent string
	RegionId     string
	RegionType   string
	RegionName   string
	IsMultiZones bool
	ZoneList     []string
}

type RegionListResponse struct {
	RegionList []RegionListRegionListResponse
}
