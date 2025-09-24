package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetClusterApi
/* 调用该接口查询集群详情。
 */type CcseGetClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetClusterApi(client *core.CtyunClient) *CcseGetClusterApi {
	return &CcseGetClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetClusterApi) Do(ctx context.Context, credential core.Credential, req *CcseGetClusterRequest) (*CcseGetClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetClusterRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
}

type CcseGetClusterResponse struct {
	StatusCode int32                            `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                           `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetClusterReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                           `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetClusterReturnObjResponse struct {
	Id          string `json:"id,omitempty"`          /*  集群ID  */
	ClusterName string `json:"clusterName,omitempty"` /*  集群名称  */
	ClusterType int32  `json:"clusterType,omitempty"` /*  集群类型，表示如下
	0：专有版
	2：托管版  */
	ClusterDesc          string `json:"clusterDesc,omitempty"`          /*  集群描述  */
	DeployMode           string `json:"deployMode,omitempty"`           /*  集群部署模式  */
	ClusterVersion       string `json:"clusterVersion,omitempty"`       /*  集群版本  */
	ClusterStatus        string `json:"clusterStatus,omitempty"`        /*  集群状态，取值：<br />normal：正常。<br/>creating：创建中。<br/>create_fail：创建失败。<br/>adjust：规模调整中。<br/>updating：升级中。<br/>suspend：暂停。<br/>deleting：删除中。<br/>deleted：已删除。<br/>delete_fail：删除失败。<br/>resetting：节点重置中。<br/>resettled：节点已重置。<br/>reset_fail：节点重置失败。<br/>upgrading：集群升级中。<br/>upgrade_fail：集群升级失败。  */
	Ipv4Ipv6             string `json:"ipv4Ipv6,omitempty"`             /*  IP协议版本，取值：<br/>ipv4：IPv4版本。<br/>ipv6：IPv6版本。  */
	ControlPlaneProtocol string `json:"controlPlaneProtocol,omitempty"` /*  控制面接口协议，取值：<br/>ipv4：IPv4版本。<br/>ipv6：IPv6版本。  */
	MasterExtraVip       string `json:"masterExtraVip,omitempty"`       /*  master外网VIP地址  */
	MasterExtraVipv6     string `json:"masterExtraVipv6,omitempty"`     /*  master业务VIPv6地址  */
	MasterIntraVip       string `json:"masterIntraVip,omitempty"`       /*  master内网VIP地址  */
	MasterIntraVipv6     string `json:"masterIntraVipv6,omitempty"`     /*  master管理VIPv6地址  */
	SecurePort           int32  `json:"securePort,omitempty"`           /*  ApiServer安全端口  */
	NonSecurePort        int32  `json:"nonSecurePort,omitempty"`        /*  ApiServer非安全端口，0表示不开启非安全端口  */
	StartPort            int32  `json:"startPort,omitempty"`            /*  节点服务起始端口  */
	EndPort              int32  `json:"endPort,omitempty"`              /*  节点服务终止端口  */
	ServiceCidr          string `json:"serviceCidr,omitempty"`          /*  Service IP地址范围  */
	ServiceCidrv6        string `json:"serviceCidrv6,omitempty"`        /*  Service IPv6地址范围  */
	PodCidr              string `json:"podCidr,omitempty"`              /*  Pod IP地址范围  */
	PodCidrv6            string `json:"podCidrv6,omitempty"`            /*  Pod IPv6地址范围  */
	CreatedTime          string `json:"createdTime,omitempty"`          /*  创建时间  */
	ModifiedTime         string `json:"modifiedTime,omitempty"`         /*  修改时间  */
	MasterNodeNum        int32  `json:"masterNodeNum,omitempty"`        /*  master节点数量  */
	SlaveNodeNum         int32  `json:"slaveNodeNum,omitempty"`         /*  slave节点数量  */
	GrafanaAddress       string `json:"grafanaAddress,omitempty"`       /*  监控面板grafana地址  */
	BizState             int32  `json:"bizState,omitempty"`             /*  业务状态，1：运行中，2：已停止，3：已注销，4：已退订，5：扩容中，6：开通中，7：已取消，9：重启中，10：节点重置中，11：升级中，13：缩容中，14：已过期(冻结、过期)，15：节点升规格中，17：创建失败，18：退订中，19：控制面升配中，20：休眠中，21：唤醒中，22：转订购模式中  */
	ChannelLabel         string `json:"channelLabel,omitempty"`         /*  渠道标签  */
	ResPoolId            string `json:"resPoolId,omitempty"`            /*  资源池ID  */
	ResPoolName          string `json:"resPoolName,omitempty"`          /*  资源池名称  */
	Eip                  string `json:"eip,omitempty"`                  /*  集群绑定弹性ip  */
	Timezone             string `json:"timezone,omitempty"`             /*  时区  */
	ContainerRuntime     string `json:"containerRuntime,omitempty"`     /*  容器运行时  */
	NetworkPlugin        string `json:"networkPlugin,omitempty"`        /*  网络插件，包括calico、cubecni  */
	KubeProxyPattern     string `json:"kubeProxyPattern,omitempty"`     /*  kube-proxy 代理模式，包括ipvs、iptables  */
	ProdInstId           string `json:"prodInstId,omitempty"`           /*  实例ID  */
	ProdId               string `json:"prodId,omitempty"`               /*  集群规格编码  */
	ExpireTime           string `json:"expireTime,omitempty"`           /*  到期时间  */
	BillMode             string `json:"billMode,omitempty"`             /*  计费类型  */
	MasterSlbIp          string `json:"masterSlbIp,omitempty"`          /*  集群ApiServer Elb IP  */
	VpcId                string `json:"vpcId,omitempty"`                /*  集群所属VPC ID  */
	VpcName              string `json:"vpcName,omitempty"`              /*  集群所属VPC名称  */
	SubnetUuid           string `json:"subnetUuid,omitempty"`           /*  集群所属子网  */
	SecurityGroupName    string `json:"securityGroupName,omitempty"`    /*  安全组名称  */
	SecurityGroupId      string `json:"securityGroupId,omitempty"`      /*  安全组ID  */
	ClusterAlias         string `json:"clusterAlias,omitempty"`         /*  集群显示名称  */
}
