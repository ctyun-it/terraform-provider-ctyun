package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcReplaceSubnetAclApi
/* 子网替换 ACL
 */type CtvpcReplaceSubnetAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcReplaceSubnetAclApi(client *core.CtyunClient) *CtvpcReplaceSubnetAclApi {
	return &CtvpcReplaceSubnetAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/replace-subnet-acl",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcReplaceSubnetAclApi) Do(ctx context.Context, credential core.Credential, req *CtvpcReplaceSubnetAclRequest) (*CtvpcReplaceSubnetAclResponse, error) {
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
	var resp CtvpcReplaceSubnetAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcReplaceSubnetAclRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池 ID  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
	SubnetID    string  `json:"subnetID,omitempty"`    /*  子网 的 ID  */
	AclID       string  `json:"aclID,omitempty"`       /*  ackl 的 ID  */
}

type CtvpcReplaceSubnetAclResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
