package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanEdgeSaftyFlowBandwidthApi
/* 安全引流带宽创建 */
type SdwanCreateSdwanEdgeSaftyFlowBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanEdgeSaftyFlowBandwidthApi(client *core.CtyunClient) *SdwanCreateSdwanEdgeSaftyFlowBandwidthApi {
	return &SdwanCreateSdwanEdgeSaftyFlowBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-safty-flow-bandwidth/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanEdgeSaftyFlowBandwidthApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanEdgeSaftyFlowBandwidthRequest) (*SdwanCreateSdwanEdgeSaftyFlowBandwidthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanEdgeSaftyFlowBandwidthRequest
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
	var resp SdwanCreateSdwanEdgeSaftyFlowBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanEdgeSaftyFlowBandwidthRequest struct {
	Account    string `json:"account"`    /*  用户账号  */
	Bandwidth  string `json:"bandwidth"`  /*  带宽  */
	StartTime  string `json:"startTime"`  /*  起始时间  */
	ExpireTime string `json:"expireTime"` /*  到期时间  */
}

type SdwanCreateSdwanEdgeSaftyFlowBandwidthResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
}
