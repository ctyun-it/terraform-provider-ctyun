package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaCreateNetworkPathApi
/* 创建网络路径
 */type CtvpcVniaCreateNetworkPathApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaCreateNetworkPathApi(client *core.CtyunClient) *CtvpcVniaCreateNetworkPathApi {
	return &CtvpcVniaCreateNetworkPathApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vnia/create-network-path",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaCreateNetworkPathApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaCreateNetworkPathRequest) (*CtvpcVniaCreateNetworkPathResponse, error) {
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
	var resp CtvpcVniaCreateNetworkPathResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaCreateNetworkPathRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池 ID  */
	Name       string `json:"name,omitempty"`       /*  路径分析名字  */
	SourceID   string `json:"sourceID,omitempty"`   /*  源设备  */
	SourceType string `json:"sourceType,omitempty"` /*  源类型，目前仅支持 ecs / internet / subnet  */
	SourcePort int32  `json:"sourcePort"`           /*  源端口, 1 - 65535  */
	TargetType string `json:"targetType,omitempty"` /*  目标类型，目前仅支持 ecs / internet / subnet / elb  */
	TargetPort int32  `json:"targetPort"`           /*  目标端口, 1 - 65535  */
	TargetID   string `json:"targetID,omitempty"`   /*  目标设备  */
	SourceIP   string `json:"sourceIP,omitempty"`   /*  源 IP  */
	TargetIP   string `json:"targetIP,omitempty"`   /*  目的 IP  */
	Protocol   string `json:"protocol,omitempty"`   /*  协议，仅支持 ICMP / TCP / UDP  */
}

type CtvpcVniaCreateNetworkPathResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
