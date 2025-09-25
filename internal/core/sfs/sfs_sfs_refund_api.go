package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsRefundApi
/* 退订文件系统
 */type SfsSfsRefundApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsRefundApi(client *core.CtyunClient) *SfsSfsRefundApi {
	return &SfsSfsRefundApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/refund",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsRefundApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsRefundRequest) (*SfsSfsRefundResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsRefundRequest
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
	var resp SfsSfsRefundResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsRefundRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	ResourceID  string `json:"resourceID,omitempty"`  /*  文件系统资源 ID  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
}

type SfsSfsRefundResponse struct {
	StatusCode  int32                          `json:"statusCode"`  /*  返回状态码(800为成功，900为订单处理中/失败)  */
	Message     string                         `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                         `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SfsSfsRefundReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                         `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
}

type SfsSfsRefundReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID"` /*  退订订单号，可以使用该订单号确认资源的最终退订状态  */
	MasterOrderNO string `json:"masterOrderNO"` /*  订单号  */
	RegionID      string `json:"regionID"`      /*  资源池 ID  */
}
