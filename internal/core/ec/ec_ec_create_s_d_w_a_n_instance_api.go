package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCreateSDWANInstanceApi
/* 添加sdwan网络实例 */
type EcEcCreateSDWANInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCreateSDWANInstanceApi(client *core.CtyunClient) *EcEcCreateSDWANInstanceApi {
	return &EcEcCreateSDWANInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/sdwan-instance/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCreateSDWANInstanceApi) Do(ctx context.Context, credential core.Credential, req *EcEcCreateSDWANInstanceRequest) (*EcEcCreateSDWANInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCreateSDWANInstanceRequest
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
	var resp EcEcCreateSDWANInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCreateSDWANInstanceRequest struct {
	EcID       string `json:"ecID"`                 /*  云间高速ID  */
	CgwID      string `json:"cgwID"`                /*  云网关ID  */
	SdwanID    string `json:"sdwanID"`              /*  sdwan ID  */
	RtbID      string `json:"rtbID"`                /*  云网关默认路由表ID  */
	Weights    *int32 `json:"weights,omitempty"`    /*  权重，sdwan默认60，无冗余实例则不传  */
	RouteLearn *int32 `json:"routeLearn,omitempty"` /*  路由学习开关，开启后云网关自动学习网络实例路由<br/>取值范围:<br/>1:学习<br/>0:不学习<br/>默认学习  */
	RouteSync  *int32 `json:"routeSync,omitempty"`  /*  路由同步开关，开启后云网关路由自动同步到网络实例<br/>取值范围:<br/>1:同步<br/>0:不同步<br/>默认同步  */
}

type EcEcCreateSDWANInstanceResponse struct {
	TraceID     *string                                   `json:"traceID"`     /*  链路追踪ID  */
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcCreateSDWANInstanceReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcCreateSDWANInstanceReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /*  操作日志id  */
}
