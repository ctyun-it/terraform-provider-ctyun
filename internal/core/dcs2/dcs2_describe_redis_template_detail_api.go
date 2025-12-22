package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeRedisTemplateDetailApi
/* 查询参数模板详情。
 */type Dcs2DescribeRedisTemplateDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeRedisTemplateDetailApi(client *core.CtyunClient) *Dcs2DescribeRedisTemplateDetailApi {
	return &Dcs2DescribeRedisTemplateDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisTemplate/describeRedisTemplateDetail",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeRedisTemplateDetailApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeRedisTemplateDetailRequest) (*Dcs2DescribeRedisTemplateDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("templateId", req.TemplateId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeRedisTemplateDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeRedisTemplateDetailRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	TemplateId string `json:"templateId,omitempty"` /*  模板ID。 可调用 <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15104&data=81&isNormal=1&vid=270">查询参数模板列表</a> 接口，使用Template表id字段。  */
}

type Dcs2DescribeRedisTemplateDetailResponse struct {
	StatusCode int32                                             `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                            `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2DescribeRedisTemplateDetailReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                                            `json:"requestId"`  /*  请求 ID。  */
	Code       string                                            `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                            `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2DescribeRedisTemplateDetailReturnObjResponse struct {
	Template *Dcs2DescribeRedisTemplateDetailReturnObjTemplateResponse `json:"template"` /*  总数。  */
	Params   []*Dcs2DescribeRedisTemplateDetailReturnObjParamsResponse `json:"params"`   /*  参数数组。  */
}

type Dcs2DescribeRedisTemplateDetailReturnObjTemplateResponse struct {
	Id          string `json:"id"`          /*  参数记录ID。  */
	Name        string `json:"name"`        /*  参数名称。  */
	Description string `json:"description"` /*  参数描述。  */
	CacheMode   string `json:"cacheMode"`   /*  适合的实例架构版本。<li>ORIGINAL_67：Redis 6.0/7.0类型。<li>ORIGINAL_5：Redis 5.0类型。<li>CLASSIC：经典版。  */
	SysTemplate *bool  `json:"sysTemplate"` /*  是否为系统模板。<li>true：系统模板。<li>false：自定义模板。  */
}

type Dcs2DescribeRedisTemplateDetailReturnObjParamsResponse struct {
	ParamName    string `json:"paramName"`    /*  参数名称。  */
	Description  string `json:"description"`  /*  参数描述。  */
	ValueRange   string `json:"valueRange"`   /*  参数范围。  */
	DefaultValue string `json:"defaultValue"` /*  默认值。  */
	NeedRestart  *bool  `json:"needRestart"`  /*  参数修改后是否需要重启实例。<li>true：需要重启。<li>false：无需重启。  */
	CurrentValue string `json:"currentValue"` /*  当前值。  */
}
