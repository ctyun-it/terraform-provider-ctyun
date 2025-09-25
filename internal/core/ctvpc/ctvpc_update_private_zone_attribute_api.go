package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdatePrivateZoneAttributeApi
/* 修改内网DNS
 */type CtvpcUpdatePrivateZoneAttributeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdatePrivateZoneAttributeApi(client *core.CtyunClient) *CtvpcUpdatePrivateZoneAttributeApi {
	return &CtvpcUpdatePrivateZoneAttributeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/private-zone/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdatePrivateZoneAttributeApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdatePrivateZoneAttributeRequest) (*CtvpcUpdatePrivateZoneAttributeResponse, error) {
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
	var resp CtvpcUpdatePrivateZoneAttributeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdatePrivateZoneAttributeRequest struct {
	ClientToken  *string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID     string  `json:"regionID,omitempty"`     /*  资源池ID  */
	VpcIDList    *string `json:"vpcIDList,omitempty"`    /*  关联的vpc,多个ID之间用半角逗号（,）隔开, 最多同时支持 5 个 VPC  */
	ZoneID       string  `json:"zoneID,omitempty"`       /*  zoneID  */
	ProxyPattern *string `json:"proxyPattern,omitempty"` /*  zone：当前可用区不进行递归解析。 record：不完全劫持，进行递归解析代理  */
	TTL          int32   `json:"TTL"`                    /*  zone ttl, 单位秒。default is 300，大于300，小于2147483647  */
}

type CtvpcUpdatePrivateZoneAttributeResponse struct {
	StatusCode  int32                                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcUpdatePrivateZoneAttributeReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcUpdatePrivateZoneAttributeReturnObjResponse struct {
	ZoneID *string `json:"zoneID,omitempty"` /*  名称  */
}
