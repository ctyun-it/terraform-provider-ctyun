package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeSecurityIpsApi
/* 查询分布式缓存Redis实例的IP白名单分组。
 */type Dcs2DescribeSecurityIpsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeSecurityIpsApi(client *core.CtyunClient) *Dcs2DescribeSecurityIpsApi {
	return &Dcs2DescribeSecurityIpsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/securityMgrServant/describeSecurityIps",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeSecurityIpsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeSecurityIpsRequest) (*Dcs2DescribeSecurityIpsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeSecurityIpsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeSecurityIpsRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
}

type Dcs2DescribeSecurityIpsResponse struct {
	StatusCode int32                                     `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                    `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2DescribeSecurityIpsReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                                    `json:"requestId"`  /*  请求 ID。  */
	Code       string                                    `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                    `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2DescribeSecurityIpsReturnObjResponse struct {
	Total int32                                           `json:"total"` /*  数量。  */
	Rows  []*Dcs2DescribeSecurityIpsReturnObjRowsResponse `json:"rows"`  /*  白名单分组列表。  */
}

type Dcs2DescribeSecurityIpsReturnObjRowsResponse struct {
	Group string `json:"group"` /*  分组。  */
	Ip    string `json:"ip"`    /*  白名单集合。  */
}
