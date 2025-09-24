package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAssociateSecondaryCidrsToVpcApi
/* VPC 绑定扩展网段
 */type CtvpcAssociateSecondaryCidrsToVpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAssociateSecondaryCidrsToVpcApi(client *core.CtyunClient) *CtvpcAssociateSecondaryCidrsToVpcApi {
	return &CtvpcAssociateSecondaryCidrsToVpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/associate-secondary-cidrs",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAssociateSecondaryCidrsToVpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAssociateSecondaryCidrsToVpcRequest) (*CtvpcAssociateSecondaryCidrsToVpcResponse, error) {
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
	var resp CtvpcAssociateSecondaryCidrsToVpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAssociateSecondaryCidrsToVpcRequest struct {
	ClientToken string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性， 长度 1 - 64  */
	ProjectID   *string  `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
	RegionID    string   `json:"regionID,omitempty"`    /*  资源池ID  */
	VpcID       string   `json:"vpcID,omitempty"`       /*  vpc id  */
	Cidrs       []string `json:"cidrs"`                 /*  是Array类型，里面的内容是String，要绑定的扩展网段ip  */
}

type CtvpcAssociateSecondaryCidrsToVpcResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
