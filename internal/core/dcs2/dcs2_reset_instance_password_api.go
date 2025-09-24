package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ResetInstancePasswordApi
/* 重置分布式缓存Redis实例的密码。
 */type Dcs2ResetInstancePasswordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ResetInstancePasswordApi(client *core.CtyunClient) *Dcs2ResetInstancePasswordApi {
	return &Dcs2ResetInstancePasswordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceManageMgrServant/resetInstancePassword",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ResetInstancePasswordApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ResetInstancePasswordRequest) (*Dcs2ResetInstancePasswordResponse, error) {
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
	var resp Dcs2ResetInstancePasswordResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ResetInstancePasswordRequest struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID  */
	NewPassword string `json:"newPassword,omitempty"` /*  新密码<li>长度8-26字符<li>必须同时包含大写字母、小写字母、数字、英文格式特殊符号(@%^*_+!$-=.) 中的三种类型<li>不能有空格  */
}

type Dcs2ResetInstancePasswordResponse struct {
	StatusCode int32                                       `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                      `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ResetInstancePasswordReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                      `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                      `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ResetInstancePasswordReturnObjResponse struct{}
