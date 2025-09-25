package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListRegionsV41Api
/* 查询租户可见的资源池列表。
 */type CtecsListRegionsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListRegionsV41Api(client *core.CtyunClient) *CtecsListRegionsV41Api {
	return &CtecsListRegionsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/region/list-regions",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListRegionsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListRegionsV41Request) (*CtecsListRegionsV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.RegionName != "" {
		ctReq.AddParam("regionName", req.RegionName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListRegionsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListRegionsV41Request struct {
	RegionName string /*  资源池名称  */
}

type CtecsListRegionsV41Response struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsListRegionsV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsListRegionsV41ReturnObjResponse struct {
	RegionList []*CtecsListRegionsV41ReturnObjRegionListResponse `json:"regionList"` /*  资源池对象  */
}

type CtecsListRegionsV41ReturnObjRegionListResponse struct {
	RegionID         string   `json:"regionID,omitempty"`     /*  资源池ID	  */
	RegionParent     string   `json:"regionParent,omitempty"` /*  资源池所属省份	  */
	RegionName       string   `json:"regionName,omitempty"`   /*  资源池名称	  */
	RegionType       string   `json:"regionType,omitempty"`   /*  资源池类型	  */
	IsMultiZones     *bool    `json:"isMultiZones"`           /*  是否多可用区资源池	  */
	ZoneList         []string `json:"zoneList"`               /*  可用区列表	  */
	RegionCode       string   `json:"regionCode,omitempty"`   /*  地域编号  */
	OpenapiAvailable *bool    `json:"openapiAvailable"`       /*  是否支持通过OpenAPI访问  */
}
