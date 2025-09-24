package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeTopHotKeysPolicyApi
/* 查询分布式缓存Redis实例热key自动分析策略。
 */type Dcs2DescribeTopHotKeysPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeTopHotKeysPolicyApi(client *core.CtyunClient) *Dcs2DescribeTopHotKeysPolicyApi {
	return &Dcs2DescribeTopHotKeysPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/keyAnalysisMgrServant/describeTopHotKeysPolicy",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeTopHotKeysPolicyApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeTopHotKeysPolicyRequest) (*Dcs2DescribeTopHotKeysPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeTopHotKeysPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeTopHotKeysPolicyRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeTopHotKeysPolicyResponse struct {
	StatusCode int32                                          `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	ReturnObj  *Dcs2DescribeTopHotKeysPolicyReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	Message    string                                         `json:"message,omitempty"`    /*  响应信息  */
	RequestId  string                                         `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                         `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeTopHotKeysPolicyReturnObjResponse struct {
	Enable_auto_scanHotkey string `json:"enable_auto_scanHotkey,omitempty"` /*  是否开启热key备份策略<li>true: 开启<li>false：关闭  */
	Schedule_days          string `json:"schedule_days,omitempty"`          /*  星期几触发  */
	Schedule_hours         string `json:"schedule_hours,omitempty"`         /*  整点时刻触发  */
}
