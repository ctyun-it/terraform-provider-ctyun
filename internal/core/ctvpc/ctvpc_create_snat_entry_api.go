package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateSnatEntryApi
/* 创建 SNAT 规则
 */type CtvpcCreateSnatEntryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateSnatEntryApi(client *core.CtyunClient) *CtvpcCreateSnatEntryApi {
	return &CtvpcCreateSnatEntryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/create-snat-entry",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateSnatEntryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateSnatEntryRequest) (*CtvpcCreateSnatEntryResponse, error) {
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
	var resp CtvpcCreateSnatEntryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateSnatEntryRequest struct {
	RegionID       string   `json:"regionID,omitempty"`       /*  区域id  */
	NatGatewayID   string   `json:"natGatewayID,omitempty"`   /*  NAT网关ID  */
	SourceSubnetID *string  `json:"sourceSubnetID,omitempty"` /*  子网id，【非自定义情况必传 sourceCIDR和sourceSubnetID二选一必传】  */
	SourceCIDR     *string  `json:"sourceCIDR,omitempty"`     /*  自定义输入VPC、交换机或ECS实例的网段，还可以输入任意网段。【自定义子网信息必传】】  */
	SnatIps        []string `json:"snatIps"`                  /*  [eip-p8cdvc4srg]  */
	ClientToken    string   `json:"clientToken,omitempty"`    /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcCreateSnatEntryResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateSnatEntryReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreateSnatEntryReturnObjResponse struct {
	Snat *CtvpcCreateSnatEntryReturnObjSnatResponse `json:"snat"` /*  业务数据  */
}

type CtvpcCreateSnatEntryReturnObjSnatResponse struct {
	Status  *string `json:"status,omitempty"`  /*  绑定状态，取值 in_progress / done  */
	Message *string `json:"message,omitempty"` /*  绑定状态提示信息  */
	SnatID  *string `json:"snatID,omitempty"`  /*  snat id, 当 status != done 时，snatID 为 null  */
}
