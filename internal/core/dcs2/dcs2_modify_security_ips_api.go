package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifySecurityIpsApi
/* 分布式缓存Redis实例添加白名单分组，重写白名单分组，删除白名单分组。
 */type Dcs2ModifySecurityIpsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifySecurityIpsApi(client *core.CtyunClient) *Dcs2ModifySecurityIpsApi {
	return &Dcs2ModifySecurityIpsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/securityMgrServant/modifySecurityIps",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifySecurityIpsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifySecurityIpsRequest) (*Dcs2ModifySecurityIpsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ModifySecurityIpsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifySecurityIpsRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Group      string `json:"group,omitempty"`      /*  白名单分组名  */
	Mode       string `json:"mode,omitempty"`       /*  修改方式<li>cover：覆盖原白名单分组<li>append：新增白名单分组<li>delete：删除该白名单分组  */
	Ip         string `json:"ip,omitempty"`         /*  白名单列表<li>您可填写IP地址(如192.168.1.1)或IP段(如192.168.1.0/24)<li>同时添加多个IP请使用英文逗号隔开如192.168.0.1,192.168.1.0/24<br>说明：当 `mode=delete` 时，此参数为空  */
}

type Dcs2ModifySecurityIpsResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                  `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ModifySecurityIpsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                  `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                  `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifySecurityIpsReturnObjResponse struct{}
