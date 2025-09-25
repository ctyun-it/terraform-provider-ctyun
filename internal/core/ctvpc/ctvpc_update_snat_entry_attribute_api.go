package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateSnatEntryAttributeApi
/* 修改 SNAT 规则接口
 */type CtvpcUpdateSnatEntryAttributeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateSnatEntryAttributeApi(client *core.CtyunClient) *CtvpcUpdateSnatEntryAttributeApi {
	return &CtvpcUpdateSnatEntryAttributeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/modify-snat-entry",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateSnatEntryAttributeApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateSnatEntryAttributeRequest) (*CtvpcUpdateSnatEntryAttributeResponse, error) {
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
	var resp CtvpcUpdateSnatEntryAttributeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateSnatEntryAttributeRequest struct {
	RegionID       string  `json:"regionID,omitempty"`       /*  区域id  */
	SNatID         string  `json:"sNatID,omitempty"`         /*  SNAT条目所在的SNAT表的ID。  */
	SourceSubnetID *string `json:"sourceSubnetID,omitempty"` /*  子网id，【非自定义情况必传】  */
	SourceCIDR     *string `json:"sourceCIDR,omitempty"`     /*  输入VPC、交换机或ECS实例的网段，还可以输入任意网段。  */
	Description    *string `json:"description,omitempty"`    /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:,'{}  */
	ClientToken    string  `json:"clientToken,omitempty"`    /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcUpdateSnatEntryAttributeResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
