package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyInstanceMajorVersionApi
/* 升级引擎大版本
 */type Dcs2ModifyInstanceMajorVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyInstanceMajorVersionApi(client *core.CtyunClient) *Dcs2ModifyInstanceMajorVersionApi {
	return &Dcs2ModifyInstanceMajorVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceManageMgrServant/modifyInstanceMajorVersion",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyInstanceMajorVersionApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyInstanceMajorVersionRequest) (*Dcs2ModifyInstanceMajorVersionResponse, error) {
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
	var resp Dcs2ModifyInstanceMajorVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyInstanceMajorVersionRequest struct {
	RegionId      string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId    string `json:"prodInstId,omitempty"`    /*  实例ID  */
	EngineVersion string `json:"engineVersion,omitempty"` /*  引擎大版本<br/><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=17922&isNormal=1&vid=270">查询实例版本信息</a> engineMajorVersionInfo表upgradableEngineVersionItems中的版本  */
}

type Dcs2ModifyInstanceMajorVersionResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string `json:"code,omitempty"`       /*  响应码描述  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
