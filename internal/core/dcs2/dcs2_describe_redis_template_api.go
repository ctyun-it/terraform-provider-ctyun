package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2DescribeRedisTemplateApi
/* 查询可用的参数模版列表。
 */type Dcs2DescribeRedisTemplateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeRedisTemplateApi(client *core.CtyunClient) *Dcs2DescribeRedisTemplateApi {
	return &Dcs2DescribeRedisTemplateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisTemplate/describeRedisTemplate",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeRedisTemplateApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeRedisTemplateRequest) (*Dcs2DescribeRedisTemplateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("type", req.RawType)
	if req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(req.PageNum), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeRedisTemplateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeRedisTemplateRequest struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	RawType  string `json:"type,omitempty"`     /*  模板类型。<li>sys：系统模板。<li>custom：自定义模板。  */
	PageNum  int32  `json:"pageNum,omitempty"`  /*  页码，默认1。  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页记录数，默认值：10。  */
}

type Dcs2DescribeRedisTemplateResponse struct {
	StatusCode int32                                       `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                      `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2DescribeRedisTemplateReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                                      `json:"requestId"`  /*  请求 ID。  */
	Code       string                                      `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                      `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2DescribeRedisTemplateReturnObjResponse struct {
	Total int32                                             `json:"total"` /*  总数。  */
	Size  int32                                             `json:"size"`  /*  本次返回数量。  */
	List  []*Dcs2DescribeRedisTemplateReturnObjListResponse `json:"list"`  /*  参数对象列表。  */
}

type Dcs2DescribeRedisTemplateReturnObjListResponse struct {
	Id          string `json:"id"`          /*  参数记录ID。  */
	Name        string `json:"name"`        /*  参数名称。  */
	Description string `json:"description"` /*  参数描述。  */
	CacheMode   string `json:"cacheMode"`   /*  适合的实例架构版本。<li>ORIGINAL_67：Redis 6.0/7.0类型。<li>ORIGINAL_5：Redis 5.0类型。  */
	SysTemplate *bool  `json:"sysTemplate"` /*  是否为系统模板。<li>true：系统模板。<li>false：自定义模板。  */
}
