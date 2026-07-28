package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcDeleteVPCNetworkApi
/* 删除VPC网络实例 */
type EcEcDeleteVPCNetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcDeleteVPCNetworkApi(client *core.CtyunClient) *EcEcDeleteVPCNetworkApi {
	return &EcEcDeleteVPCNetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/vpc-instance/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcDeleteVPCNetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcDeleteVPCNetworkRequest) (*EcEcDeleteVPCNetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcDeleteVPCNetworkRequest
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
	var resp EcEcDeleteVPCNetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcDeleteVPCNetworkRequest struct {
	VpcID string `json:"vpcID"` /*  vpc ID  */
	EcID  string `json:"ecID"`  /*  云间高速实例ID  */
}

type EcEcDeleteVPCNetworkResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcDeleteVPCNetworkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcEcDeleteVPCNetworkReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /*  操作日志id  */
}
