package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2DescribeInstancesClusterMemberInfoApi
/* 批量查询租户的缓存redis集群节点信息。
 */type Dcs2DescribeInstancesClusterMemberInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeInstancesClusterMemberInfoApi(client *core.CtyunClient) *Dcs2DescribeInstancesClusterMemberInfoApi {
	return &Dcs2DescribeInstancesClusterMemberInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeInstancesClusterMemberInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeInstancesClusterMemberInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeInstancesClusterMemberInfoRequest) (*Dcs2DescribeInstancesClusterMemberInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.ProjectId != "" {
		ctReq.AddParam("projectId", req.ProjectId)
	}
	if req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(req.PageIndex), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeInstancesClusterMemberInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeInstancesClusterMemberInfoRequest struct {
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProjectId string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br>说明：默认值为"0"  */
	PageIndex int32  /*  当前页码  */
	PageSize  int32  /*  每页大小  */
}

type Dcs2DescribeInstancesClusterMemberInfoResponse struct {
	StatusCode int32                                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeInstancesClusterMemberInfoReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeInstancesClusterMemberInfoReturnObjResponse struct {
	Total     int32                                                               `json:"total,omitempty"` /*  响应码描述  */
	Instances []*Dcs2DescribeInstancesClusterMemberInfoReturnObjInstancesResponse `json:"instances"`       /*  redis集群列表  */
}

type Dcs2DescribeInstancesClusterMemberInfoReturnObjInstancesResponse struct {
	ProdInstId string                                                                  `json:"prodInstId,omitempty"` /*  实例ID  */
	Total      int32                                                                   `json:"total,omitempty"`      /*  响应码描述  */
	Rows       []*Dcs2DescribeInstancesClusterMemberInfoReturnObjInstancesRowsResponse `json:"rows"`                 /*  返回redis集群详细信息  */
}

type Dcs2DescribeInstancesClusterMemberInfoReturnObjInstancesRowsResponse struct {
	RedisSetName string                                                                       `json:"redisSetName,omitempty"` /*  redis集群名  */
	RedisSetInfo string                                                                       `json:"RedisSetInfo,omitempty"` /*  redis集群节点信息  */
	LastestTime  string                                                                       `json:"lastestTime,omitempty"`  /*  更新时间戳  */
	IsRWsep      string                                                                       `json:"isRWsep,omitempty"`      /*  是否读写分离<li>true：开启<li>false：关闭  */
	Nodes        []*Dcs2DescribeInstancesClusterMemberInfoReturnObjInstancesRowsNodesResponse `json:"nodes"`                  /*  redis节点集合，数据见RedisNode  */
	SlotInfo     string                                                                       `json:"slotInfo,omitempty"`     /*  槽位信息  */
	IsAuth       string                                                                       `json:"isAuth,omitempty"`       /*  是否加密鉴权  */
}

type Dcs2DescribeInstancesClusterMemberInfoReturnObjInstancesRowsNodesResponse struct {
	MasterName   string `json:"masterName,omitempty"`   /*  主节点名称  */
	FragmentName string `json:"fragmentName,omitempty"` /*  分片名称  */
	StartSlot    string `json:"startSlot,omitempty"`    /*  开始槽位  */
	EndSlot      string `json:"endSlot,omitempty"`      /*  结束槽位  */
	ConnUrl      string `json:"connUrl,omitempty"`      /*  连接地址  */
	AzName       string `json:"azName,omitempty"`       /*  可用区名称  */
	Status       string `json:"status,omitempty"`       /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	Enabled      *bool  `json:"enabled"`                /*  是否可用状态  */
	VpcUrl       string `json:"vpcUrl,omitempty"`       /*  vpc连接地址  */
}
