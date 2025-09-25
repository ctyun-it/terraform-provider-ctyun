package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeClusterMemberInfoApi
/* 查询集群节点信息
 */type Dcs2DescribeClusterMemberInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeClusterMemberInfoApi(client *core.CtyunClient) *Dcs2DescribeClusterMemberInfoApi {
	return &Dcs2DescribeClusterMemberInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeClusterMemberInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeClusterMemberInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeClusterMemberInfoRequest) (*Dcs2DescribeClusterMemberInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeClusterMemberInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeClusterMemberInfoRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeClusterMemberInfoResponse struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                       `json:"message,omitempty"`    /*  响应信息  */
	RequestId  string                                       `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                       `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
	Total      int32                                        `json:"total,omitempty"`      /*  响应码描述  */
	Rows       []*Dcs2DescribeClusterMemberInfoRowsResponse `json:"rows"`                 /*  返回redis集群数据集，数据见rows  */
}

type Dcs2DescribeClusterMemberInfoRowsResponse struct {
	RedisSetName       string                                            `json:"redisSetName,omitempty"`       /*  redis集群名  */
	RedisSetInfo       string                                            `json:"RedisSetInfo,omitempty"`       /*  redis集群节点信息  */
	LastestTime        string                                            `json:"lastestTime,omitempty"`        /*  更新时间戳  */
	IsRWsep            string                                            `json:"isRWsep,omitempty"`            /*  是否读写分离<li>true：开启<li>false：关闭  */
	Nodes              []*Dcs2DescribeClusterMemberInfoRowsNodesResponse `json:"nodes"`                        /*  redis节点集合，数据见RedisNode  */
	SlotInfo           string                                            `json:"slotInfo,omitempty"`           /*  槽位信息  */
	Requirepass        string                                            `json:"Requirepass,omitempty"`        /*  redis节点密码  */
	RequirepassEncrypt string                                            `json:"requirepassEncrypt,omitempty"` /*  redis节点密文  */
	IsAuth             string                                            `json:"isAuth,omitempty"`             /*  是否加密鉴权  */
}

type Dcs2DescribeClusterMemberInfoRowsNodesResponse struct {
	MasterName      string                                                      `json:"masterName,omitempty"`      /*  主节点名称  */
	StartSlot       string                                                      `json:"startSlot,omitempty"`       /*  开始槽位  */
	EndSlot         string                                                      `json:"endSlot,omitempty"`         /*  结束槽位  */
	ConnUrl         string                                                      `json:"connUrl,omitempty"`         /*  连接地址  */
	Status          string                                                      `json:"status,omitempty"`          /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	Enabled         *bool                                                       `json:"enabled"`                   /*  是否可用状态  */
	FragmentName    string                                                      `json:"fragmentName,omitempty"`    /*  分片名称  */
	SlaveNodes      []*Dcs2DescribeClusterMemberInfoRowsNodesSlaveNodesResponse `json:"slaveNodes"`                /*  从节点数组，数据见SlaveNode  */
	VpcUrl          string                                                      `json:"vpcUrl,omitempty"`          /*  vpc连接地址  */
	SpuInstDeployId string                                                      `json:"spuInstDeployId,omitempty"` /*  节点ID  */
	AzId            string                                                      `json:"azId,omitempty"`            /*  可用区ID  */
	AzName          string                                                      `json:"azName,omitempty"`          /*  可用区名称  */
	NodeVpcIp       string                                                      `json:"nodeVpcIp,omitempty"`       /*  节点VPCIP  */
}

type Dcs2DescribeClusterMemberInfoRowsNodesSlaveNodesResponse struct {
	SlaveName       string `json:"slaveName,omitempty"`       /*  redis从节点名称  */
	ConnUrl         string `json:"connUrl,omitempty"`         /*  连接地址  */
	Status          string `json:"status,omitempty"`          /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	Enabled         *bool  `json:"enabled"`                   /*  是否可用状态  */
	VpcUrl          string `json:"vpcUrl,omitempty"`          /*  从节点vpc连接地址  */
	SpuInstDeployId string `json:"spuInstDeployId,omitempty"` /*  节点ID  */
	AzId            string `json:"azId,omitempty"`            /*  可用区ID  */
	AzName          string `json:"azName,omitempty"`          /*  可用区名称  */
	NodeVpcIp       string `json:"nodeVpcIp,omitempty"`       /*  节点VpcIp  */
}
