package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcDisassociatePrivateZoneFromVpcApi
/* 内网取消 DNS 关联 VPC
 */type CtvpcDisassociatePrivateZoneFromVpcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDisassociatePrivateZoneFromVpcApi(client *core.CtyunClient) *CtvpcDisassociatePrivateZoneFromVpcApi {
	return &CtvpcDisassociatePrivateZoneFromVpcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone/disassociate-vpc",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDisassociatePrivateZoneFromVpcApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDisassociatePrivateZoneFromVpcRequest) (*CtvpcDisassociatePrivateZoneFromVpcResponse, error) {
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
	var resp CtvpcDisassociatePrivateZoneFromVpcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDisassociatePrivateZoneFromVpcRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池ID  */
	ZoneID      string  `json:"zoneID,omitempty"`      /*  zoneID  */
	VpcIDList   string  `json:"vpcIDList,omitempty"`   /*  关联的vpc,多个vpc用逗号隔开, 最多支持 5 个 VPC  */
}

type CtvpcDisassociatePrivateZoneFromVpcResponse struct {
	StatusCode  int32                                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDisassociatePrivateZoneFromVpcReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcDisassociatePrivateZoneFromVpcReturnObjResponse struct {
	ZoneID *string `json:"zoneID,omitempty"` /*  名称  */
}
