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
	_, err := ctReq.WriteJson(struct {
		*Dcs2CreateRedisTemplateRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
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
	RegionId string                                  `json:"regionId,omitempty"` /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	Template *Dcs2CreateRedisTemplateTemplateRequest `json:"template"`           /*  模板。  */
	Params   []*Dcs2CreateRedisTemplateParamsRequest `json:"params"`             /*  参数列表。  */
}

type Dcs2CreateRedisTemplateTemplateRequest struct {
	Name        string `json:"name,omitempty"`        /*  参数名称。  */
	Description string `json:"description,omitempty"` /*  参数描述。  */
	CacheMode   string `json:"cacheMode,omitempty"`   /*  适合的实例架构版本，可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15104&isNormal=1&vid=270">查询参数模板列表</a> 接口，使用Template表cacheMode字段。  */
	SysTemplate *bool  `json:"sysTemplate"`           /*  是否为系统模板。<li>true：系统模板。<li>false：自定义模板。  */
}

type Dcs2CreateRedisTemplateParamsRequest struct {
	ParamName    string `json:"paramName,omitempty"`    /*  参数名称。  */
	CurrentValue string `json:"currentValue,omitempty"` /*  目标值。参数的取值范围可参考<a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15097&isNormal=1&vid=270">查询实例配置参数</a> 返回的param表valueRange字段。  */
}

type Dcs2CreateRedisTemplateResponse struct {
	StatusCode int32                                     `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                    `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2CreateRedisTemplateReturnObjResponse `json:"returnObj"`  /*  无返回数据，空对象。  */
	RequestId  string                                    `json:"requestId"`  /*  请求 ID。  */
	Code       string                                    `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                    `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2CreateRedisTemplateReturnObjResponse struct{}
