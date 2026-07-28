package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCheckSDWANApi
/* SDWAN入云校验网段，针对已经加入云间高速的SDWAN，当SDWAN发生子网变动，检验变动后的网段是否和云间高速其他网段冲突。 */
type EcEcCheckSDWANApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCheckSDWANApi(client *core.CtyunClient) *EcEcCheckSDWANApi {
	return &EcEcCheckSDWANApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/sdwan/check-cidr",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCheckSDWANApi) Do(ctx context.Context, credential core.Credential, req *EcEcCheckSDWANRequest) (*EcEcCheckSDWANResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCheckSDWANRequest
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
	var resp EcEcCheckSDWANResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCheckSDWANRequest struct {
	SdwanID string   `json:"sdwanID"` /*  sdwan ID  */
	CIDRs   []string `json:"CIDRs"`   /*  sdwan新增的v4网段，与v6CIDRs不能同时为空  */
	V6CIDRs []string `json:"v6CIDRs"` /*  sdwan新增的v6网段，与CIDRs不能同时为空  */
}

type EcEcCheckSDWANResponse struct {
	TraceID     *string                          `json:"traceID"`     /*  链路追踪ID  */
	StatusCode  *int32                           `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcCheckSDWANReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcCheckSDWANReturnObjResponse struct {
	IsConflict *bool `json:"isConflict"` /*  sdwan网段是否冲突，返回ture则冲突  */
}
