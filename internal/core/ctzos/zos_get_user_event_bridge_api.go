package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetUserEventBridgeApi
/* 用于查询对象存储用户级事件总线状态。
 */type ZosGetUserEventBridgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetUserEventBridgeApi(client *core.CtyunClient) *ZosGetUserEventBridgeApi {
	return &ZosGetUserEventBridgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-user-event-bridge",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetUserEventBridgeApi) Do(ctx context.Context, credential core.Credential, req *ZosGetUserEventBridgeRequest) (*ZosGetUserEventBridgeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetUserEventBridgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetUserEventBridgeRequest struct {
	RegionID string /*  资源池 ID  */
}

type ZosGetUserEventBridgeResponse struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回码<br>取值范围：800 成功  */
	Message     string                                  `json:"message,omitempty"`     /*  状态描述  */
	Description string                                  `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetUserEventBridgeReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                  `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetUserEventBridgeReturnObjResponse struct {
	UserEventBridgeEnabled *bool `json:"userEventBridgeEnabled"` /*  对象存储用户事件总线开启状态，false为关闭，ture为开启  */
}
