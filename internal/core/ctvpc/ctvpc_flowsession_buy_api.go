package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcFlowsessionBuyApi
/* 购买流量会话
 */type CtvpcFlowsessionBuyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcFlowsessionBuyApi(client *core.CtyunClient) *CtvpcFlowsessionBuyApi {
	return &CtvpcFlowsessionBuyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/flowsession/buy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcFlowsessionBuyApi) Do(ctx context.Context, credential core.Credential, req *CtvpcFlowsessionBuyRequest) (*CtvpcFlowsessionBuyResponse, error) {
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
	var resp CtvpcFlowsessionBuyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcFlowsessionBuyRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  区域ID  */
	ClientToken    string `json:"clientToken,omitempty"`    /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	MirrorFilterID string `json:"mirrorFilterID,omitempty"` /*  过滤条件 ID  */
	SrcPort        string `json:"srcPort,omitempty"`        /*  源弹性网卡 ID  */
	DstPort        string `json:"dstPort,omitempty"`        /*  目的弹性网卡 ID  */
	SubnetID       string `json:"subnetID,omitempty"`       /*  子网 ID  */
	Vni            int32  `json:"vni"`                      /*  VXLAN 网络标识符, 0 - 1677215  */
	Name           string `json:"name,omitempty"`           /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	CycleType      string `json:"cycleType,omitempty"`      /*  仅支持按需  */
}

type CtvpcFlowsessionBuyResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcFlowsessionBuyReturnObjResponse `json:"returnObj"`             /*  object  */
}

type CtvpcFlowsessionBuyReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id。  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  可以为 null。  */
	RegionID             *string `json:"regionID,omitempty"`             /*  可用区id。  */
	NatGatewayID         *string `json:"natGatewayID,omitempty"`         /*  flow_session ID  */
}
