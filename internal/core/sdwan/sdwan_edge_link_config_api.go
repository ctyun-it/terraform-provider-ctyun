package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanEdgeLinkConfigApi
/* EDGE链路配置 */
type SdwanEdgeLinkConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanEdgeLinkConfigApi(client *core.CtyunClient) *SdwanEdgeLinkConfigApi {
	return &SdwanEdgeLinkConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge/link-config",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanEdgeLinkConfigApi) Do(ctx context.Context, credential core.Credential, req *SdwanEdgeLinkConfigRequest) (*SdwanEdgeLinkConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanEdgeLinkConfigRequest
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
	var resp SdwanEdgeLinkConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanEdgeLinkConfigRequest struct {
	EdgeID string                             `json:"edgeID"`          /*  edge id  */
	Links  []*SdwanEdgeLinkConfigLinksRequest `json:"links,omitempty"` /*  链路配置 必填  */
}

type SdwanEdgeLinkConfigLinksRequest struct {
	Port         *string `json:"port,omitempty"`         /*  本参数表示端口名称<br/>取值范围:<br/>WAN1:WAN1<br/>WAN2:WAN2<br/>LAN:LAN<br/>LTE:LTE  */
	TransportNet *string `json:"transportNet,omitempty"` /*  本参数表示传输网络类型<br/><br/>取值范围:<br/>internet:internet<br/>mpls:mpls<br/>lte:lte  */
	Level        *string `json:"level,omitempty"`        /*  本参数表示主备链路配置<br/>取值范围：<br/>master:主<br/>slave:备  */
}

type SdwanEdgeLinkConfigResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
