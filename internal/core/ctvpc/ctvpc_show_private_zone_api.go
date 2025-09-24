package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowPrivateZoneApi
/* 内网DNS详情
 */type CtvpcShowPrivateZoneApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowPrivateZoneApi(client *core.CtyunClient) *CtvpcShowPrivateZoneApi {
	return &CtvpcShowPrivateZoneApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/private-zone/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowPrivateZoneApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowPrivateZoneRequest) (*CtvpcShowPrivateZoneResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("zoneID", req.ZoneID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowPrivateZoneResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowPrivateZoneRequest struct {
	RegionID string /*  资源池ID  */
	ZoneID   string /*  zoneID  */
}

type CtvpcShowPrivateZoneResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowPrivateZoneReturnObjResponse `json:"returnObj"`             /*  object  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowPrivateZoneReturnObjResponse struct {
	ZoneID          *string                                                 `json:"zoneID,omitempty"`       /*  名称  */
	Name            *string                                                 `json:"name,omitempty"`         /*  zone名称  */
	Description     *string                                                 `json:"description,omitempty"`  /*  描述  */
	ProxyPattern    *string                                                 `json:"proxyPattern,omitempty"` /*  zone, record  */
	TTL             int32                                                   `json:"TTL"`                    /*  zone ttl, default is 300  */
	VpcAssociations []*CtvpcShowPrivateZoneReturnObjVpcAssociationsResponse `json:"vpcAssociations"`        /*  vpc关联信息数组  */
	CreatedAt       *string                                                 `json:"createdAt,omitempty"`    /*  创建时间  */
	UpdatedAt       *string                                                 `json:"updatedAt,omitempty"`    /*  更新时间  */
}

type CtvpcShowPrivateZoneReturnObjVpcAssociationsResponse struct {
	VpcAssociation *CtvpcShowPrivateZoneReturnObjVpcAssociationsVpcAssociationResponse `json:"vpcAssociation"` /*  vpc关联信息对象  */
}

type CtvpcShowPrivateZoneReturnObjVpcAssociationsVpcAssociationResponse struct {
	VpcID   *string `json:"vpcID,omitempty"`   /*  vpc  */
	VpcName *string `json:"vpcName,omitempty"` /*  vpcName  */
}
