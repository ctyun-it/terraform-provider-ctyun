package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeHotKeysApi
/* 创建查询分布式缓存Redis实例的热Key任务，当前key请求次数超过3000，则认为是热key，返回任务ID。
 */type Dcs2DescribeHotKeysApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeHotKeysApi(client *core.CtyunClient) *Dcs2DescribeHotKeysApi {
	return &Dcs2DescribeHotKeysApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/keyAnalysisMgrServant/describeHotKeys",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeHotKeysApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeHotKeysRequest) (*Dcs2DescribeHotKeysResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeHotKeysResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeHotKeysRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeHotKeysResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeHotKeysReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeHotKeysReturnObjResponse struct {
	TaskId string `json:"taskId,omitempty"` /*  任务ID  */
}
