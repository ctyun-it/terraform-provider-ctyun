package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeRedisTemplateDetailApi
/* 查询在不同版本的参数模版中支持设置的参数列表。
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
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	TemplateId string /*  模板ID  */
}

type Dcs2DescribeRedisTemplateDetailResponse struct {
	StatusCode int32                                             `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                            `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeRedisTemplateDetailReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                            `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                            `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                            `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeRedisTemplateDetailReturnObjResponse struct {
	Template *Dcs2DescribeRedisTemplateDetailReturnObjTemplateResponse `json:"template"` /*  总数  */
	Params   []*Dcs2DescribeRedisTemplateDetailReturnObjParamsResponse `json:"params"`   /*  参数数组  */
}

type Dcs2DescribeRedisTemplateDetailReturnObjTemplateResponse struct {
	Id          string `json:"id,omitempty"`          /*  参数记录ID  */
	Name        string `json:"name,omitempty"`        /*  参数名称  */
	Description string `json:"description,omitempty"` /*  参数描述  */
	CacheMode   string `json:"cacheMode,omitempty"`   /*  适合的实例架构版本  */
	SysTemplate *bool  `json:"sysTemplate"`           /*  是否为系统模板<li>true：系统模板<li>false：自定义模板  */
}

type Dcs2DescribeRedisTemplateDetailReturnObjParamsResponse struct {
	ParamName    string `json:"paramName,omitempty"`    /*  参数名称  */
	Description  string `json:"description,omitempty"`  /*  参数描述  */
	ValueRange   string `json:"valueRange,omitempty"`   /*  参数范围  */
	DefaultValue string `json:"defaultValue,omitempty"` /*  默认值  */
	NeedRestart  *bool  `json:"needRestart"`            /*  参数修改后是否需要重启实例<li>true：需要重启<li>false：无需重启  */
	CurrentValue string `json:"currentValue,omitempty"` /*  当前值  */
}
