package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcModifyEndpointServiceIpVersionApi
/* 转换终端节点服务IP版本
 */type CtvpcModifyEndpointServiceIpVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcModifyEndpointServiceIpVersionApi(client *core.CtyunClient) *CtvpcModifyEndpointServiceIpVersionApi {
	return &CtvpcModifyEndpointServiceIpVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/modify-endpoint-service-ip-version",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcModifyEndpointServiceIpVersionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcModifyEndpointServiceIpVersionRequest) (*CtvpcModifyEndpointServiceIpVersionResponse, error) {
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
	var resp CtvpcModifyEndpointServiceIpVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcModifyEndpointServiceIpVersionRequest struct {
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池ID  */
	EndpointServiceID string  `json:"endpointServiceID,omitempty"` /*  终端节点服务id  */
	IpVersion         int32   `json:"ipVersion"`                   /*  0:ipv4, 1:ipv6（暂不支持）, 2:双栈  */
	SubnetID          *string `json:"subnetID,omitempty"`          /*  子网ID，overlay反向型终端节点服务必填  */
	Ipv6Address       *string `json:"ipv6Address,omitempty"`       /*  反向型终端节点服务ipv6中转ip  */
	InstanceIDV6      *string `json:"instanceIDV6,omitempty"`      /*  后端ipv6实例id，服务后端为vip时必填  */
	UnderlayIp6       *string `json:"underlayIp6,omitempty"`       /*  ipv6 underlay ip，服务为天翼云服务时必填  */
}

type CtvpcModifyEndpointServiceIpVersionResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
