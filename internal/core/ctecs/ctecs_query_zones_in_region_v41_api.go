package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryZonesInRegionV41Api
/* 查询单个资源池的可用区信息
 */type CtecsQueryZonesInRegionV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryZonesInRegionV41Api(client *core.CtyunClient) *CtecsQueryZonesInRegionV41Api {
	return &CtecsQueryZonesInRegionV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/region/get-zones",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryZonesInRegionV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryZonesInRegionV41Request) (*CtecsQueryZonesInRegionV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryZonesInRegionV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryZonesInRegionV41Request struct {
	RegionID string /*  资源池ID  */
}

type CtecsQueryZonesInRegionV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                       `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                       `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQueryZonesInRegionV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryZonesInRegionV41ReturnObjResponse struct {
	ZoneList []*CtecsQueryZonesInRegionV41ReturnObjZoneListResponse `json:"zoneList"` /*  可用区列表  */
}

type CtecsQueryZonesInRegionV41ReturnObjZoneListResponse struct {
	Name          string `json:"name,omitempty"`          /*  可用区名称，其他需要可用区参数的接口需要依赖该名称结果  */
	AzDisplayName string `json:"azDisplayName,omitempty"` /*  可用区展示名  */
}
