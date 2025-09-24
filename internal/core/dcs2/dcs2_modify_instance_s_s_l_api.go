package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyInstanceSSLApi
/* 打开或关闭SSL加密,或更新SSL证书有效期。注意：开启或关闭SSL开关、更新SSL证书有效期将会重启实例导致业务短暂不可用，需要等待一定时间方可生效，并且需要应用修改SSL的连接方式，请慎重操作!
 */type Dcs2ModifyInstanceSSLApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyInstanceSSLApi(client *core.CtyunClient) *Dcs2ModifyInstanceSSLApi {
	return &Dcs2ModifyInstanceSSLApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/securityMgrServant/modifyInstanceSSL",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyInstanceSSLApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyInstanceSSLRequest) (*Dcs2ModifyInstanceSSLResponse, error) {
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
	var resp Dcs2ModifyInstanceSSLResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyInstanceSSLRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	SslEnabled string `json:"sslEnabled,omitempty"` /*  操作类型<li>Disable: 关闭SSL<li>Enable: 开启SSL）<li>Update: 更新证书有效期  */
}

type Dcs2ModifyInstanceSSLResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                  `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ModifyInstanceSSLReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                  `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                  `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifyInstanceSSLReturnObjResponse struct{}
