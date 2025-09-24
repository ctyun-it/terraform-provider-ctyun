package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeMemoryInfoApi
/* 查询Redis节点内存
 */type Dcs2DescribeMemoryInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeMemoryInfoApi(client *core.CtyunClient) *Dcs2DescribeMemoryInfoApi {
	return &Dcs2DescribeMemoryInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeMemoryInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeMemoryInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeMemoryInfoRequest) (*Dcs2DescribeMemoryInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("ip", req.Ip)
	ctReq.AddParam("port", req.Port)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeMemoryInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeMemoryInfoRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	Ip         string /*  redis节点Ip  */
	Port       string /*  redis节点端口  */
}

type Dcs2DescribeMemoryInfoResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeMemoryInfoReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeMemoryInfoReturnObjResponse struct {
	FreeMemory string `json:"freeMemory,omitempty"` /*  可用内存  */
	MaxMemory  string `json:"maxMemory,omitempty"`  /*  最大内存  */
}
