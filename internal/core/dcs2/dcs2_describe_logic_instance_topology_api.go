package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeLogicInstanceTopologyApi
/* 查询分布式缓存Redis实例的逻辑拓扑结构。
 */type Dcs2DescribeLogicInstanceTopologyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeLogicInstanceTopologyApi(client *core.CtyunClient) *Dcs2DescribeLogicInstanceTopologyApi {
	return &Dcs2DescribeLogicInstanceTopologyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeLogicInstanceTopology",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeLogicInstanceTopologyApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeLogicInstanceTopologyRequest) (*Dcs2DescribeLogicInstanceTopologyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeLogicInstanceTopologyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeLogicInstanceTopologyRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeLogicInstanceTopologyResponse struct {
	StatusCode int32                                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeLogicInstanceTopologyReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeLogicInstanceTopologyReturnObjResponse struct {
	RedisNodes  []*Dcs2DescribeLogicInstanceTopologyReturnObjRedisNodesResponse  `json:"redisNodes"`  /*  redis节点集合，见RedisNode  */
	AccessNodes []*Dcs2DescribeLogicInstanceTopologyReturnObjAccessNodesResponse `json:"accessNodes"` /*  接入机节点集合，见AccessNode  */
}

type Dcs2DescribeLogicInstanceTopologyReturnObjRedisNodesResponse struct {
	ConnUrl    string                                                                `json:"connUrl,omitempty"`    /*  连接地址  */
	StartSlot  string                                                                `json:"startSlot,omitempty"`  /*  开始槽位  */
	EndSlot    string                                                                `json:"endSlot,omitempty"`    /*  结束槽位  */
	VpcUrl     string                                                                `json:"vpcUrl,omitempty"`     /*  vpc连接地址  */
	Status     string                                                                `json:"status,omitempty"`     /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	MasterName string                                                                `json:"masterName,omitempty"` /*  主节点名称  */
	Slaves     []*Dcs2DescribeLogicInstanceTopologyReturnObjRedisNodesSlavesResponse `json:"slaves"`               /*  从节点集合，见SlaveNode  */
}

type Dcs2DescribeLogicInstanceTopologyReturnObjAccessNodesResponse struct {
	ClientConnNum    int32  `json:"clientConnNum,omitempty"`    /*  客户端连接数  */
	ClientMaxConnNum int32  `json:"clientMaxConnNum,omitempty"` /*  接入机最大并发连接数  */
	ConnNumPerRedis  string `json:"connNumPerRedis,omitempty"`  /*  Redis节点每秒连接数  */
	ConnUrl          string `json:"connUrl,omitempty"`          /*  连接地址  */
	ProxyName        string `json:"proxyName,omitempty"`        /*  代理名称  */
	VpcUrl           string `json:"vpcUrl,omitempty"`           /*  vpc连接地址  */
}

type Dcs2DescribeLogicInstanceTopologyReturnObjRedisNodesSlavesResponse struct {
	SlaveName string `json:"slaveName,omitempty"` /*  redis从节点名称  */
	ConnUrl   string `json:"connUrl,omitempty"`   /*  连接地址  */
	Status    string `json:"status,omitempty"`    /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	Enabled   string `json:"enabled,omitempty"`   /*  是否可用状态  */
	VpcUrl    string `json:"vpcUrl,omitempty"`    /*  从节点vpc连接地址  */
}
