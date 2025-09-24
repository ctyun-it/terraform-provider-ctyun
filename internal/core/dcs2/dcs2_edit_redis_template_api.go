package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2EditRedisTemplateApi
/* 修改自定义参数模板
 */type Dcs2EditRedisTemplateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2EditRedisTemplateApi(client *core.CtyunClient) *Dcs2EditRedisTemplateApi {
	return &Dcs2EditRedisTemplateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisTemplate/editRedisTemplate",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2EditRedisTemplateApi) Do(ctx context.Context, credential core.Credential, req *Dcs2EditRedisTemplateRequest) (*Dcs2EditRedisTemplateResponse, error) {
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
	var resp Dcs2EditRedisTemplateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2EditRedisTemplateRequest struct {
	RegionId string                                /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	Template *Dcs2EditRedisTemplateTemplateRequest `json:"template"` /*  模板  */
	Params   []*Dcs2EditRedisTemplateParamsRequest `json:"params"`   /*  参数列表  */
}

type Dcs2EditRedisTemplateTemplateRequest struct {
	Id          string `json:"id,omitempty"`          /*  参数记录ID  */
	Name        string `json:"name,omitempty"`        /*  参数名称  */
	Description string `json:"description,omitempty"` /*  参数描述  */
	CacheMode   string `json:"cacheMode,omitempty"`   /*  适合的实例架构版本<br/><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15104&isNormal=1&vid=270">查询参数模板列表</a> Template表cacheMode字段  */
	SysTemplate *bool  `json:"sysTemplate"`           /*  是否系统模板  */
}

type Dcs2EditRedisTemplateParamsRequest struct {
	ParamName     string `json:"paramName,omitempty"`     /*  参数名称  */
	OriginalValue string `json:"originalValue,omitempty"` /*  当前值  */
	CurrentValue  string `json:"currentValue,omitempty"`  /*  修改的目标值  */
}

type Dcs2EditRedisTemplateResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                  `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2EditRedisTemplateReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                  `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                  `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2EditRedisTemplateReturnObjResponse struct{}
