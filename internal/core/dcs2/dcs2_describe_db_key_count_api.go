package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeDbKeyCountApi
/* 查询DB编号中缓存key的数量。
 */type Dcs2DescribeDbKeyCountApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeDbKeyCountApi(client *core.CtyunClient) *Dcs2DescribeDbKeyCountApi {
	return &Dcs2DescribeDbKeyCountApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeDbKeyCount",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeDbKeyCountApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeDbKeyCountRequest) (*Dcs2DescribeDbKeyCountResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groups", req.Groups)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeDbKeyCountResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeDbKeyCountRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	Groups     string /*  DB编号，多个分组用逗号间隔  */
}

type Dcs2DescribeDbKeyCountResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeDbKeyCountReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeDbKeyCountReturnObjResponse struct {
	Total int32                                          `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeDbKeyCountReturnObjRowsResponse `json:"rows"`            /*  rows  */
}

type Dcs2DescribeDbKeyCountReturnObjRowsResponse struct {
	GroupName string `json:"groupName,omitempty"` /*  DB编号  */
	KeyNumber string `json:"keyNumber,omitempty"` /*  key数量  */
}
