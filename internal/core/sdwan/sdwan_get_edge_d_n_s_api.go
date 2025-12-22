package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeDNSApi
/* 查找智能网关 dns地址信息 */
type SdwanGetEdgeDNSApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeDNSApi(client *core.CtyunClient) *SdwanGetEdgeDNSApi {
	return &SdwanGetEdgeDNSApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-dns/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeDNSApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeDNSRequest) (*SdwanGetEdgeDNSResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeDNSResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeDNSRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanGetEdgeDNSResponse struct {
	StatusCode  int32                             `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                           `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeDNSReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                           `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeDNSReturnObjResponse struct {
	Result       []*SdwanGetEdgeDNSReturnObjResultResponse `json:"result"`       /*  查询edge dns  */
	TotalCount   int32                                     `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                     `json:"currentCount"` /*  页码  */
}

type SdwanGetEdgeDNSReturnObjResultResponse struct {
	DNS       *string `json:"DNS"`       /*  DNS  */
	DNSSwitch *string `json:"DNSSwitch"` /*  本参数表示DNS开关<br/>取值范围：<br/>true：开启。<br/>false: 关闭。  */
}
