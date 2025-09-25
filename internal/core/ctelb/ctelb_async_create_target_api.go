package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbAsyncCreateTargetApi
/* 创建后端服务
 */type CtelbAsyncCreateTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbAsyncCreateTargetApi(client *core.CtyunClient) *CtelbAsyncCreateTargetApi {
	return &CtelbAsyncCreateTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/async-create-vm",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbAsyncCreateTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbAsyncCreateTargetRequest) (*CtelbAsyncCreateTargetResponse, error) {
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
	var resp CtelbAsyncCreateTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbAsyncCreateTargetRequest struct {
	ClientToken   string                                  `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID      string                                  `json:"regionID,omitempty"`      /*  区域ID  */
	TargetGroupID string                                  `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	Targets       []*CtelbAsyncCreateTargetTargetsRequest `json:"targets"`                 /*  后端服务主机  */
}

type CtelbAsyncCreateTargetTargetsRequest struct {
	InstanceID   string `json:"instanceID,omitempty"`   /*  后端服务主机 id  */
	ProtocolPort int32  `json:"protocolPort,omitempty"` /*  后端服务监听端口，1-65535  */
	InstanceType string `json:"instanceType,omitempty"` /*  后端服务主机类型，仅支持vm类型  */
	Weight       int32  `json:"weight,omitempty"`       /*  后端服务主机权重: 1 - 256  */
	Address      string `json:"address,omitempty"`      /*  后端服务主机主网卡所在的 IP  */
}

type CtelbAsyncCreateTargetResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbAsyncCreateTargetReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbAsyncCreateTargetReturnObjResponse struct{}
