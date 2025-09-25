package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbDisassociateEipFromLoadBalancerApi
/* 解绑弹性 IP
 */type CtelbDisassociateEipFromLoadBalancerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbDisassociateEipFromLoadBalancerApi(client *core.CtyunClient) *CtelbDisassociateEipFromLoadBalancerApi {
	return &CtelbDisassociateEipFromLoadBalancerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/disassociate-eip-from",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbDisassociateEipFromLoadBalancerApi) Do(ctx context.Context, credential core.Credential, req *CtelbDisassociateEipFromLoadBalancerRequest) (*CtelbDisassociateEipFromLoadBalancerResponse, error) {
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
	var resp CtelbDisassociateEipFromLoadBalancerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbDisassociateEipFromLoadBalancerRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	ID          string `json:"ID,omitempty"`          /*  负载均衡ID, 该字段后续废弃  */
	ElbID       string `json:"elbID,omitempty"`       /*  负载均衡ID, 推荐使用该字段, 当同时使用 ID 和 elbID 时，优先使用 elbID  */
	EipID       string `json:"eipID,omitempty"`       /*  弹性公网IP的ID  */
}

type CtelbDisassociateEipFromLoadBalancerResponse struct {
	StatusCode  int32                                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbDisassociateEipFromLoadBalancerReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbDisassociateEipFromLoadBalancerReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  负载均衡ID  */
}
