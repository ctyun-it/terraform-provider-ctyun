package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateSecurityGroupIngressV41Api
/* 该接口提供用户创建安全组入向规则功能<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br /><b>注意事项：</b><br />&emsp;&emsp;入向规则：direction参数填写”ingress“
 */type CtecsCreateSecurityGroupIngressV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateSecurityGroupIngressV41Api(client *core.CtyunClient) *CtecsCreateSecurityGroupIngressV41Api {
	return &CtecsCreateSecurityGroupIngressV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/vpc/create-security-group-ingress",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateSecurityGroupIngressV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreateSecurityGroupIngressV41Request) (*CtecsCreateSecurityGroupIngressV41Response, error) {
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
	var resp CtecsCreateSecurityGroupIngressV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateSecurityGroupIngressV41Request struct {
	RegionID           string                                                         `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	SecurityGroupID    string                                                         `json:"securityGroupID,omitempty"` /*  安全组ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	SecurityGroupRules []*CtecsCreateSecurityGroupIngressV41SecurityGroupRulesRequest `json:"securityGroupRules"`        /*  规则信息列表，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关配置  */
	ClientToken        string                                                         `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性。长度为1-64字符，要求单个云平台账户内唯一，使用同一个clientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时  */
}

type CtecsCreateSecurityGroupIngressV41SecurityGroupRulesRequest struct {
	Direction             string `json:"direction,omitempty"`             /*  规则方向，入方向则填写ingress  */
	RemoteType            int32  `json:"remoteType,omitempty"`            /*  remote 类型，0 表示使用 cidr，1 表示使用远端安全组，默认为 0  */
	RemoteSecurityGroupID string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组 id  */
	Action                string `json:"action,omitempty"`                /*  授权策略，取值范围：accept（允许），drop（拒绝）。  */
	Priority              int32  `json:"priority,omitempty"`              /*  优先级，取值范围：[1, 100] <br />注：取值越小优先级越大。  */
	Protocol              string `json:"protocol,omitempty"`              /*  网络协议，取值范围：ANY（任意）、TCP、UDP、ICMP(v4)  */
	Ethertype             string `json:"ethertype,omitempty"`             /*  IP类型，取值范围：IPv4、IPv6  */
	DestCidrIp            string `json:"destCidrIp,omitempty"`            /*  远端地址:0.0.0.0/0  */
	Description           string `json:"description,omitempty"`           /*  安全组规则描述信息，满足以下规则：<br />① 长度0-128字符，<br />② 支持拉丁字母、中文、数字, 特殊字符<br />！@#￥%……&*（）——-+={}《》？：“”【】、；‘'，。、<br />不能以 http: / https: 开头  */
	RawRange              string `json:"range,omitempty"`                 /*  安全组开放的传输层协议相关的源端端口范围  */
}

type CtecsCreateSecurityGroupIngressV41Response struct {
	StatusCode  int32                                                `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                               `json:"description,omitempty"` /*  中文描述信息  */
	ErrorCode   string                                               `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                               `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	ReturnObj   *CtecsCreateSecurityGroupIngressV41ReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtecsCreateSecurityGroupIngressV41ReturnObjResponse struct {
	SgRuleIDs []string `json:"sgRuleIDs"` /*  安全组规则 id 列表  */
}
