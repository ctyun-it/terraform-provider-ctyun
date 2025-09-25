package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcSubnetDisassociateRouteTableApi
/* 子网解绑路由表，仅支持 3.0 资源池。
 */type CtvpcSubnetDisassociateRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcSubnetDisassociateRouteTableApi(client *core.CtyunClient) *CtvpcSubnetDisassociateRouteTableApi {
	return &CtvpcSubnetDisassociateRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/subnet-disassociate-route-table",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcSubnetDisassociateRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcSubnetDisassociateRouteTableRequest) (*CtvpcSubnetDisassociateRouteTableResponse, error) {
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
	var resp CtvpcSubnetDisassociateRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcSubnetDisassociateRouteTableRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SubnetID string `json:"subnetID,omitempty"` /*  子网 的 ID  */
}

type CtvpcSubnetDisassociateRouteTableResponse struct {
	StatusCode  int32                                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcSubnetDisassociateRouteTableReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcSubnetDisassociateRouteTableReturnObjResponse struct {
	SubnetID *string `json:"subnetID,omitempty"` /*  vpc 示例 ID  */
}
