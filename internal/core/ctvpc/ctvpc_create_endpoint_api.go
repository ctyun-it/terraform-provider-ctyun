package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateEndpointApi
/* 创建终端节点服务
 */type CtvpcCreateEndpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateEndpointApi(client *core.CtyunClient) *CtvpcCreateEndpointApi {
	return &CtvpcCreateEndpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/create-endpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateEndpointApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateEndpointRequest) (*CtvpcCreateEndpointResponse, error) {
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
	var resp CtvpcCreateEndpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateEndpointRequest struct {
	ClientToken       string    `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string    `json:"regionID,omitempty"`          /*  资源池ID  */
	CycleType         string    `json:"cycleType,omitempty"`         /*  收费类型：只能填写 on_demand  */
	EndpointServiceID string    `json:"endpointServiceID,omitempty"` /*  终端节点关联的终端节点服务  */
	IpVersion         int32     `json:"ipVersion"`                   /*  0:ipv4, 1:ipv6（暂不支持）, 2:双栈，默认0  */
	EndpointName      string    `json:"endpointName,omitempty"`      /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	SubnetID          string    `json:"subnetID,omitempty"`          /*  子网id  */
	VpcID             string    `json:"vpcID,omitempty"`             /*  vpc-xxxx  */
	IP                *string   `json:"IP,omitempty"`                /*  ipv4 vpc address  */
	IP6               *string   `json:"IP6,omitempty"`               /*  ipv6 vpc address  */
	WhitelistFlag     int32     `json:"whitelistFlag"`               /*  白名单开关 1.开启 0.关闭，默认1  */
	Whitelist         []*string `json:"whitelist"`                   /*  白名单  */
	Whitelist6        []*string `json:"whitelist6"`                  /*  ipv6白名单  */
	EnableDns         *bool     `json:"enableDns"`                   /*  是否开启dns, true:开启,false:关闭  */
	PayVoucherPrice   *string   `json:"payVoucherPrice,omitempty"`   /*  代金券金额，支持到小数点后两位  */
	DeleteProtection  *bool     `json:"deleteProtection"`            /*  是否开启删除保护, true:开启,false:关闭，不传默认关闭  */
	ProtectionService *string   `json:"protectionService,omitempty"` /*  删除保护使能服务  */
}

type CtvpcCreateEndpointResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateEndpointReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcCreateEndpointReturnObjResponse struct {
	MasterOrderID        *string `json:"masterOrderID,omitempty"`        /*  订单id  */
	MasterOrderNO        *string `json:"masterOrderNO,omitempty"`        /*  订单编号  */
	RegionID             *string `json:"regionID,omitempty"`             /*  资源池ID  */
	MasterResourceStatus *string `json:"masterResourceStatus,omitempty"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     *string `json:"masterResourceID,omitempty"`     /*  资源id  */
	EndpointID           *string `json:"endpointID,omitempty"`           /*  终端节点ID, 当 masterResourceStatus 不为 started 时，该取值可为 null  */
}
