package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcFlowsessionUpdateTargetApi
/* 修改流量会话目的
 */type CtvpcFlowsessionUpdateTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcFlowsessionUpdateTargetApi(client *core.CtyunClient) *CtvpcFlowsessionUpdateTargetApi {
	return &CtvpcFlowsessionUpdateTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flowsession/update-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcFlowsessionUpdateTargetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcFlowsessionUpdateTargetRequest) (*CtvpcFlowsessionUpdateTargetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcFlowsessionUpdateTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcFlowsessionUpdateTargetRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  区域ID  */
	FlowSessionID string `json:"flowSessionID,omitempty"` /*  流量镜像会话 ID  */
	TargetPortID  string `json:"targetPortID,omitempty"`  /*  目的网卡 ID  */
}

type CtvpcFlowsessionUpdateTargetResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcFlowsessionUpdateTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcFlowsessionUpdateTargetReturnObjResponse struct {
	FlowSessionID *string `json:"flowSessionID,omitempty"` /*  流量镜像会话 ID  */
}
