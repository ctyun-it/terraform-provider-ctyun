package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateSdwanEdgeSaftyFlowBandwidthApi
/* 安全引流带宽修改 */
type SdwanUpdateSdwanEdgeSaftyFlowBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateSdwanEdgeSaftyFlowBandwidthApi(client *core.CtyunClient) *SdwanUpdateSdwanEdgeSaftyFlowBandwidthApi {
	return &SdwanUpdateSdwanEdgeSaftyFlowBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-safty-flow-bandwidth/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateSdwanEdgeSaftyFlowBandwidthApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateSdwanEdgeSaftyFlowBandwidthRequest) (*SdwanUpdateSdwanEdgeSaftyFlowBandwidthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateSdwanEdgeSaftyFlowBandwidthRequest
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
	var resp SdwanUpdateSdwanEdgeSaftyFlowBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateSdwanEdgeSaftyFlowBandwidthRequest struct {
	InstanceID string  `json:"instanceID"`           /*  实例id  */
	Bandwidth  *string `json:"bandwidth,omitempty"`  /*  带宽（带宽和时间一次只能改一个）  */
	StartTime  *string `json:"startTime,omitempty"`  /*  起始时间（起始时间和到期时间一起填）  */
	ExpireTime *string `json:"expireTime,omitempty"` /*  到期时间（起始时间和到期时间一起填）  */
}

type SdwanUpdateSdwanEdgeSaftyFlowBandwidthResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
}
