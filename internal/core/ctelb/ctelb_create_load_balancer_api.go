package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbCreateLoadBalancerApi
/* 创建负载均衡实例
 */type CtelbCreateLoadBalancerApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbCreateLoadBalancerApi(client *core.CtyunClient) *CtelbCreateLoadBalancerApi {
	return &CtelbCreateLoadBalancerApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/create-loadbalancer",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbCreateLoadBalancerApi) Do(ctx context.Context, credential core.Credential, req *CtelbCreateLoadBalancerRequest) (*CtelbCreateLoadBalancerResponse, error) {
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
	var resp CtelbCreateLoadBalancerResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbCreateLoadBalancerRequest struct {
	ClientToken      string `json:"clientToken,omitempty"`      /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID         string `json:"regionID,omitempty"`         /*  区域ID  */
	ProjectID        string `json:"projectID,omitempty"`        /*  企业项目 ID，默认为'0'  */
	VpcID            string `json:"vpcID,omitempty"`            /*  vpc的ID  */
	SubnetID         string `json:"subnetID,omitempty"`         /*  子网的ID  */
	Name             string `json:"name,omitempty"`             /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description      string `json:"description,omitempty"`      /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	EipID            string `json:"eipID,omitempty"`            /*  弹性公网IP的ID。当resourceType=external为必填  */
	SlaName          string `json:"slaName,omitempty"`          /*  lb的规格名称,支持elb.s1.small和elb.default，默认为elb.default  */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型。internal：内网负载均衡，external：公网负载均衡  */
	PrivateIpAddress string `json:"privateIpAddress,omitempty"` /*  负载均衡的私有IP地址，不指定则自动分配  */
	DeleteProtection *bool  `json:"deleteProtection"`           /*  删除保护。false（不开启）、true（开）。 默认：不开启  */
}

type CtelbCreateLoadBalancerResponse struct {
	StatusCode  int32                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbCreateLoadBalancerReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                    `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbCreateLoadBalancerReturnObjResponse struct {
	ID string `json:"ID,omitempty"` /*  负载均衡ID  */
}
