package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQuerySummaryInRegionV41Api
/* 查询资源池概况，比如地域，多az信息，支持的cpu架构，资源池占用类型，资源池版本信息等
 */type CtecsQuerySummaryInRegionV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQuerySummaryInRegionV41Api(client *core.CtyunClient) *CtecsQuerySummaryInRegionV41Api {
	return &CtecsQuerySummaryInRegionV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/region/get-summary",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQuerySummaryInRegionV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQuerySummaryInRegionV41Request) (*CtecsQuerySummaryInRegionV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQuerySummaryInRegionV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQuerySummaryInRegionV41Request struct {
	RegionID string /*  资源池ID  */
}

type CtecsQuerySummaryInRegionV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQuerySummaryInRegionV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQuerySummaryInRegionV41ReturnObjResponse struct {
	RegionID         string   `json:"regionID,omitempty"`      /*  资源池ID  */
	RegionParent     string   `json:"regionParent,omitempty"`  /*  资源池所属省份  */
	RegionName       string   `json:"regionName,omitempty"`    /*  资源池名称  */
	RegionType       string   `json:"regionType,omitempty"`    /*  资源池类型  */
	IsMultiZones     *bool    `json:"isMultiZones"`            /*  是否多可用区资源池  */
	ZoneList         []string `json:"zoneList"`                /*  可用区列表  */
	CpuArches        []string `json:"cpuArches"`               /*  资源池cpu架构信息  */
	RegionVersion    string   `json:"regionVersion,omitempty"` /*  资源池版本  */
	Dedicated        *bool    `json:"dedicated"`               /*  是否是专属资源池，账户可能访问的是一个自己可见的专属资源池  */
	Province         string   `json:"province,omitempty"`      /*  省份  */
	City             string   `json:"city,omitempty"`          /*  城市  */
	OpenapiAvailable *bool    `json:"openapiAvailable"`        /*  是否支持通过OpenAPI访问  */
}
