package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteVpcApi
/* 删除专有网络删除专有网络之前，需要先删除所有子网，且需要删除子网内所有的云资源，包括ECS、弹性裸金属服务器、弹性负载均衡、NAT网关、高可用虚拟 IP 等，需要将子网内的占用IP的资源全部释放。
 */type CtvpcDeleteVpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteVpcApi(client *core.CtyunClient) *CtvpcDeleteVpcApi {
	return &CtvpcDeleteVpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteVpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteVpcRequest) (*CtvpcDeleteVpcResponse, error) {
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
	var resp CtvpcDeleteVpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteVpcRequest struct {
	RegionID  string  `json:"regionID,omitempty"`  /*  资源池 ID  */
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目 ID，默认为0  */
	VpcID     string  `json:"vpcID,omitempty"`     /*  VPC 的 ID  */
}

type CtvpcDeleteVpcResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
