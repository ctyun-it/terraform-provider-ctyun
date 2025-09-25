package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsRefundSfsApi
/* 退订文件系统
 */type SfsSfsRefundSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsRefundSfsApi(client *core.CtyunClient) *SfsSfsRefundSfsApi {
	return &SfsSfsRefundSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/refund-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsRefundSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsRefundSfsRequest) (*SfsSfsRefundSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsRefundSfsRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsRefundSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsRefundSfsRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	SfsUID      string `json:"sfsUID,omitempty"`      /*  弹性文件功能系统唯一 ID  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
}

type SfsSfsRefundSfsResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800为成功，900为失败/订单处理中)  */
	Message     string                              `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                              `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsRefundSfsReturnObjResponse   `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
	ErrorDetail *SfsSfsRefundSfsErrorDetailResponse `json:"errorDetail"` /*  错误明细。一般情况下，会对订单侧(bss)的弹性文件订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。其他订单侧(bss)的错误，以sfs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
}

type SfsSfsRefundSfsReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID"` /*  退订订单号，可以使用该订单号确认资源的最终退订状态  */
	MasterOrderNO string `json:"masterOrderNO"` /*  订单号  */
	RegionID      string `json:"regionID"`      /*  资源所属资源池 ID  */
}

type SfsSfsRefundSfsErrorDetailResponse struct {
	BssErrCode       string `json:"bssErrCode"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中  */
	BssErrMsg        string `json:"bssErrMsg"`        /*  bss错误信息，包含于bss格式化JSON错误信息中  */
	BssOrigErr       string `json:"bssOrigErr"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息  */
	BssErrPrefixHint string `json:"bssErrPrefixHint"` /*  bss格式化JSON错误信息的前置提示信息  */
}
