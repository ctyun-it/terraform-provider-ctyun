package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcGetZoneBindVpcsApi
/* 获取 zone 绑定的VPC列表
 */type CtvpcGetZoneBindVpcsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetZoneBindVpcsApi(client *core.CtyunClient) *CtvpcGetZoneBindVpcsApi {
	return &CtvpcGetZoneBindVpcsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/list-vpcs",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetZoneBindVpcsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetZoneBindVpcsRequest) (*CtvpcGetZoneBindVpcsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("zoneID", req.ZoneID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcGetZoneBindVpcsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetZoneBindVpcsRequest struct {
	ZoneID   string /*  zoneID  */
	RegionID string /*  资源池ID  */
}

type CtvpcGetZoneBindVpcsResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGetZoneBindVpcsReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcGetZoneBindVpcsReturnObjResponse struct {
	VpcAssociations []*CtvpcGetZoneBindVpcsReturnObjVpcAssociationsResponse `json:"vpcAssociations"` /*  dns 记录  */
}

type CtvpcGetZoneBindVpcsReturnObjVpcAssociationsResponse struct {
	VpcID   *string `json:"vpcID,omitempty"`   /*  vpc  */
	VpcName *string `json:"vpcName,omitempty"` /*  vpcName  */
}
