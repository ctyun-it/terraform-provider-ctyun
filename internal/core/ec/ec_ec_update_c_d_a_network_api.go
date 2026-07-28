package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcUpdateCDANetworkApi
/* 修改CDA网络实例 */
type EcEcUpdateCDANetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcUpdateCDANetworkApi(client *core.CtyunClient) *EcEcUpdateCDANetworkApi {
	return &EcEcUpdateCDANetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cda-instance/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcUpdateCDANetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcUpdateCDANetworkRequest) (*EcEcUpdateCDANetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcUpdateCDANetworkRequest
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
	var resp EcEcUpdateCDANetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcUpdateCDANetworkRequest struct {
	InstanceID string      `json:"instanceID"` /*  指定网络实例ID  */
	CdaInfo    interface{} `json:"cdaInfo"`    /*  云专线信息，json格式，仅支持更新子网  */
}

type EcEcUpdateCDANetworkResponse struct {
	StatusCode  *int32  `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string `json:"traceID"`     /*  链路追踪ID  */
}
