package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2RestartInstanceApi
/* 重启分布式缓存Redis实例
 */type Dcs2RestartInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2RestartInstanceApi(client *core.CtyunClient) *Dcs2RestartInstanceApi {
	return &Dcs2RestartInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceManageMgrServant/restartInstance",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2RestartInstanceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2RestartInstanceRequest) (*Dcs2RestartInstanceResponse, error) {
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
	var resp Dcs2RestartInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2RestartInstanceRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  产品实例标识  */
}

type Dcs2RestartInstanceResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string `json:"code,omitempty"`       /*  响应码描述  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
