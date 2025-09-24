package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreateTargetApi
/* 创建后端服务
 */type CtelbCreateTargetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreateTargetApi(client *core.CtyunClient) *CtelbCreateTargetApi {
	return &CtelbCreateTargetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-target",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreateTargetApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreateTargetRequest) (*CtelbCreateTargetResponse, error) {
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
	var resp CtelbCreateTargetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreateTargetRequest struct {
	ClientToken   string `json:"clientToken,omitempty"`   /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	InstanceIP    string `json:"instanceIP,omitempty"`    /*  后端服务 ip  */
	RegionID      string `json:"regionID,omitempty"`      /*  区域ID  */
	TargetGroupID string `json:"targetGroupID,omitempty"` /*  后端服务组ID  */
	InstanceType  string `json:"instanceType,omitempty"`  /*  实例类型。取值范围：VM、BM、ECI、IP  */
	InstanceID    string `json:"instanceID,omitempty"`    /*  实例ID  */
	ProtocolPort  int32  `json:"protocolPort,omitempty"`  /*  协议端口。取值范围：1-65535  */
	Weight        int32  `json:"weight,omitempty"`        /*  权重。取值范围：1-256，默认为100  */
}

type CtelbCreateTargetResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbCreateTargetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbCreateTargetReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  后端服务组ID  */
}
