package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcQueryToCyclePriceApi
/* 查询转包周期价格
 */type CtvpcQueryToCyclePriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcQueryToCyclePriceApi(client *core.CtyunClient) *CtvpcQueryToCyclePriceApi {
	return &CtvpcQueryToCyclePriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/order/query-to-cycle-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcQueryToCyclePriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcQueryToCyclePriceRequest) (*CtvpcQueryToCyclePriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcQueryToCyclePriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcQueryToCyclePriceRequest struct {
	ResourceID   string `json:"resourceID,omitempty"`   /*  资源ID  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域ID  */
	CycleType    string `json:"cycleType,omitempty"`    /*  周期类型  */
	CycleCount   int32  `json:"cycleCount"`             /*  周期数,cycleType=month时,1<=cycleCount<=12,cycleType=year时,1<=cycleCount<=3  */
}

type CtvpcQueryToCyclePriceResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	Message     *string                                  `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtvpcQueryToCyclePriceReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcQueryToCyclePriceReturnObjResponse struct {
	TotalPrice float32 `json:"totalPrice"` /*  总价  */
	FinalPrice float32 `json:"finalPrice"` /*  最终价格  */
}
