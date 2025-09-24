package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbAsyncCreateLoadbalanceApi
/* 创建负载均衡实例，该接口为异步接口，第一次请求会返回资源在创建中，需要用户发起多次请求，直到 status 为 done 为止。
 */type CtelbAsyncCreateLoadbalanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbAsyncCreateLoadbalanceApi(client *core.CtyunClient) *CtelbAsyncCreateLoadbalanceApi {
	return &CtelbAsyncCreateLoadbalanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/async-create-loadbalance",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbAsyncCreateLoadbalanceApi) Do(ctx context.Context, credential core.Credential, req *CtelbAsyncCreateLoadbalanceRequest) (*CtelbAsyncCreateLoadbalanceResponse, error) {
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
	var resp CtelbAsyncCreateLoadbalanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbAsyncCreateLoadbalanceRequest struct {
	ClientToken      string `json:"clientToken,omitempty"`      /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID         string `json:"regionID,omitempty"`         /*  区域ID  */
	VpcID            string `json:"vpcID,omitempty"`            /*  vpc的ID  */
	SubnetID         string `json:"subnetID,omitempty"`         /*  子网的ID  */
	Name             string `json:"name,omitempty"`             /*  唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	EipID            string `json:"eipID,omitempty"`            /*  弹性公网IP的ID。当resourceType=external为必填  */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型。internal：内网负载均衡，external：公网负载均衡  */
	PrivateIpAddress string `json:"privateIpAddress,omitempty"` /*  负载均衡的私有IP地址，不指定则自动分配  */
}

type CtelbAsyncCreateLoadbalanceResponse struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbAsyncCreateLoadbalanceReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                        `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbAsyncCreateLoadbalanceReturnObjResponse struct {
	Status        string `json:"status,omitempty"`        /*  创建进度: in_progress / done  */
	Message       string `json:"message,omitempty"`       /*  进度说明  */
	LoadBalanceID string `json:"loadBalanceID,omitempty"` /*  负载均衡ID，可能为 null  */
}
