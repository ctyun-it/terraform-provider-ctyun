package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCreateApi
/* 创建云间高速 */
type EcEcCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCreateApi(client *core.CtyunClient) *EcEcCreateApi {
	return &EcEcCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/express-connect/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCreateApi) Do(ctx context.Context, credential core.Credential, req *EcEcCreateRequest) (*EcEcCreateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCreateRequest
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
	var resp EcEcCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCreateRequest struct {
	EcName        string  `json:"ecName"`                  /*  名称  */
	EcDescription *string `json:"ecDescription,omitempty"` /*  描述信息  */
}

type EcEcCreateResponse struct {
	StatusCode  *int32                       `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Description *string                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                      `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcCreateResponseReturnObj `json:"returnObj"`   /*  详细信息 */
	Message     *string                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */

}
type EcEcCreateResponseReturnObj struct {
	EcID          *string `json:"ecID"`          /*  云间高速实例ID  */
	EcName        *string `json:"ecName"`        /*  名称  */
	EcDescription *string `json:"ecDescription"` /*  描述信息  */
	Status        *int32  `json:"status"`        /*  运行状态，<br/>取值范围:<br/>0:创建中<br/>2:运行中<br/>18:删除中<br/>21:设置中<br/>22:更新带宽中<br/>24:更新中  */
	CreateDate    *string `json:"createDate"`    /*  创建时间  */
	Error         *string `json:"error"`         /*  业务细分码，为product.module.code三段式码  */
}
