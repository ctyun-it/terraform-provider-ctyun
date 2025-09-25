package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryProductsInRegionV41Api
/* 查询一个资源池支持的云产品信息列表，以及云产品的产品特性信息。
 */type CtecsQueryProductsInRegionV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryProductsInRegionV41Api(client *core.CtyunClient) *CtecsQueryProductsInRegionV41Api {
	return &CtecsQueryProductsInRegionV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/region/get-products",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryProductsInRegionV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryProductsInRegionV41Request) (*CtecsQueryProductsInRegionV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryProductsInRegionV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryProductsInRegionV41Request struct {
	RegionID string /*  资源池ID  */
}

type CtecsQueryProductsInRegionV41Response struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQueryProductsInRegionV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryProductsInRegionV41ReturnObjResponse struct {
	AzList []*CtecsQueryProductsInRegionV41ReturnObjAzListResponse `json:"azList"` /*  az分区列表  */
}

type CtecsQueryProductsInRegionV41ReturnObjAzListResponse struct {
	AzName        string                                                       `json:"azName,omitempty"`        /*  可用区名称  */
	AzDisplayName string                                                       `json:"azDisplayName,omitempty"` /*  可用区展示名  */
	Details       *CtecsQueryProductsInRegionV41ReturnObjAzListDetailsResponse `json:"details"`                 /*  可用区详细信息  */
}

type CtecsQueryProductsInRegionV41ReturnObjAzListDetailsResponse struct {
	StorageType []*CtecsQueryProductsInRegionV41ReturnObjAzListDetailsStorageTypeResponse `json:"storageType"` /*  不同az可用区的存储类型  */
}

type CtecsQueryProductsInRegionV41ReturnObjAzListDetailsStorageTypeResponse struct {
	RawType string `json:"type,omitempty"` /*  存储类型  */
	Name    string `json:"name,omitempty"` /*  类型名称  */
}
