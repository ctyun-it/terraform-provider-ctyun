package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsRenewEbsApi
/* 支持包周期的数据盘续订。
 */type EbsRenewEbsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsRenewEbsApi(client *core.CtyunClient) *EbsRenewEbsApi {
	return &EbsRenewEbsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/renew-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsRenewEbsApi) Do(ctx context.Context, credential core.Credential, req *EbsRenewEbsRequest) (*EbsRenewEbsResponse, error) {
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
	var resp EbsRenewEbsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsRenewEbsRequest struct {
	DiskID    string  `json:"diskID,omitempty"`    /*  云硬盘ID。  */
	RegionID  *string `json:"regionID,omitempty"`  /*  资源池ID。如本地语境支持保存regionID，那么建议传递。  */
	CycleType string  `json:"cycleType,omitempty"` /*  包周期类型，取值范围：
	●year：包年。
	●month：包月。  */
	CycleCount int32 `json:"cycleCount"` /*  包周期数。
	周期为年（year）时，最大支持续订3年。
	周期为月（month）时，最大支持续订36个月。  */
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。  */
}

type EbsRenewEbsResponse struct {
	StatusCode  int32                           `json:"statusCode"`            /*  返回状态码(800为成功，900为处理中/失败)。  */
	Message     *string                         `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                         `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsRenewEbsReturnObjResponse   `json:"returnObj"`             /*  返回结构体。  */
	ErrorCode   *string                         `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                         `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
	ErrorDetail *EbsRenewEbsErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。<br> 其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
}

type EbsRenewEbsReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID,omitempty"` /*  订单ID。调用方在拿到masterOrderID之后，<br/>在若干错误情况下，可以使用masterOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO *string `json:"masterOrderNO,omitempty"` /*  订单号。  */
	RegionID      *string `json:"regionID,omitempty"`      /*  资源池ID。  */
}

type EbsRenewEbsErrorDetailResponse struct {
	BssErrCode       *string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中。  */
	BssErrMsg        *string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中。  */
	BssOrigErr       *string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息。  */
	BssErrPrefixHint *string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息。  */
}
