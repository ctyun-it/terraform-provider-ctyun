package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsEcsGetAutoRenewConfigApi
/* 查询一台包周期付费类型（包年包月）的云主机自动续订配置信息。<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /> */
type CtecsEcsGetAutoRenewConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsEcsGetAutoRenewConfigApi(client *core.CtyunClient) *CtecsEcsGetAutoRenewConfigApi {
	return &CtecsEcsGetAutoRenewConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/get-auto-renew-config",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsEcsGetAutoRenewConfigApi) Do(ctx context.Context, credential core.Credential, req *CtecsEcsGetAutoRenewConfigRequest) (*CtecsEcsGetAutoRenewConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsEcsGetAutoRenewConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtecsEcsGetAutoRenewConfigRequest struct {
	RegionID   string `json:"regionID"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}

type CtecsEcsGetAutoRenewConfigResponse struct {
	StatusCode  int    `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string `json:"errorCode"`   /*  错误码，为product.module.code三段式码  */
	Error       string `json:"error"`       /*  错误码，为product.module.code三段式码  */
	Message     string `json:"message"`     /*  英文描述信息  */
	Description string `json:"description"` /*  中文描述信息  */
	ReturnObj   struct {
		InstanceID          string `json:"instanceID"`
		AutoRenewStatus     int    `json:"autoRenewStatus"`
		AutoRenewCycleCount int    `json:"autoRenewCycleCount"`
		AutoRenewCycleType  string `json:"autoRenewCycleType"`
	} `json:"returnObj"`
}
