package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Dcs2DescribeDedicatedClusterInstanceListApi
/* 查询专属集群产品实例列表。
 */type Dcs2DescribeDedicatedClusterInstanceListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeDedicatedClusterInstanceListApi(client *core.CtyunClient) *Dcs2DescribeDedicatedClusterInstanceListApi {
	return &Dcs2DescribeDedicatedClusterInstanceListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeDedicatedClusterInstanceList",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeDedicatedClusterInstanceListApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeDedicatedClusterInstanceListRequest) (*Dcs2DescribeDedicatedClusterInstanceListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(req.PageIndex), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.InstanceName != "" {
		ctReq.AddParam("instanceName", req.InstanceName)
	}
	if req.ProdInstId != "" {
		ctReq.AddParam("prodInstId", req.ProdInstId)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeDedicatedClusterInstanceListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeDedicatedClusterInstanceListRequest struct {
	RegionId     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	PageIndex    int32  /*  当前页码  */
	PageSize     int32  /*  每页大小  */
	InstanceName string /*  实例名称  */
	ProdInstId   string /*  实例ID  */
}

type Dcs2DescribeDedicatedClusterInstanceListResponse struct {
	StatusCode int32                                                      `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                                     `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeDedicatedClusterInstanceListReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                                     `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                                     `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                                     `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeDedicatedClusterInstanceListReturnObjResponse struct {
	Total int32                                                            `json:"total,omitempty"` /*  总数  */
	Rows  []*Dcs2DescribeDedicatedClusterInstanceListReturnObjRowsResponse `json:"rows"`            /*  UserInfo  */
}

type Dcs2DescribeDedicatedClusterInstanceListReturnObjRowsResponse struct {
	ProdInstId        string `json:"prodInstId,omitempty"`        /*  实例ID  */
	User              string `json:"user,omitempty"`              /*  实例ID  */
	InstanceName      string `json:"instanceName,omitempty"`      /*  实例名称  */
	AccessSets        string `json:"accessSets,omitempty"`        /*  接入机集群名  */
	RedisSets         string `json:"redisSets,omitempty"`         /*  REDIS集群名  */
	Vip               string `json:"vip,omitempty"`               /*  vip地址  */
	Vipv6             string `json:"vipv6,omitempty"`             /*  IPV6格式的vip地址  */
	Status            int32  `json:"status,omitempty"`            /*  实例状态<li>0：有效<li>1：开通中<li>2：暂停<li>3：变更中<li>4：开通失败<li>5：停止中<li>6：已停止<li>8：已退订  */
	VipPort           int32  `json:"vipPort,omitempty"`           /*  vip访问端口  */
	CapacityInfo      string `json:"capacityInfo,omitempty"`      /*  容量信息  */
	PayType           int32  `json:"payType,omitempty"`           /*  付费类型<li>0：包年/包月<li>1: 按需  */
	PayTypeName       string `json:"payTypeName,omitempty"`       /*  付费类型名  */
	ElasticIpBind     int32  `json:"elasticIpBind,omitempty"`     /*  是否绑定弹性IP<li>0：未绑定<li>1：已绑定  */
	ElasticIp         string `json:"elasticIp,omitempty"`         /*  弹性IP  */
	OuterElasticIpId  string `json:"outerElasticIpId,omitempty"`  /*  弹性IP ID  */
	Whitelists        string `json:"whitelists,omitempty"`        /*  访问白名单  */
	Expiration        string `json:"expiration,omitempty"`        /*  过期时间  */
	EngineVersion     string `json:"engineVersion,omitempty"`     /*  引擎版本  */
	EngineVersionName string `json:"engineVersionName,omitempty"` /*  引擎版本名  */
	ArchType          string `json:"archType,omitempty"`          /*  架构类型<li>1：集群版<li>2：标准版<li>3：直连Cluster版<li>4：容量版<li>5：Proxy集群版  */
	NodeType          string `json:"nodeType,omitempty"`          /*  节点类型<li>1：双副本<li>2：单副本  */
	SecurityGroup     string `json:"securityGroup,omitempty"`     /*  安全组  */
	NetName           string `json:"netName,omitempty"`           /*  vpc网络名称  */
	Subnet            string `json:"subnet,omitempty"`            /*  子网名称  */
	CreateTime        string `json:"createTime,omitempty"`        /*  创建时间  */
	ExpTime           string `json:"expTime,omitempty"`           /*  过期时间  */
	ArchTypeName      string `json:"archTypeName,omitempty"`      /*  架构类型  */
	NodeTypeName      string `json:"nodeTypeName,omitempty"`      /*  节点类型名  */
	TplName           string `json:"tplName,omitempty"`           /*  模板名称  */
	TplCode           string `json:"tplCode,omitempty"`           /*  模板编码  */
	MaintenanceTime   string `json:"maintenanceTime,omitempty"`   /*  维护时间  */
	StatusName        string `json:"statusName,omitempty"`        /*  状态名称  */
	ProtectionStatus  bool   `json:"protectionStatus"`            /*  实例退订保护状态  */
	CopiesCount       string `json:"copiesCount,omitempty"`       /*  副本数  */
	ReadReplica       string `json:"readReplica,omitempty"`       /*  只读副本数  */
	ShardCount        string `json:"shardCount,omitempty"`        /*  分片数  */
	CpuArchType       string `json:"cpuArchType,omitempty"`       /*  cpu架构<li>x86<li>arm  */
	OuterProjectId    string `json:"outerProjectId,omitempty"`    /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br>说明：默认值为"0"  */
	OuterProjectName  string `json:"outerProjectName,omitempty"`  /*  项目名  */
	Description       string `json:"description,omitempty"`       /*  实例描述信息  */
}
