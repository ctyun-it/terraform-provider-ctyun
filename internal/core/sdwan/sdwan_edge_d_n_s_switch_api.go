package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanEdgeDNSSwitchApi
/* 智能网关 dns开关 */
type SdwanEdgeDNSSwitchApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanEdgeDNSSwitchApi(client *core.CtyunClient) *SdwanEdgeDNSSwitchApi {
	return &SdwanEdgeDNSSwitchApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-dns/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanEdgeDNSSwitchApi) Do(ctx context.Context, credential core.Credential, req *SdwanEdgeDNSSwitchRequest) (*SdwanEdgeDNSSwitchResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanEdgeDNSSwitchRequest
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
	var resp SdwanEdgeDNSSwitchResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanEdgeDNSSwitchRequest struct {
	EdgeID    string `json:"edgeID"`    /*  智能网关ID  */
	DnsSwitch string `json:"dnsSwitch"` /*  DNS开关  */
}

type SdwanEdgeDNSSwitchResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
