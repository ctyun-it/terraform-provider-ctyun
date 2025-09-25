package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsModifySecurityGroupIngressV41Api
/* 该接口提供用户修改安全组入方向规则描述的功能<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsModifySecurityGroupIngressV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsModifySecurityGroupIngressV41Api(client *core.CtyunClient) *CtecsModifySecurityGroupIngressV41Api {
	return &CtecsModifySecurityGroupIngressV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/vpc/modify-security-group-ingress",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsModifySecurityGroupIngressV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsModifySecurityGroupIngressV41Request) (*CtecsModifySecurityGroupIngressV41Response, error) {
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
	var resp CtecsModifySecurityGroupIngressV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsModifySecurityGroupIngressV41Request struct {
	RegionID              string `json:"regionID,omitempty"`              /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	SecurityGroupID       string `json:"securityGroupID,omitempty"`       /*  安全组ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	SecurityGroupRuleID   string `json:"securityGroupRuleID,omitempty"`   /*  安全组入向规则ID，获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4818&data=94">查询用户安全组详情</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4823&data=94">创建安全组入向规则</a><br />注：在多可用区类型资源池下，安全组ID通常以“sg-”开头，非多可用区类型资源池安全组ID为uuid格式  */
	Description           string `json:"description,omitempty"`           /*  安全组规则描述信息，满足以下规则：<br />① 长度0-128字符，<br />② 支持拉丁字母、中文、数字, 特殊字符<br />！@#￥%……&*（）——-+={}《》？：“”【】、；‘'，。、<br />不能以 http: / https: 开头  */
	ClientToken           string `json:"clientToken,omitempty"`           /*  客户端存根，用于保证订单幂等性。长度为1-64字符，要求单个云平台账户内唯一，使用同一个clientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时  */
	Action                string `json:"action,omitempty"`                /*  拒绝策略:允许-accept 拒绝-drop  */
	Priority              int32  `json:"priority,omitempty"`              /*  优先级:1~100，取值越小优先级越大  */
	Protocol              string `json:"protocol,omitempty"`              /*  协议: ANY、TCP、UDP、ICMP(v4)  */
	RemoteSecurityGroupID string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组id  */
	DestCidrIp            string `json:"destCidrIp,omitempty"`            /*  cidr  */
	RemoteType            int32  `json:"remoteType,omitempty"`            /*  远端类型，0 表示 destCidrIp，1 表示 remoteSecurityGroupID, 2 表示 prefixlistID，默认为 0  */
	PrefixListID          string `json:"prefixListID,omitempty"`          /*  前缀列表  */
}

type CtecsModifySecurityGroupIngressV41Response struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string `json:"message,omitempty"`     /*  英文描述信息  */
	Description string `json:"description,omitempty"` /*  中文描述信息  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
}
