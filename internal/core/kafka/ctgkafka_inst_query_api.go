package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaInstQueryApi
/* 查询实例。
 */type CtgkafkaInstQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaInstQueryApi(client *core.CtyunClient) *CtgkafkaInstQueryApi {
	return &CtgkafkaInstQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaInstQueryApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaInstQueryRequest) (*CtgkafkaInstQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.ProdInstId != "" {
		ctReq.AddParam("prodInstId", req.ProdInstId)
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	if req.ExactMatchName != nil {
		ctReq.AddParam("exactMatchName", strconv.FormatBool(*req.ExactMatchName))
	}
	if req.Status != 0 {
		ctReq.AddParam("status", strconv.FormatInt(int64(req.Status), 10))
	}
	if req.OuterProjectId != "" {
		ctReq.AddParam("outerProjectId", req.OuterProjectId)
	}
	if req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(req.PageNum), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaInstQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaInstQueryRequest struct {
	RegionId       string `json:"regionId,omitempty"`       /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId     string `json:"prodInstId,omitempty"`     /*  实例ID，如果传入，则返回指定实例的详细信息。  */
	Name           string `json:"name,omitempty"`           /*  实例名称。  */
	ExactMatchName *bool  `json:"exactMatchName"`           /*  是否精确匹配实例名称，当name有值时本参数有效。<br><li>true：精确查询<br><li>false：模糊查询（默认）  */
	Status         int32  `json:"status,omitempty"`         /*  实例状态：<br><li>1：运行中<br><li>2：已过期<br><li>3：已注销<br><li>4：变更中<br><li>5：已退订<br><li>6：开通中<br><li>7：已取消<br><li>8：已停止<br><li>9：弹性IP处理中<br><li>10：重启中<br><li>11：重启失败<br><li>12：升级中<br><li>13：已欠费<br><li>101：开通失败  */
	OuterProjectId string `json:"outerProjectId,omitempty"` /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目。  */
	PageNum        int32  `json:"pageNum,omitempty"`        /*  分页中的页数，默认1，范围1-40000。  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaInstQueryResponse struct {
	StatusCode string                              `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                              `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaInstQueryReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaInstQueryReturnObjResponse struct {
	Total int32                                     `json:"total,omitempty"` /*  总记录数。  */
	Data  []*CtgkafkaInstQueryReturnObjDataResponse `json:"data"`            /*  实例信息。  */
}

type CtgkafkaInstQueryReturnObjDataResponse struct {
	ProdInstId       string                                            `json:"prodInstId,omitempty"`       /*  实例ID。  */
	InstanceName     string                                            `json:"instanceName,omitempty"`     /*  实例名称。  */
	Status           int32                                             `json:"status,omitempty"`           /*  状态：<br><li>1：正常<br><li>2：暂停<br><li>3：注销  */
	MqEngineType     string                                            `json:"mqEngineType,omitempty"`     /*  产品引擎类型：<br><li>kafka<br><li>pulsar  */
	PartitionNum     int32                                             `json:"partitionNum,omitempty"`     /*  分区数上限。  */
	Subnet           string                                            `json:"subnet,omitempty"`           /*  子网名称。  */
	Space            string                                            `json:"space,omitempty"`            /*  单个代理（节点）磁盘存储空间，单位GB。  */
	SecurityGroup    string                                            `json:"securityGroup,omitempty"`    /*  安全组ID。  */
	Network          string                                            `json:"network,omitempty"`          /*  网络名称。  */
	BillMode         string                                            `json:"billMode,omitempty"`         /*  付费类型。<br><li>1：包周期<br><li>2：按需  */
	ExpireTime       string                                            `json:"expireTime,omitempty"`       /*  过期时间。  */
	CreateTime       string                                            `json:"createTime,omitempty"`       /*  创建时间。  */
	UpdateTime       string                                            `json:"updateTime,omitempty"`       /*  更新时间。  */
	NodeList         []*CtgkafkaInstQueryReturnObjDataNodeListResponse `json:"nodeList"`                   /*  代理（节点）列表信息，当请求参数传入prodInstId时有值  */
	Version          string                                            `json:"version,omitempty"`          /*  实例版本。  */
	RegionCode       string                                            `json:"regionCode,omitempty"`       /*  资源池编码。  */
	RegionName       string                                            `json:"regionName,omitempty"`       /*  资源池名称。  */
	OuterProjectId   string                                            `json:"outerProjectId,omitempty"`   /*  企业项目ID。  */
	OuterProjectName string                                            `json:"outerProjectName,omitempty"` /*  企业项目名称。  */
	DomainEndpoint   string                                            `json:"domainEndpoint,omitempty"`   /*  实例域名接入点。  */
	Remark           string                                            `json:"remark,omitempty"`           /*  实例备注。  */
	ClusterType      int32                                             `json:"clusterType,omitempty"`      /*  实例类型。<br><li>1：单机版<br><li>2：集群版  */
	Specifications   string                                            `json:"specifications,omitempty"`   /*  实例规格。  */
	VpcId            string                                            `json:"vpcId,omitempty"`            /*  VPC ID。  */
	SubnetId         string                                            `json:"subnetId,omitempty"`         /*  子网ID。  */
	AutoReassign     int32                                             `json:"autoReassign,omitempty"`     /*  扩容时是否自动迁移主题。 <br><li>1：是<br><li>2：否<br><li>空值表示不开启  */
	Ipv6Enable       int32                                             `json:"ipv6Enable,omitempty"`       /*  是否开启IPV6。 <br><li>1：是<br><li>2：否<br><li>空值表示不开启  */
	ElasticEnable    string                                            `json:"elasticEnable,omitempty"`    /*  是否开启弹性带宽。 <br><li>1：是<br><li>2：否<br><li>空值表示不开启  */
	Vip              string                                            `json:"vip,omitempty"`              /*  实例vip。  */
	Vipv6            string                                            `json:"vipv6,omitempty"`            /*  实例vipv6。  */
	DiskType         string                                            `json:"diskType,omitempty"`         /*  磁盘类型。  */
	Protocols        *CtgkafkaInstQueryReturnObjDataProtocolsResponse  `json:"protocols"`                  /*  实例接入点信息，当请求参数传入prodInstId时有值。  */
}

type CtgkafkaInstQueryReturnObjDataNodeListResponse struct {
	TenantId       string `json:"tenantId,omitempty"`       /*  租户ID。  */
	TenantName     string `json:"tenantName,omitempty"`     /*  租户名称。  */
	VpcIp          string `json:"vpcIp,omitempty"`          /*  节点VPC IP地址。  */
	VpcIpv6        string `json:"vpcIpv6,omitempty"`        /*  节点VPC IP地址。  */
	ElasticIp      string `json:"elasticIp,omitempty"`      /*  弹性IP，绑定后有值。  */
	ElasticIpState int32  `json:"elasticIpState,omitempty"` /*  弹性IP绑定状态。<br><li>0：未绑定<br><li>1：绑定成功<br><li>2：绑定失败<br><li>3：处理中<br><li>4：解绑失败  */
	CreateTime     string `json:"createTime,omitempty"`     /*  节点创建时间。  */
	UpdateTime     string `json:"updateTime,omitempty"`     /*  节点更新时间。  */
	ServerSeq      int32  `json:"serverSeq,omitempty"`      /*  节点ID。  */
	State          int32  `json:"state,omitempty"`          /*  节点专题。<br><li>0：停止<br><li>1：启动  */
	AzId           string `json:"azId,omitempty"`           /*  节点所在可用区ID。  */
	AzName         string `json:"azName,omitempty"`         /*  节点所在可用区名称。  */
	EcsId          string `json:"ecsId,omitempty"`          /*  节点ID。  */
	VpcPort        string `json:"vpcPort,omitempty"`        /*  公共接入点（PLAINTEXT）端口。  */
	SaslPort       string `json:"saslPort,omitempty"`       /*  安全接入点（SASL_PLAINTEXT）端口。  */
	ListenNodePort string `json:"listenNodePort,omitempty"` /*  SSL接入点(SASL_SSL)端口。  */
	Ipv6PlainPort  string `json:"ipv6PlainPort,omitempty"`  /*  IPV6公共接入点（PLAINTEXT）端口。  */
	Ipv6SaslPort   string `json:"ipv6SaslPort,omitempty"`   /*  IPV6安全接入点（SASL_PLAINTEXT）端口。  */
	Ipv6SslPort    string `json:"ipv6SslPort,omitempty"`    /*  IPV6 SSL接入点(SASL_SSL)端口。  */
	HttpPort       string `json:"httpPort,omitempty"`       /*  HTTP接入点端口。  */
	SaslNodePort   string `json:"saslNodePort,omitempty"`   /*  弹性IP安全接入点（SASL_PLAINTEXT）端口。  */
	DomainEndpoint string `json:"domainEndpoint,omitempty"` /*  域名。  */
}

type CtgkafkaInstQueryReturnObjDataProtocolsResponse struct {
	PlainTextAddr         string `json:"plainTextAddr,omitempty"`         /*  公共接入点。  */
	PlainTextIpv6Addr     string `json:"plainTextIpv6Addr,omitempty"`     /*  IPv6公共接入点。  */
	SaslPlainTextAddr     string `json:"saslPlainTextAddr,omitempty"`     /*  安全接入点。  */
	SaslPlainTextIpv6Addr string `json:"saslPlainTextIpv6Addr,omitempty"` /*  IPv6安全接入点。  */
	SaslSslAddr           string `json:"saslSslAddr,omitempty"`           /*  SSL安全接入点。  */
	SaslSslIpv6Addr       string `json:"saslSslIpv6Addr,omitempty"`       /*  IPv6 SSL安全接入点。  */
	HttpAddr              string `json:"httpAddr,omitempty"`              /*  HTTP接入点。  */
	HttpIpv6Addr          string `json:"httpIpv6Addr,omitempty"`          /*  IPv6 HTTP接入点。  */
	PublicSaslTextAddr    string `json:"publicSaslTextAddr,omitempty"`    /*  公网接入点。  */
	DomainEndpointAddr    string `json:"domainEndpointAddr,omitempty"`    /*  域名接入点。  */
	EnableIPv6            *bool  `json:"enableIPv6"`                      /*  是否开启IPv6  */
}
