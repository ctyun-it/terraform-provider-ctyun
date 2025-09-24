package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDisassociateSecondaryCidrsFromVpcApi
/* VPC解绑扩展网段。
 */type CtvpcDisassociateSecondaryCidrsFromVpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDisassociateSecondaryCidrsFromVpcApi(client *core.CtyunClient) *CtvpcDisassociateSecondaryCidrsFromVpcApi {
	return &CtvpcDisassociateSecondaryCidrsFromVpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/disassociate-secondary-cidrs",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDisassociateSecondaryCidrsFromVpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDisassociateSecondaryCidrsFromVpcRequest) (*CtvpcDisassociateSecondaryCidrsFromVpcResponse, error) {
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
	var resp CtvpcDisassociateSecondaryCidrsFromVpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDisassociateSecondaryCidrsFromVpcRequest struct {
	ClientToken string   `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	ProjectID   *string  `json:"projectID,omitempty"`   /*  企业项目 ID，默认为"0"  */
	RegionID    string   `json:"regionID,omitempty"`    /*  资源池ID  */
	VpcID       string   `json:"vpcID,omitempty"`       /*  vpc id  */
	Cidrs       []string `json:"cidrs"`                 /*  是Array类型，里面的内容是String，要解绑的扩展网段 cidr  */
}

type CtvpcDisassociateSecondaryCidrsFromVpcResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
