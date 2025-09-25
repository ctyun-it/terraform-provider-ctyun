package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseCreateClusterApi
/* 调用该接口创建Kubernetes集群
 */type CcseCreateClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseCreateClusterApi(client *core.CtyunClient) *CcseCreateClusterApi {
	return &CcseCreateClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseCreateClusterApi) Do(ctx context.Context, credential core.Credential, req *CcseCreateClusterRequest) (*CcseCreateClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseCreateClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseCreateClusterRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	ResPoolId string `json:"resPoolId,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>  */
	ClusterBaseInfo *CcseCreateClusterClusterBaseInfoRequest `json:"clusterBaseInfo"` /*  集群基本信息  */
	MasterHost      *CcseCreateClusterMasterHostRequest      `json:"masterHost"`      /*  master节点基本信息，专有版必填，托管版时不传  */
	SlaveHost       *CcseCreateClusterSlaveHostRequest       `json:"slaveHost"`       /*  worker节点基本信息  */
}

type CcseCreateClusterClusterBaseInfoRequest struct {
	SubnetUuid string `json:"subnetUuid,omitempty"` /*  子网id，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&amp;api=8659&amp;data=94">查询子网列表</a>
	注：在多可用区类型资源池下，subnetID通常以“subnet-”开头，非多可用区类型资源池subnetID为uuid格式  */
	NetworkPlugin string `json:"networkPlugin,omitempty"` /*  网络插件，可选calico和cubecni。您可查看<a href="https://www.ctyun.cn/document/10083472/10520760">容器网络插件说明</a>选择
	注：calico需要申请白名单  */
	ClusterDomain     string   `json:"clusterDomain,omitempty"` /*  集群本地域名，集群创建完成后，暂无法修改，请慎重配置  */
	PodSubnetUuidList []string `json:"podSubnetUuidList"`       /*  pod子网id列表，网络插件选择cubecni必传
	 */
	AutoGenerateSecurityGroup *bool  `json:"autoGenerateSecurityGroup"`   /*  是否自动创建安全组，默认false  */
	SecurityGroupUuid         string `json:"securityGroupUuid,omitempty"` /*  安全组id，autoGenerateSecurityGroup=false必填
	安全组id您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94&vid=88">查询用户安全组列表</a>
	注：在多可用区类型资源池下，安全组ID通常以“sg-”开头，非多可用区类型资源池安全组ID为uuid格式  */
	StartPort          int64  `json:"startPort,omitempty"`      /*  节点服务开始端口，可选范围30000-65535  */
	EndPort            int64  `json:"endPort,omitempty"`        /*  节点服务终止端口，可选范围30000-65535  */
	EnableApiServerEip *bool  `json:"enableApiServerEip"`       /*  是否开启ApiServerEip，默认false，若开启将自动创建按需计费类型的eip。  */
	EnableSnat         *bool  `json:"enableSnat"`               /*  是否开启nat网关，默认false，若开启将自动创建按需计费类型的nat网关。  */
	NatGatewaySpec     string `json:"natGatewaySpec,omitempty"` /*  enableSnat=true必填，nat网关规格：small，medium，large，xlarge，您可查看不同<a href="https://www.ctyun.cn/document/10026759/10043996">类型规格</a>的详细说明  */
	ElbProdCode        string `json:"elbProdCode,omitempty"`    /*  ApiServer的ELB类型，standardI（标准I型） ,standardII（标准II型）, enhancedI（增强I型）, enhancedII（增强II型） , higherI（高阶I型）
	您可查看不同<a href="https://www.ctyun.cn/document/10026756/10032048">类型规格</a>的详细说明
	*/
	NodeLabels *CcseCreateClusterClusterBaseInfoNodeLabelsRequest `json:"nodeLabels"`        /*  节点标签  */
	PodCidr    string                                             `json:"podCidr,omitempty"` /*  pod网络cidr，使用cubecni作为网络插件时，podCidr传值为vpc cidr，vpc cidr通过如下方式查询：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8658&data=94&isNormal=1&vid=88">查询VPC列表</a>
	注：网络插件为calico时，podCidr与vpcCidr与serviceCidr不能重叠。选择cubecni时，podCidr（vpcCidr）与serviceCidr不能重叠。  */
	InstallNginxIngress *bool  `json:"installNginxIngress"`          /*  是否安装nginx ingress插件  */
	NginxIngressLBSpec  string `json:"nginxIngressLBSpec,omitempty"` /*  installNginxIngress=true必填，nginx ingress ELB的规格:standardI（标准I型） ,standardII（标准II型）, enhancedI（增强I型）, enhancedII（增强II型） , higherI（高阶I型）
	您可查看不同<a href="https://www.ctyun.cn/document/10026756/10032048">规格的详细说明</a>  */
	NginxIngressLBNetWork string                                            `json:"nginxIngressLBNetWork,omitempty"` /*  installNginxIngress=true必填，nginx ingress访问方式：external（公网），internal（内网），当选择公网时将自动创建eip额外产生eip相关费用  */
	BillMode              string                                            `json:"billMode,omitempty"`              /*  计费模式：1为包周期，2为按需  */
	CycleType             string                                            `json:"cycleType,omitempty"`             /*  订购周期类型：3表示按月订购，按需订购不传  */
	CycleCnt              int32                                             `json:"cycleCnt,omitempty"`              /*  订购时长：cycleType为3时，cycleCnt为1表示订购1个月  */
	AutoRenewStatus       *bool                                             `json:"autoRenewStatus"`                 /*  是否自动续费，默认false  */
	AutoRenewCycleType    string                                            `json:"autoRenewCycleType,omitempty"`    /*  自动续费周期类型，3（按月），5（按年）  */
	AutoRenewCycleCount   string                                            `json:"autoRenewCycleCount,omitempty"`   /*  自动续期续期时长  */
	ContainerRuntime      string                                            `json:"containerRuntime,omitempty"`      /*  容器运行时,可选containerd、docker，您可查看<a href="https://www.ctyun.cn/document/10083472/10902208">容器运行时说明</a>  */
	Timezone              string                                            `json:"timezone,omitempty"`              /*  时区  */
	ClusterVersion        string                                            `json:"clusterVersion,omitempty"`        /*  集群版本1.23.3 ，1.25.6 ，1.27.8，1.29.3，您可查看<a href="https://www.ctyun.cn/document/10083472/10650447">集群版本说明</a>选择  */
	DeployType            string                                            `json:"deployType,omitempty"`            /*  部署模式：单可用区为single，多可用区为multi。  */
	AzInfos               []*CcseCreateClusterClusterBaseInfoAzInfosRequest `json:"azInfos"`                         /*  可用区信息，包括可用区编码，该可用区下master节点数量。专有版时必填  */
	ServiceCidr           string                                            `json:"serviceCidr,omitempty"`           /*  服务cidr。网络插件为calico时，podCidr与vpcCidr与serviceCidr不能重叠。选择cubecni时，podCidr（vpcCidr）与serviceCidr不能重叠。  */
	VpcUuid               string                                            `json:"vpcUuid,omitempty"`               /*  虚拟私有云ID，通过以下方式查询：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94&vid=88">查询用户vpc列表</a>
	注：在多可用区类型资源池下，vpcID通常为“vpc-”开头，非多可用区类型资源池vpcID为uuid格式  */
	ClusterName              string `json:"clusterName,omitempty"`      /*  集群名字  */
	KubeProxy                string `json:"kubeProxy,omitempty"`        /*  kubeProxy类型：iptables或ipvs，您可查看<a href="https://www.ctyun.cn/document/10083472/10915725">iptables与IPVS如何选择</a>  */
	PluginCstorcsiAk         string `json:"pluginCstorcsiAk,omitempty"` /*  CSI插件AK  */
	PluginCstorcsiSk         string `json:"pluginCstorcsiSk,omitempty"` /*  CSI插件SK  */
	PluginCstorcsiEnabled    *bool  `json:"pluginCstorcsiEnabled"`      /*  是否安装csi插件  */
	PluginCcseMonitorEnabled *bool  `json:"pluginCcseMonitorEnabled"`   /*  是否安装监控插件  */
	ClusterSeries            string `json:"clusterSeries,omitempty"`    /*  集群系列，cce.standard-专有版，cce.managed-托管版，您可查看<a href="https://www.ctyun.cn/document/10083472/10892150">产品定义</a>选择  */
	ProjectId                string `json:"projectId,omitempty"`        /*  集群所属企业项目id，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目
	注：默认值为"0"  */
	EnableAls                 *bool  `json:"enableAls"`                   /*  是否开启日志插件  */
	EnableAlsCubeEventer      *bool  `json:"enableAlsCubeEventer"`        /*  是否开启事件采集插件  */
	PluginNginxIngressEnabled *bool  `json:"pluginNginxIngressEnabled"`   /*  是否开启nginxingress  */
	CustomScriptBase64        *bool  `json:"customScriptBase64"`          /*  节点自定义脚本是否采用base64编码，默认为false  */
	HostScript                string `json:"hostScript,omitempty"`        /*  节点自定义脚本，如果customScriptBase64=false，此参数传明文；若customScriptBase64=true，此参数传base64后内容  */
	EnablePostUserScript      *bool  `json:"enablePostUserScript"`        /*  节点部署后置脚本是否采用base64编码  */
	PostUserScript            string `json:"postUserScript,omitempty"`    /*  节点部署后置脚本，节点自定义脚本，如果enablePostUserScript=false，此参数传明文；若enablePostUserScript=true，此参数传base64后内容  */
	EnableHostName            *bool  `json:"enableHostName"`              /*  自定义节点名称是否开启  */
	HostNamePrefix            string `json:"hostNamePrefix,omitempty"`    /*  自定义节点名称前缀  */
	HostNamePostfix           string `json:"hostNamePostfix,omitempty"`   /*  自定义节点名称后缀  */
	NodeTaints                string `json:"nodeTaints,omitempty"`        /*  节点污点，格式为 [{\"key\":\"{key}\",\"value\":\"{value}\",\"effect\":\"{effect}\"}]，上述的{key}、{value}、{effect}替换为所需字段。effect枚举包括NoSchedule、PreferNoSchedule、NoExecute  */
	NodeUnschedulable         *bool  `json:"nodeUnschedulable"`           /*  节点不可调度  */
	ClusterDesc               string `json:"clusterDesc,omitempty"`       /*  集群描述  */
	DeleteProtection          *bool  `json:"deleteProtection"`            /*  集群删除保护  */
	SeriesType                string `json:"seriesType,omitempty"`        /*  托管版集群规格，托管版集群必填，单实例-managedbase，多实例-managedpro。单/多实例指控制面是否高可用，生产环境建议使用多实例  */
	Ipvlan                    *bool  `json:"ipvlan"`                      /*  基于IPVLAN做弹性网卡共享，主机镜像只有使用CtyunOS系统才能生效  */
	NetworkPolicy             *bool  `json:"networkPolicy"`               /*  提供基于策略的网络访问控制  */
	IpStackType               string `json:"ipStackType,omitempty"`       /*  ip栈类型，ipv4或dual。默认ipv4  */
	ServiceCidrV6             string `json:"serviceCidrV6,omitempty"`     /*  ipStackType=dual时必填，IPV6网段的掩码必须在112到120之间，如fc00::/112  */
	ClusterAlias              string `json:"clusterAlias,omitempty"`      /*  集群显示名称  */
	NodeScale                 string `json:"nodeScale,omitempty"`         /*  托管版集群节点规模。当seriesType=managedbase时，选择节点规模，可选10；当seriesType=managedpro时，选择节点规模，可选为50，200，1000，2000  */
	SyncNodeLabels            *bool  `json:"syncNodeLabels"`              /*  节点池节点标签的改动将会被同步更新到存量节点，默认false  */
	SyncNodeTaints            *bool  `json:"syncNodeTaints"`              /*  节点池节点污点的改动将会被同步更新到存量节点，默认false  */
	EnableAffinityGroup       *bool  `json:"enableAffinityGroup"`         /*  是否默认使用反亲和性的云主机组创建节点池，默认false  */
	AffinityGroupUuid         string `json:"affinityGroupUuid,omitempty"` /*  是否使用已有反亲和性的云主机组创建节点池，enableAffinityGroup=true时填写，不填系统使用默认配置创建
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8324&data=87&vid=81">查询云主机组列表或者详情</a>  */
	DelegateName           string                                                 `json:"delegateName,omitempty"` /*  工作节点使用已有委托，不填使用系统默认创建委托  */
	ResourceLabels         *CcseCreateClusterClusterBaseInfoResourceLabelsRequest `json:"resourceLabels"`         /*  云主机资源标签  */
	CpuManagerPolicyEnable *bool                                                  `json:"cpuManagerPolicyEnable"` /*  是否开启cpu管理策略，默认false  */
	San                    string                                                 `json:"san,omitempty"`          /*  自定义证书SAN，您可查看<a href="https://www.ctyun.cn/document/10083472/10915719">自定义集群APIServer证书SAN</a>  */
	CustomCAEnable         *bool                                                  `json:"customCAEnable"`         /*  是否开启集群CA，默认false  */
	CustomCA               string                                                 `json:"customCA,omitempty"`     /*  集群CA  */
}

type CcseCreateClusterMasterHostRequest struct {
	Cpu         int32  `json:"cpu,omitempty"`         /*  cpu  */
	Mem         int32  `json:"mem,omitempty"`         /*  内存  */
	ItemDefName string `json:"itemDefName,omitempty"` /*  实例规格名称
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87&isNormal=1&vid=81">查询主机规格资源</a>  */
	ItemDefType string                                         `json:"itemDefType,omitempty"` /*  实例规格类型  */
	Size        int32                                          `json:"size,omitempty"`        /*  master节点数量  */
	SysDisk     *CcseCreateClusterMasterHostSysDiskRequest     `json:"sysDisk"`               /*  系统盘信息  */
	DataDisks   []*CcseCreateClusterMasterHostDataDisksRequest `json:"dataDisks"`             /*  数据盘信息  */
}

type CcseCreateClusterSlaveHostRequest struct {
	Cpu         int32  `json:"cpu,omitempty"`         /*  cpu  */
	Mem         int32  `json:"mem,omitempty"`         /*  内存  */
	ItemDefName string `json:"itemDefName,omitempty"` /*  实例规格名称
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87&isNormal=1&vid=81">查询主机规格资源</a>  */
	ItemDefType     string                                        `json:"itemDefType,omitempty"`     /*  实例规格族  */
	Size            int32                                         `json:"size,omitempty"`            /*  worker节点数量  */
	SysDisk         *CcseCreateClusterSlaveHostSysDiskRequest     `json:"sysDisk"`                   /*  系统盘信息  */
	DataDisks       []*CcseCreateClusterSlaveHostDataDisksRequest `json:"dataDisks"`                 /*  数据盘信息  */
	ForeignMirrorId string                                        `json:"foreignMirrorId,omitempty"` /*  镜像id，worker节点为ecs类型时必填  */
	MirrorType      int32                                         `json:"mirrorType"`                /*  镜像类型，0-私有，1-公有。  */
	MirrorName      string                                        `json:"mirrorName,omitempty"`      /*  镜像名称，worker节点为ebm类型必填。您可查看<a href="https://www.ctyun.cn/document/10083472/11004475">节点规格和节点镜像</a>  */
	AzInfos         []*CcseCreateClusterSlaveHostAzInfosRequest   `json:"azInfos"`                   /*  可用区信息，包括可用区编码，可用区worker节点数量  */
}

type CcseCreateClusterClusterBaseInfoNodeLabelsRequest struct{}

type CcseCreateClusterClusterBaseInfoAzInfosRequest struct {
	AzName string `json:"azName,omitempty"` /*  可用区编码查询：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&vid=81">资源池可用区查询</a>  */
	Size int32 `json:"size,omitempty"` /*  该可用区节点数量  */
}

type CcseCreateClusterClusterBaseInfoResourceLabelsRequest struct{}

type CcseCreateClusterMasterHostSysDiskRequest struct {
	ItemDefName string `json:"itemDefName,omitempty"` /*  系统盘规格，云硬盘类型，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	Size int32 `json:"size,omitempty"` /*  系统盘大小，单位为G  */
}

type CcseCreateClusterMasterHostDataDisksRequest struct {
	ItemDefName string `json:"itemDefName,omitempty"` /*  数据盘规格名称，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	Size      int32  `json:"size,omitempty"`      /*  数据盘大小，单位为G  */
	DecTypeId string `json:"decTypeId,omitempty"` /*  专属存储池id，需用专属资源池权限
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=4&api=12772&data=113&isNormal=1&vid=106">查询存储池信息详情</a>  */
}

type CcseCreateClusterSlaveHostSysDiskRequest struct {
	ItemDefName string `json:"itemDefName,omitempty"` /*  系统盘规格，云硬盘类型，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	Size int32 `json:"size,omitempty"` /*  系统盘大小，单位为G  */
}

type CcseCreateClusterSlaveHostDataDisksRequest struct {
	ItemDefName string `json:"itemDefName,omitempty"` /*  数据盘规格名称，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	Size      int32  `json:"size,omitempty"`      /*  数据盘大小，单位为G  */
	DecTypeId string `json:"decTypeId,omitempty"` /*  专属存储池id，需用专属资源池权限
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=4&api=12772&data=113&isNormal=1&vid=106">查询存储池信息详情</a>  */
}

type CcseCreateClusterSlaveHostAzInfosRequest struct {
	AzName string `json:"azName,omitempty"` /*  可用区编码查询：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&vid=81">资源池可用区查询</a>  */
	Size int32 `json:"size,omitempty"` /*  该可用区节点数量  */
}

type CcseCreateClusterResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                              `json:"message,omitempty"`    /*  提示信息  */
	Error      string                              `json:"error,omitempty"`      /*  错误码  */
	ReturnObj  *CcseCreateClusterReturnObjResponse `json:"returnObj"`            /*  返回对象  */
}

type CcseCreateClusterReturnObjResponse struct {
	OrderId string `json:"orderId,omitempty"` /*  订单id  */
	OrderNo string `json:"orderNo,omitempty"` /*  订单编号  */
}
