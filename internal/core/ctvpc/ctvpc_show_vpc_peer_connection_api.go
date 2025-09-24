package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowVpcPeerConnectionApi
/* 查询对等连接详情
 */type CtvpcShowVpcPeerConnectionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowVpcPeerConnectionApi(client *core.CtyunClient) *CtvpcShowVpcPeerConnectionApi {
	return &CtvpcShowVpcPeerConnectionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/get-vpc-peer-connection-attribute",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowVpcPeerConnectionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowVpcPeerConnectionRequest) (*CtvpcShowVpcPeerConnectionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("instanceID", req.InstanceID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowVpcPeerConnectionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowVpcPeerConnectionRequest struct {
	InstanceID string /*  对等连接 ID  */
	RegionID   string /*  区域id  */
}

type CtvpcShowVpcPeerConnectionResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowVpcPeerConnectionReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowVpcPeerConnectionReturnObjResponse struct{}
