package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcResourceBindLabelApi
/* 资源绑定标签
 */type CtvpcResourceBindLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcResourceBindLabelApi(client *core.CtyunClient) *CtvpcResourceBindLabelApi {
	return &CtvpcResourceBindLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/labels/resource_bind_label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcResourceBindLabelApi) Do(ctx context.Context, credential core.Credential, req *CtvpcResourceBindLabelRequest) (*CtvpcResourceBindLabelResponse, error) {
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
	var resp CtvpcResourceBindLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcResourceBindLabelRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域ID  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型，resourceType only support vpc / subnet / acl / security_group / route_table / havip / port  / multicast_domain / vpc_peer / vpce_endpoint / vpce_endpoint_service / ipv6_gateway / elb /                private_nat / nat / eip / bandwidth /ipv6_bandwidth  */
	ResourceID   string `json:"resourceID,omitempty"`   /*  资源 ID  */
	LabelKey     string `json:"labelKey,omitempty"`     /*  标签 key  */
	LabelValue   string `json:"labelValue,omitempty"`   /*  标签 取值  */
}

type CtvpcResourceBindLabelResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
