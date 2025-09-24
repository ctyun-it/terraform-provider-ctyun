package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeInstancesOverviewApi
/* 查询分布式缓存Redis实例基础详情。
 */type Dcs2DescribeInstancesOverviewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeInstancesOverviewApi(client *core.CtyunClient) *Dcs2DescribeInstancesOverviewApi {
	return &Dcs2DescribeInstancesOverviewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeInstancesOverview",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeInstancesOverviewApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeInstancesOverviewRequest) (*Dcs2DescribeInstancesOverviewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeInstancesOverviewResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeInstancesOverviewRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeInstancesOverviewResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeInstancesOverviewReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeInstancesOverviewReturnObjResponse struct {
	Total    int32                                                   `json:"total,omitempty"` /*  总数  */
	Rows     []*Dcs2DescribeInstancesOverviewReturnObjRowsResponse   `json:"rows"`            /*  节点信息redisSubSet  */
	UserInfo *Dcs2DescribeInstancesOverviewReturnObjUserInfoResponse `json:"userInfo"`        /*  UserInfo  */
	Nodes    []*Dcs2DescribeInstancesOverviewReturnObjNodesResponse  `json:"nodes"`           /*  MasterNode  */
}

type Dcs2DescribeInstancesOverviewReturnObjRowsResponse struct {
	UserName      string   `json:"userName,omitempty"`      /*  userName  */
	Name          string   `json:"name,omitempty"`          /*  节点名称  */
	AccessSetName string   `json:"accessSetName,omitempty"` /*  接入机名  */
	RedisSetNames []string `json:"redisSetNames"`           /*  redis集群名  */
}

type Dcs2DescribeInstancesOverviewReturnObjUserInfoResponse struct {
	DataDiskType  string `json:"dataDiskType"`
	ShardMemSize  string `json:"shardMemSize"`
	PaasInstAttrs []struct {
		AttrKey string `json:"attrKey"`
		AttrVal string `json:"attrVal"`
	} `json:"paasInstAttrs"`
	ProdInstId            string                                                          `json:"prodInstId,omitempty"`            /*  实例ID  */
	User                  string                                                          `json:"user,omitempty"`                  /*  产品实例标识  */
	InstanceName          string                                                          `json:"instanceName,omitempty"`          /*  实例名称  */
	AccessSets            string                                                          `json:"accessSets,omitempty"`            /*  接入机集群名  */
	RedisSets             string                                                          `json:"redisSets,omitempty"`             /*  REDIS集群名  */
	RedisSubSet           []string                                                        `json:"redisSubSet"`                     /*  redisSubSet  */
	Vip                   string                                                          `json:"vip,omitempty"`                   /*  vip地址  */
	Vipv6                 string                                                          `json:"vipv6,omitempty"`                 /*  vipv6  */
	Status                int32                                                           `json:"status,omitempty"`                /*  实例状态<li>0：有效<li>1：开通中<li>2：暂停<li>3：变更中<li>4：开通失败<li>5：停止中<li>6：已停止<li>8：已退订  */
	StatusName            string                                                          `json:"statusName,omitempty"`            /*  状态英文名  */
	ProtectionStatus      *bool                                                           `json:"protectionStatus"`                /*  是否开启实例退订保护  */
	VipPort               int32                                                           `json:"vipPort,omitempty"`               /*  所属资源池  */
	CapacityInfo          string                                                          `json:"capacityInfo,omitempty"`          /*  容量信息  */
	Capacity              string                                                          `json:"capacity,omitempty"`              /*  实例规格容量 单位G  */
	ShardCount            string                                                          `json:"shardCount,omitempty"`            /*  分片数  */
	PayType               int32                                                           `json:"payType,omitempty"`               /*  付费类型<li>0：包年/包月<li>1: 按需  */
	PayTypeName           string                                                          `json:"payTypeName,omitempty"`           /*  付费类型  */
	ElasticIpBind         int32                                                           `json:"elasticIpBind,omitempty"`         /*  是否绑定弹性IP<li>0:未绑定<li>1：已绑定  */
	ElasticIp             string                                                          `json:"elasticIp,omitempty"`             /*  弹性IP  */
	OuterElasticIpId      string                                                          `json:"outerElasticIpId,omitempty"`      /*  弹性IP ID  */
	ConnectionAddress     string                                                          `json:"connectionAddress,omitempty"`     /*  连接地址  */
	Ipv6ConnectionAddress string                                                          `json:"ipv6ConnectionAddress,omitempty"` /*  ipv6连接地址  */
	Whitelists            string                                                          `json:"whitelists,omitempty"`            /*  访问白名单  */
	Expiration            string                                                          `json:"expiration,omitempty"`            /*  过期时间  */
	EngineVersion         string                                                          `json:"engineVersion,omitempty"`         /*  引擎版本  */
	EngineVersionName     string                                                          `json:"engineVersionName,omitempty"`     /*  引擎版本名  */
	ArchType              string                                                          `json:"archType,omitempty"`              /*  架构类型<li>1：集群版<li>2：标准版<li>3：直连Cluster版<li>4：容量版<li>5：Proxy集群版  */
	NodeType              string                                                          `json:"nodeType,omitempty"`              /*  节点类型<li>1：双副本<li>2：单副本  */
	SecurityGroup         string                                                          `json:"securityGroup,omitempty"`         /*  安全组  */
	NetName               string                                                          `json:"netName,omitempty"`               /*  vpc网络名称  */
	Subnet                string                                                          `json:"subnet,omitempty"`                /*  子网名称  */
	CreateTime            string                                                          `json:"createTime,omitempty"`            /*  创建时间  */
	ExpTime               string                                                          `json:"expTime,omitempty"`               /*  过期时间  */
	ArchTypeName          string                                                          `json:"archTypeName,omitempty"`          /*  架构类型  */
	NodeTypeName          string                                                          `json:"nodeTypeName,omitempty"`          /*  节点类型名  */
	TplName               string                                                          `json:"tplName,omitempty"`               /*  模板名称  */
	TplCode               string                                                          `json:"tplCode,omitempty"`               /*  模板编码  */
	MaintenanceTime       string                                                          `json:"maintenanceTime,omitempty"`       /*  维护时间  */
	AzList                []*Dcs2DescribeInstancesOverviewReturnObjUserInfoAzListResponse `json:"azList"`                          /*  可用区列表  */
	Description           string                                                          `json:"description,omitempty"`           /*  实例描述信息  */
	EnableMultiRead       *bool                                                           `json:"enableMultiRead"`                 /*  是否开启读写分离  */
	CopiesCount           string                                                          `json:"copiesCount,omitempty"`           /*  副本数  */
	ReadReplica           string                                                          `json:"readReplica,omitempty"`           /*  只读副本数  */
	AutoScaleFlag         *bool                                                           `json:"autoScaleFlag"`                   /*  是否开启弹性伸缩  */
	CpuArchType           string                                                          `json:"cpuArchType,omitempty"`           /*  cpu架构<li>x86<li>arm  */
}

