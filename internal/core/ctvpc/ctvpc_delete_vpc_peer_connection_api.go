package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDeleteVpcPeerConnectionApi
/* 删除对等连接
 */type CtvpcDeleteVpcPeerConnectionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDeleteVpcPeerConnectionApi(client *core.CtyunClient) *CtvpcDeleteVpcPeerConnectionApi {
	return &CtvpcDeleteVpcPeerConnectionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/delete-vpc-peer-connection",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDeleteVpcPeerConnectionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDeleteVpcPeerConnectionRequest) (*CtvpcDeleteVpcPeerConnectionResponse, error) {
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
	var resp CtvpcDeleteVpcPeerConnectionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDeleteVpcPeerConnectionRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	InstanceID  string `json:"instanceID,omitempty"`  /*  对等连接的唯一值  */
	RegionID    string `json:"regionID,omitempty"`    /*  区域id  */
}

type CtvpcDeleteVpcPeerConnectionResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDeleteVpcPeerConnectionReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                        `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcDeleteVpcPeerConnectionReturnObjResponse struct {
	Status  *string `json:"status,omitempty"`  /*  创建对等链接状态，取值 in_progress / done  */
	Message *string `json:"message,omitempty"` /*  创建状态  */
}
