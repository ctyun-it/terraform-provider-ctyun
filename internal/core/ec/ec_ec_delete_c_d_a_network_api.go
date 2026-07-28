package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcDeleteCDANetworkApi
/* 删除CDA网络实例 */
type EcEcDeleteCDANetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcDeleteCDANetworkApi(client *core.CtyunClient) *EcEcDeleteCDANetworkApi {
	return &EcEcDeleteCDANetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cda-instance/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcDeleteCDANetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcDeleteCDANetworkRequest) (*EcEcDeleteCDANetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcDeleteCDANetworkRequest
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
	var resp EcEcDeleteCDANetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcDeleteCDANetworkRequest struct {
	InstanceID string  `json:"instanceID"`      /*  指定网络实例ID  */
	CdaID      *string `json:"cdaID,omitempty"` /*  专线ID  */
}

type EcEcDeleteCDANetworkResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcDeleteCDANetworkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcDeleteCDANetworkReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /*  操作日志id  */
}
