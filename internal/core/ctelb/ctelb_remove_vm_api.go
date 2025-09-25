package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbRemoveVmApi
/* 删除后端服务
 */type CtelbRemoveVmApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbRemoveVmApi(client *core.CtyunClient) *CtelbRemoveVmApi {
	return &CtelbRemoveVmApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/remove-vm",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbRemoveVmApi) Do(ctx context.Context, credential core.Credential, req *CtelbRemoveVmRequest) (*CtelbRemoveVmResponse, error) {
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
	var resp CtelbRemoveVmResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbRemoveVmRequest struct {
	RegionID      string   `json:"regionID,omitempty"`      /*  区域ID  */
	TargetGroupID string   `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	TargetIDs     []string `json:"targetIDs"`               /*  后端服务 ID 列表  */
}

type CtelbRemoveVmResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
