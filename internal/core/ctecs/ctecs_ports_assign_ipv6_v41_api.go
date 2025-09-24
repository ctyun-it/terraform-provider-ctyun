package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsPortsAssignIpv6V41Api
/* 网卡关联多个IPv6地址<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsPortsAssignIpv6V41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsPortsAssignIpv6V41Api(client *core.CtyunClient) *CtecsPortsAssignIpv6V41Api {
	return &CtecsPortsAssignIpv6V41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/ports/assign-ipv6",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsPortsAssignIpv6V41Api) Do(ctx context.Context, credential core.Credential, req *CtecsPortsAssignIpv6V41Request) (*CtecsPortsAssignIpv6V41Response, error) {
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
	var resp CtecsPortsAssignIpv6V41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsPortsAssignIpv6V41Request struct {
	ClientToken        string   `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个clientToken值，则代表为同一个请求。保留时间为24小时  */
	RegionID           string   `json:"regionID,omitempty"`           /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	NetworkInterfaceID string   `json:"networkInterfaceID,omitempty"` /*  网卡ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10197673">弹性网卡-弹性网卡基本知识</a>来了解弹性网卡<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5802&data=94">查询弹性网卡列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5789&data=94">创建弹性网卡</a>  */
	Ipv6AddressesCount int32    `json:"ipv6AddressesCount,omitempty"` /*  Ipv6地址数量，新增自动分配地址的IPv6的数量，与ipv6Addresses二选一  */
	Ipv6Addresses      []string `json:"ipv6Addresses"`                /*  IPv6地址列表，新增指定地址的IPv6列表，与ipv6AddressesCount二选一  */
}

type CtecsPortsAssignIpv6V41Response struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string `json:"message,omitempty"`     /*  英文描述信息  */
	Description string `json:"description,omitempty"` /*  中文描述信息  */
}
