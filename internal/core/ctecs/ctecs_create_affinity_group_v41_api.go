package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateAffinityGroupV41Api
/* 该接口提供用户创建云主机组的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;用户配额：确认个人在不同资源池下资源配额，可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9714&data=87">用户配额查询</a>接口进行查询<br />
 */type CtecsCreateAffinityGroupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateAffinityGroupV41Api(client *core.CtyunClient) *CtecsCreateAffinityGroupV41Api {
	return &CtecsCreateAffinityGroupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/affinity-group/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateAffinityGroupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreateAffinityGroupV41Request) (*CtecsCreateAffinityGroupV41Response, error) {
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
	var resp CtecsCreateAffinityGroupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateAffinityGroupV41Request struct {
	RegionID          string `json:"regionID,omitempty"`          /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AffinityGroupName string `json:"affinityGroupName,omitempty"` /*  云主机组名称，满足以下规则：长度在1-64个字符，只能由中文、英文字母、数字、下划线_、中划线-、点.组成  */
	PolicyType        int32  `json:"policyType"`                  /*  云主机组策略类型。<br />取值范围：<br />0：强制反亲和性，<br />1：强制亲和性，<br />2：反亲和性，<br />3：亲和性，<br />4：电力反亲和性<br />注：默认值2  */
}

type CtecsCreateAffinityGroupV41Response struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsCreateAffinityGroupV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsCreateAffinityGroupV41ReturnObjResponse struct {
	AffinityGroupID     string                                                           `json:"affinityGroupID,omitempty"`   /*  云主机组ID  */
	AffinityGroupName   string                                                           `json:"affinityGroupName,omitempty"` /*  云主机组名称  */
	AffinityGroupPolicy *CtecsCreateAffinityGroupV41ReturnObjAffinityGroupPolicyResponse `json:"affinityGroupPolicy"`         /*  云主机组策略  */
}

type CtecsCreateAffinityGroupV41ReturnObjAffinityGroupPolicyResponse struct {
	PolicyType     int32  `json:"policyType,omitempty"`     /*  云主机组策略类型。<br />取值范围：<br />0：强制反亲和性，<br />1：强制亲和性，<br />2：反亲和性，<br />3：亲和性，<br />4：电力反亲和性  */
	PolicyTypeName string `json:"policyTypeName,omitempty"` /*  云主机组策略类型名称<br />取值范围：<br />anti-affinity：强制反亲和性，<br />affinity：强制亲和性，<br />soft-anti-affinity：反亲和性，<br />soft-affinity：亲和性，<br />power-anti-affinity：电力反亲和性  */
}
