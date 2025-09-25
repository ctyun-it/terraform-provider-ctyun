package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbDeleteAccessControlApi
/* 删除访问控制策略
 */type CtelbDeleteAccessControlApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbDeleteAccessControlApi(client *core.CtyunClient) *CtelbDeleteAccessControlApi {
	return &CtelbDeleteAccessControlApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/delete-access-control",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbDeleteAccessControlApi) Do(ctx context.Context, credential core.Credential, req *CtelbDeleteAccessControlRequest) (*CtelbDeleteAccessControlResponse, error) {
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
	var resp CtelbDeleteAccessControlResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbDeleteAccessControlRequest struct {
	ClientToken     string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID        string `json:"regionID,omitempty"`        /*  区域ID  */
	ID              string `json:"ID,omitempty"`              /*  访问控制ID, 该字段后续废弃  */
	AccessControlID string `json:"accessControlID,omitempty"` /*  访问控制ID, 推荐使用该字段, 当同时使用 ID 和 accessControlID 时，优先使用 accessControlID  */
}

type CtelbDeleteAccessControlResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
