package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcUpdateGatewayApi
/* 修改云网关描述 */
type EcEcUpdateGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcUpdateGatewayApi(client *core.CtyunClient) *EcEcUpdateGatewayApi {
	return &EcEcUpdateGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cloud-gateway/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcUpdateGatewayApi) Do(ctx context.Context, credential core.Credential, req *EcEcUpdateGatewayRequest) (*EcEcUpdateGatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcUpdateGatewayRequest
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
	var resp EcEcUpdateGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcUpdateGatewayRequest struct {
	CgwID          string  `json:"cgwID"`                    /*  云网关实例ID  */
	CgwName        *string `json:"cgwName,omitempty"`        /*  云网关名称  */
	CgwDescription *string `json:"cgwDescription,omitempty"` /*  云网关描述  */
}

type EcEcUpdateGatewayResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                             `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcUpdateGatewayReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcUpdateGatewayReturnObjResponse struct {
	CgwID          *string `json:"cgwID"`          /*  云网关实例ID  */
	CgwName        *string `json:"cgwName"`        /*  云网关名称  */
	CgwDescription *string `json:"cgwDescription"` /*  云网关描述  */
}
