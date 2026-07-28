package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCreateGatewayApi
/* 创建云网关 */
type EcEcCreateGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCreateGatewayApi(client *core.CtyunClient) *EcEcCreateGatewayApi {
	return &EcEcCreateGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cloud-gateway/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCreateGatewayApi) Do(ctx context.Context, credential core.Credential, req *EcEcCreateGatewayRequest) (*EcEcCreateGatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCreateGatewayRequest
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
	var resp EcEcCreateGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCreateGatewayRequest struct {
	CgwName     string  `json:"cgwName"`                  /*  云网关名称  */
	Description *string `json:"cgwDescription,omitempty"` /*  云网关描述  */
	Region      *int32  `json:"region,omitempty"`         /*  地域信息，取值如下<br/>1：中国大陆（默认）<br/>2:亚太  */
	DcName      string  `json:"dcName"`                   /*  资源池名称  */
	DcID        string  `json:"dcID"`                     /*  资源池ID信息  */
	EcID        string  `json:"ecID"`                     /*  云间高速实例ID  */
	DcType      string  `json:"dcType"`                   /*  资源池类型，<br/>取值范围:<br/>'CNP':CNP资源池<br/>'MAZ':MAZ资源池<br/>'PRVT':私有云资源池  */
}

type EcEcCreateGatewayResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                             `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcCreateGatewayReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcCreateGatewayReturnObjResponse struct {
	CgwID      *string `json:"cgwID"`      /*  云网关实例ID  */
	CreateDate *string `json:"createDate"` /*  创建时间  */
}
