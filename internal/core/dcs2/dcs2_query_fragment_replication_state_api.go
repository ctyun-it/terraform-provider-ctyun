package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2QueryFragmentReplicationStateApi
/* 查询分片副本状态。
 */type Dcs2QueryFragmentReplicationStateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryFragmentReplicationStateApi(client *core.CtyunClient) *Dcs2QueryFragmentReplicationStateApi {
	return &Dcs2QueryFragmentReplicationStateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/queryFragmentReplicationState",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryFragmentReplicationStateApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryFragmentReplicationStateRequest) (*Dcs2QueryFragmentReplicationStateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("fragmentName", req.FragmentName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryFragmentReplicationStateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryFragmentReplicationStateRequest struct {
	RegionId     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId   string /*  实例 ID  */
	FragmentName string /*  分片名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7740&isNormal=1&vid=270">查询集群节点信息</a> RedisNode表fragmentName字段  */
}

type Dcs2QueryFragmentReplicationStateResponse struct {
	StatusCode int32                                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2QueryFragmentReplicationStateReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2QueryFragmentReplicationStateReturnObjResponse struct {
	Role   string `json:"role,omitempty"`   /*  角色<li>master：主节点<li>slave： 从节点  */
	VpcIp  string `json:"vpcIp,omitempty"`  /*  vpc ip  */
	Status string `json:"status,omitempty"` /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	AzName string `json:"azName,omitempty"` /*  可用区名称  */
}
