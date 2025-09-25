package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2CreateRedisTemplateApi
/* 创建自定义参数模板。
 */type Dcs2CreateRedisTemplateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2CreateRedisTemplateApi(client *core.CtyunClient) *Dcs2CreateRedisTemplateApi {
	return &Dcs2CreateRedisTemplateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisTemplate/createRedisTemplate",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2CreateRedisTemplateApi) Do(ctx context.Context, credential core.Credential, req *Dcs2CreateRedisTemplateRequest) (*Dcs2CreateRedisTemplateResponse, error) {
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
	var resp Dcs2CreateRedisTemplateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CreateRedisTemplateRequest struct {
	RegionId   string                                  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string                                  `json:"prodInstId,omitempty"` /*  实例ID  */
	Params     []*Dcs2CreateRedisTemplateParamsRequest `json:"params"`               /*  参数列表  */
}

type Dcs2CreateRedisTemplateParamsRequest struct {
	ParamName     string `json:"paramName,omitempty"`     /*  参数名称  */
	OriginalValue string `json:"originalValue,omitempty"` /*  当前值  */
	CurrentValue  string `json:"currentValue,omitempty"`  /*  修改的目标值  */
}

type Dcs2CreateRedisTemplateResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2CreateRedisTemplateReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2CreateRedisTemplateReturnObjResponse struct{}
