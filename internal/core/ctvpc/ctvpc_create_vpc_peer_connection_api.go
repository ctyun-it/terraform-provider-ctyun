package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateVpcPeerConnectionApi
/* 创建对等连接
 */type CtvpcCreateVpcPeerConnectionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateVpcPeerConnectionApi(client *core.CtyunClient) *CtvpcCreateVpcPeerConnectionApi {
	return &CtvpcCreateVpcPeerConnectionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-vpc-peer-connection",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateVpcPeerConnectionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateVpcPeerConnectionRequest) (*CtvpcCreateVpcPeerConnectionResponse, error) {
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
	var resp CtvpcCreateVpcPeerConnectionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateVpcPeerConnectionRequest struct {
	ClientToken    string  `json:"clientToken,omitempty"`    /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RequestVpcID   string  `json:"requestVpcID,omitempty"`   /*  本端vpc id  */
	RequestVpcName string  `json:"requestVpcName,omitempty"` /*  本端vpc的名称  */
	RequestVpcCidr string  `json:"requestVpcCidr,omitempty"` /*  本端vpc的网段  */
	AcceptVpcID    string  `json:"acceptVpcID,omitempty"`    /*  对端的vpc id  */
	AcceptEmail    *string `json:"acceptEmail,omitempty"`    /*  对端vpc账户的邮箱，当建立跨帐号的对等链接是，需要提交  */
	Name           string  `json:"name,omitempty"`           /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	RegionID       string  `json:"regionID,omitempty"`       /*  区域id  */
}

type CtvpcCreateVpcPeerConnectionResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateVpcPeerConnectionReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                        `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateVpcPeerConnectionReturnObjResponse struct {
	Status     *string `json:"status,omitempty"`     /*  创建对等链接状态，取值 in_progress / done  */
	Message    *string `json:"message,omitempty"`    /*  创建状态  */
	InstanceID *string `json:"instanceID,omitempty"` /*  对等链接 ID  */
}
