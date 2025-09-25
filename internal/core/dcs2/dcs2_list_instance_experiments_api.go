package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ListInstanceExperimentsApi
/* 查询分布式缓存Redis实例故障列表
 */type Dcs2ListInstanceExperimentsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ListInstanceExperimentsApi(client *core.CtyunClient) *Dcs2ListInstanceExperimentsApi {
	return &Dcs2ListInstanceExperimentsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/inject/listInstanceExperiments",
			ContentType:  "",
		},
	}
}

func (a *Dcs2ListInstanceExperimentsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ListInstanceExperimentsRequest) (*Dcs2ListInstanceExperimentsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.Page != "" {
		ctReq.AddParam("page", req.Page)
	}
	if req.Size != "" {
		ctReq.AddParam("size", req.Size)
	}
	if req.ActionCode != "" {
		ctReq.AddParam("actionCode", req.ActionCode)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ListInstanceExperimentsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ListInstanceExperimentsRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
	Page       string /*  当前页,最小值为1, 默认1  */
	Size       string /*  当前页大小,范围[1, 100]，默认10  */
	ActionCode string /*  故障类型<li>memory-load: 内存负载<li>cpu-fullload：CPU满载注入<li>disk-burn: 磁盘IO Utilization<li>node-shutdown: 主机宕机<li>network-loss: 网络丢包  */
}

type Dcs2ListInstanceExperimentsResponse struct {
	StatusCode int32                                         `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                        `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ListInstanceExperimentsReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                        `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                        `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                        `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ListInstanceExperimentsReturnObjResponse struct {
	Total string `json:"total,omitempty"` /*  总条数  */
	List  string `json:"list,omitempty"`  /*  Array of Objects  */
}
