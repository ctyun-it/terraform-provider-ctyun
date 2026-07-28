package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcListCDANetworkApi
/* 查询CDA网络实例 */
type EcEcListCDANetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListCDANetworkApi(client *core.CtyunClient) *EcEcListCDANetworkApi {
	return &EcEcListCDANetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/cda-instance/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListCDANetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcListCDANetworkRequest) (*EcEcListCDANetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcListCDANetworkRequest
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
	var resp EcEcListCDANetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListCDANetworkRequest struct {
	EcID       string `json:"ecID"`       /*  云间高速实例ID  */
	CgwID      string `json:"cgwID"`      /*  云网关ID  */
	CdaID      string `json:"cdaID"`      /*  云专线ID  */
	InstanceID string `json:"instanceID"` /*  指定网络实例ID  */
	Status     string `json:"status"`     /*  指定状态描述<br/>取值包括：<br/>'creating'：加载中<br/>'running'：已连接<br/>'removing'：卸载中<br/>'flushing'：路由待更新<br/>'error'：失败  */
}

type EcEcListCDANetworkResponse struct {
	StatusCode  *int32  `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述   */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述    */
	TraceID     *string `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *string `json:"returnObj"`   /*  返回参数  */
}