type Dcs2DescribeInstancesOverviewReturnObjNodesResponse struct {
	MasterName      string                                                           `json:"masterName,omitempty"`      /*  节点名称  */
	StartSlot       string                                                           `json:"startSlot,omitempty"`       /*  起始槽位  */
	EndSlot         string                                                           `json:"endSlot,omitempty"`         /*  结束槽位  */
	ConnUrl         string                                                           `json:"connUrl,omitempty"`         /*  管理IP连接地址  */
	Status          string                                                           `json:"status,omitempty"`          /*  节点状态<li>CACHE.COMM.STATUS：正常<li>CACHE.DIAT.PREP：扩容数据准备<li>CACHE.DIAT.PROCESS：执行扩容数据<li>CACHE.DIAT.DEL：删除<li>CACHE.PROB.SWIT：故障待切换  */
	Enabled         *bool                                                            `json:"enabled"`                   /*  是否可用  */
	SlaveNodes      []*Dcs2DescribeInstancesOverviewReturnObjNodesSlaveNodesResponse `json:"slaveNodes"`                /*  从节点  */
	VpcUrl          string                                                           `json:"vpcUrl,omitempty"`          /*  vpc地址  */
	SpuInstDeployId string                                                           `json:"spuInstDeployId,omitempty"` /*  模块ID  */
	AzId            string                                                           `json:"azId,omitempty"`            /*  可用区编码  */
	AzName          string                                                           `json:"azName,omitempty"`          /*  可用区名称  */
}

type Dcs2DescribeInstancesOverviewReturnObjUserInfoAzListResponse struct {
	AzId      string `json:"azId,omitempty"`      /*  可用区ID  */
	AzName    string `json:"azName,omitempty"`    /*  可用区名称  */
	AzEngName string `json:"azEngName,omitempty"` /*  可用区英文名称  */
}

type Dcs2DescribeInstancesOverviewReturnObjNodesSlaveNodesResponse struct {
	SlaveName       string `json:"slaveName,omitempty"`       /*  slaveName  */
	ConnUrl         string `json:"connUrl,omitempty"`         /*  管理IP连接地址  */
	Status          string `json:"status,omitempty"`          /*  节点状态，  */
	Enabled         *bool  `json:"enabled"`                   /*  是否可用  */
	VpcUrl          string `json:"vpcUrl,omitempty"`          /*  vpc地址  */
	SpuInstDeployId string `json:"spuInstDeployId,omitempty"` /*  模块ID  */
	AzId            string `json:"azId,omitempty"`            /*  可用区编码  */
	AzName          string `json:"azName,omitempty"`          /*  可用区名称  */
}
