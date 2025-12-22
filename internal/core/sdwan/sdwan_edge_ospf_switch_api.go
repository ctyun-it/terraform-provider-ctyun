package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanEdgeOspfSwitchApi
/* 智能网关ospf开关 */
type SdwanEdgeOspfSwitchApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanEdgeOspfSwitchApi(client *core.CtyunClient) *SdwanEdgeOspfSwitchApi {
	return &SdwanEdgeOspfSwitchApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-ospf-switch/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanEdgeOspfSwitchApi) Do(ctx context.Context, credential core.Credential, req *SdwanEdgeOspfSwitchRequest) (*SdwanEdgeOspfSwitchResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanEdgeOspfSwitchRequest
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
	var resp SdwanEdgeOspfSwitchResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanEdgeOspfSwitchRequest struct {
	EdgeID     string `json:"edgeID"`     /*  智能网关ID  */
	OspfSwitch string `json:"ospfSwitch"` /*  本参数表示ospf开关<br/>取值范围：<br/>true： 开启<br/>false：关闭  */
}

type SdwanEdgeOspfSwitchResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
