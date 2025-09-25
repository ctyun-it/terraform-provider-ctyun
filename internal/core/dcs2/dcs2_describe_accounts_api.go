package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeAccountsApi
/* 查看分布式缓存Redis实例的帐户列表信息。
 */type Dcs2DescribeAccountsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeAccountsApi(client *core.CtyunClient) *Dcs2DescribeAccountsApi {
	return &Dcs2DescribeAccountsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/userMgr/describeAccounts",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeAccountsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeAccountsRequest) (*Dcs2DescribeAccountsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeAccountsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeAccountsRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeAccountsResponse struct {
	StatusCode int32                                  `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                 `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeAccountsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                 `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                 `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeAccountsReturnObjResponse struct {
	Total int32                                        `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeAccountsReturnObjRowsResponse `json:"rows"`            /*  账号集合，见Account  */
}

type Dcs2DescribeAccountsReturnObjRowsResponse struct {
	Name               string `json:"name,omitempty"`               /*  账户名称  */
	RawType            string `json:"type,omitempty"`               /*  账户类型<li>ro:只读<li>rw：读写<li>dba:DBA所有权限<li>wo:只写<li>sync:同步权限  */
	AccountDescription string `json:"accountDescription,omitempty"` /*  账户描述信息  */
}
