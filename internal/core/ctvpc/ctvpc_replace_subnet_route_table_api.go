package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcReplaceSubnetRouteTableApi
/* 子网更换路由表，子网必须关联一张路由表。创建VPC后会自动生成一张默认路由表，新建子网时，会关联到默认路由表，子网可以更换其他路由表。
 */type CtvpcReplaceSubnetRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcReplaceSubnetRouteTableApi(client *core.CtyunClient) *CtvpcReplaceSubnetRouteTableApi {
	return &CtvpcReplaceSubnetRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/replace-subnet-route-table",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcReplaceSubnetRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcReplaceSubnetRouteTableRequest) (*CtvpcReplaceSubnetRouteTableResponse, error) {
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
	var resp CtvpcReplaceSubnetRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcReplaceSubnetRouteTableRequest struct {
	ClientToken  *string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID     string  `json:"regionID,omitempty"`     /*  资源池 ID  */
	ProjectID    *string `json:"projectID,omitempty"`    /*  企业项目 ID，默认为0  */
	SubnetID     string  `json:"subnetID,omitempty"`     /*  子网 的 ID  */
	RouteTableID string  `json:"routeTableID,omitempty"` /*  路由表的 ID  */
}

type CtvpcReplaceSubnetRouteTableResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
