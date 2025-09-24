package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcPortReplaceSubnetApi
/* 修改内网IP
 */type CtvpcPortReplaceSubnetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPortReplaceSubnetApi(client *core.CtyunClient) *CtvpcPortReplaceSubnetApi {
	return &CtvpcPortReplaceSubnetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/change-private-ip",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPortReplaceSubnetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPortReplaceSubnetRequest) (*CtvpcPortReplaceSubnetResponse, error) {
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
	var resp CtvpcPortReplaceSubnetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPortReplaceSubnetRequest struct {
	RegionID           string  `json:"regionID,omitempty"`           /*  资源池ID  */
	NetworkInterfaceID string  `json:"networkInterfaceID,omitempty"` /*  网卡ID  */
	IpAddress          *string `json:"ipAddress,omitempty"`          /*  内网 IP 地址，如果不传为自动分配  */
	SubnetID           string  `json:"subnetID,omitempty"`           /*  子网 ID  */
	InstanceID         string  `json:"instanceID,omitempty"`         /*  弹性云主机 ID  */
}

type CtvpcPortReplaceSubnetResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
