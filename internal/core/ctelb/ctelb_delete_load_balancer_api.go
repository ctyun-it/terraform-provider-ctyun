package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbDeleteLoadBalancerApi
/* 删除负载均衡实例
 */type CtelbDeleteLoadBalancerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbDeleteLoadBalancerApi(client *core.CtyunClient) *CtelbDeleteLoadBalancerApi {
	return &CtelbDeleteLoadBalancerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/delete-loadbalancer",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbDeleteLoadBalancerApi) Do(ctx context.Context, credential core.Credential, req *CtelbDeleteLoadBalancerRequest) (*CtelbDeleteLoadBalancerResponse, error) {
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
	var resp CtelbDeleteLoadBalancerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbDeleteLoadBalancerRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	ProjectID   string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为'0'  */
	ID          string `json:"ID,omitempty"`          /*  负载均衡ID, 该字段后续废弃  */
	ElbID       string `json:"elbID,omitempty"`       /*  负载均衡ID, 推荐使用该字段, 当同时使用 ID 和 elbID 时，优先使用 elbID  */
}

type CtelbDeleteLoadBalancerResponse struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbDeleteLoadBalancerReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbDeleteLoadBalancerReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  负载均衡ID  */
}
