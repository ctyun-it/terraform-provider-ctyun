package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateSameConfigInstancesApi
/* 支持创建一台或多台与已有云主机相同配置的云主机<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />&emsp;&emsp;计费模式：确认开通云主机的计费模式，详细查看<a href="https://www.ctyun.cn/document/10026730/10030877">计费模式</a><br />&emsp;&emsp;产品选型：购买弹性云主机前，请先阅读<a href="https://www.ctyun.cn/document/10026730/10118193">规格说明</a>了解弹性云主机的选型基本信息，并通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87">查询一个或多个云主机规格资源</a>接口，获取当前资源池可用云主机规格信息<br />&emsp;&emsp;网络规划：规划云主机的网络环境，详细查看<a href="https://www.ctyun.cn/document/10026730/10237410">弹性云主机-网络</a><br /><b>注意事项：</b><br />&emsp;&emsp;成本估算：了解云主机的<a href="https://www.ctyun.cn/document/10026730/10028700">计费项</a><br />&emsp;&emsp;用户配额：确认个人在不同资源池下资源配额，可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9714&data=87">用户配额查询</a>接口进行查询<br />&emsp;&emsp;异步接口：该接口为异步接口，下单过后会拿到主订单ID（masterOrderID），后续可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9607&data=87&isNormal=1">根据masterOrderID查询云主机ID</a>，使用主订单ID来对订单情况与开通成功后的资源ID进行查询<br />&emsp;&emsp;企业项目：保证资源隔离，需确保所有资源与云主机处于相同企业项目下才可进行挂载操作，如网络资源：VPC、安全组、弹性IP、子网<br />&emsp;&emsp;标签绑定：云主机绑定标签时，标签键和值存在的情况下，绑定对应标签；不存在的情况下，创建新的标签并绑定云主机。主机创建完成后，云主机变为运行状态，此时标签仍可能未绑定，需等待一段时间（0-10分钟）<br />&emsp;&emsp;监控安装：在云服务器创建成功后，3-5分钟内将完成详细监控Agent安装，即开启云服务器CPU，内存，网络，磁盘，进程等指标详细监控，若不开启，则无任何监控数据<br />&emsp;&emsp;代金券：只支持预付费用户抵扣包周期资源的金额，且不可超过包周期资源的金额<br />&emsp;&emsp;相同配置：体现在主机配置方面与原云主机一致，包括付费方式、规格、镜像、磁盘（包含是否加密和加密密钥）、虚拟私有云、安全组、网卡（网卡和网卡数量，但内网IP地址自动分配）、云主机组、登录方式（若采用密码方式登录，则需要重新配置密码；若采用密钥方式登录，则使用原密钥，若原密钥失效，则需要重新选择密钥）、企业项目、标签，但所有的字段均可修改。在弹性公网IP（默认不使用）、购买时长、购买数量、是否开启监控、用户自定义数据方面需要重新设置
 */type CtecsCreateSameConfigInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateSameConfigInstancesApi(client *core.CtyunClient) *CtecsCreateSameConfigInstancesApi {
	return &CtecsCreateSameConfigInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/create-same-config-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateSameConfigInstancesApi) Do(ctx context.Context, credential core.Credential, req *CtecsCreateSameConfigInstancesRequest) (*CtecsCreateSameConfigInstancesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsCreateSameConfigInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateSameConfigInstancesRequest struct {
	ClientToken     string                                                  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时  */
	RegionID        string                                                  `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>   */
	AzName          string                                                  `json:"azName,omitempty"`          /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区,无需填写本字段）  */
	InstanceID      string                                                  `json:"instanceID,omitempty"`      /*  已有的云主机ID，获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>       */
	InstanceName    string                                                  `json:"instanceName,omitempty"`    /*  云主机名称。不同操作系统下，云主机名称规则有差异。<br />Windows：长度为2-15个字符，允许使用大小写字母、数字或连字符（-）。不能以连字符（-）开头或结尾，不能连续使用连字符（-），也不能仅使用数字；<br />其他操作系统：长度为2-64字符，允许使用点（.）分隔字符成多段，每段允许使用大小写字母、数字或连字符（-），但不能连续使用点号（.）或连字符（-），不能以点号（.）或连字符（-）开头或结尾，也不能仅使用数字   */
	DisplayName     string                                                  `json:"displayName,omitempty"`     /*  云主机显示名称，长度为2-63字符  */
	FlavorID        string                                                  `json:"flavorID,omitempty"`        /*  云主机规格ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10118193">规格说明</a>了解弹性云主机的选型基本信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87">查询一个或多个云主机规格资源</a><br />注：同一规格名称在不同资源池不同可用区的规格ID是不同的，调用前需确认规格ID是否归属当前资源池，多可用区资源池确认是否归属当前可用区  */
	ImageID         string                                                  `json:"imageID,omitempty"`         /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10030151">镜像概述</a>来了解云主机镜像<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89">查询可以使用的镜像资源</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4765&data=89">创建私有镜像（云主机系统盘）</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5230&data=89">创建私有镜像（云主机数据盘）</a><br />注：同一镜像名称在不同资源池的镜像ID是不同的，调用前需确认镜像ID是否归属当前资源池  */
	BootDiskType    string                                                  `json:"bootDiskType,omitempty"`    /*  系统盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO，<br />FAST-SSD：极速型SSD<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>来了解磁盘类型及其对应性能指标   */
	BootDiskSize    int32                                                   `json:"bootDiskSize,omitempty"`    /*  系统盘大小单位为GiB，取值范围：[40, 2048]，注：创建云主机过程中会存在单位转换，因此该参数只能传入整型，如果填写小数值则会被取整，影响到涉及计费   */
	VpcID           string                                                  `json:"vpcID,omitempty"`           /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94">查询VPC列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94">创建VPC</a><br />注：在多可用区类型资源池下，vpcID通常为“vpc-”开头，非多可用区类型资源池vpcID为uuid格式     */
	OnDemand        *bool                                                   `json:"onDemand"`                  /*  购买方式，取值范围：<br />false：按周期，<br />true：按需<br />您可以查看<a href="https://www.ctyun.cn/document/10026730/10030877">计费模式</a>了解云主机的计费模式<br />注：按周期（false）创建云主机需要同时指定cycleCount和cycleType参数    */
	NetworkCardList []*CtecsCreateSameConfigInstancesNetworkCardListRequest `json:"networkCardList"`           /*  网卡信息列表，您可以查看<a href="https://www.ctyun.cn/document/10026730/10225195">弹性网卡概述</a>了解弹性网卡相关信息  */
	ExtIP           string                                                  `json:"extIP,omitempty"`           /*  是否使用弹性公网IP，取值范围：<br />0：不使用，<br />1：自动分配，<br />2：使用已有<br />注：自动分配弹性公网，默认分配IPv4弹性公网，需填写带宽大小，如需ipv6请填写弹性IP版本（即参数extIP="1"时，需填写参数bandwidth、ipVersion，ipVersion含默认值ipv4）；<br />使用已有弹性公网，请填写弹性公网IP的ID，默认为ipv4版本，如使用已有ipv6，请填写弹性ip版本（即参数extIP="2"时，需填写eipID或ipv6AddressID，同时ipv6情况下请填写ipVersion）   */
	ProjectID       string                                                  `json:"projectID,omitempty"`       /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br />  */
	SecGroupList    []string                                                `json:"secGroupList"`              /*  安全组ID列表，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a><br />注：在多可用区类型资源池下，安全组ID通常以“sg-”开头，非多可用区类型资源池安全组ID为uuid格式；  */
	DataDiskList    []*CtecsCreateSameConfigInstancesDataDiskListRequest    `json:"dataDiskList"`              /*  数据盘信息列表，注：同一云主机下最多可挂载8块数据盘  */
	IpVersion       string                                                  `json:"ipVersion,omitempty"`       /*  弹性IP版本，取值范围：<br />ipv4：v4地址，<br />ipv6：v6地址<br />不指定默认为ipv4。注：请先确认该资源池是否支持ipv6（多可用区类资源池暂不支持）  */
	Bandwidth       int32                                                   `json:"bandwidth,omitempty"`       /*  带宽大小，单位为Mbit/s，取值范围：[1, 2000]   */
	Ipv6AddressID   string                                                  `json:"ipv6AddressID,omitempty"`   /*  弹性公网IPv6的ID；填写该参数时请填写ipVersion为ipv6  */
	EipID           string                                                  `json:"eipID,omitempty"`           /*  弹性公网IP的ID，您可以查看<a href="https://www.ctyun.cn/document/10026753/10026909">产品定义-弹性IP</a>来了解弹性公网IP <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8652&data=94&isNormal=1&vid=88">查询指定地域已创建的弹性 IP</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5723&data=94&vid=88">创建弹性 IP</a>  */
	AffinityGroupID string                                                  `json:"affinityGroupID,omitempty"` /*  云主机组ID，获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8324&data=87">查询云主机组列表或者详情</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8316&data=87"> 创建云主机组</a>   */
	KeyPairID       string                                                  `json:"keyPairID,omitempty"`       /*  密钥对ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10230540">密钥对</a>来了解密钥对相关内容 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87">查询一个或多个密钥对</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8344&data=87">创建一对SSH密钥对</a>  */
	UserPassword    string                                                  `json:"userPassword,omitempty"`    /*  用户密码，满足以下规则：<br />长度在8～30个字符；<br />必须包含大写字母、小写字母、数字以及特殊符号中的三项；<br />特殊符号可选：()`-!@#$%^&*_-+=｜{}[]:;'<>,.?/且不能以斜线号 / 开头；<br />不能包含3个及以上连续字符；<br />Linux镜像不能包含镜像用户名（root）、用户名的倒序（toor）、用户名大小写变化（如RoOt、rOot等）；<br />Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator等）  */
	CycleCount      int32                                                   `json:"cycleCount,omitempty"`      /*  订购时长，该参数需要与cycleType一同使用<br />注：最长订购周期为60个月（5年）；cycleType与cycleCount一起填写；按量付费（即onDemand为true）时，无需填写该参数（填写无效）  */
	CycleType       string                                                  `json:"cycleType,omitempty"`       /*  订购周期类型，取值范围：<br />MONTH：按月，<br />YEAR：按年<br />注：cycleType与cycleCount一起填写；按量付费（即onDemand为true）时，无需填写该参数（填写无效）  */
	AutoRenewStatus int32                                                   `json:"autoRenewStatus,omitempty"` /*  是否自动续订，取值范围：<br />0：不续费，<br />1：自动续费<br />注：按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年      */
	UserData        string                                                  `json:"userData,omitempty"`        /*  用户自定义数据，需要以Base64方式编码，Base64编码后的长度限制为1-16384字符。  */
	OrderCount      int32                                                   `json:"orderCount,omitempty"`      /*  购买数量，取值范围：[1, 50]。注：不填则默认为1   */
	PayVoucherPrice float32                                                 `json:"payVoucherPrice"`           /*  代金券，满足以下规则：<br />两位小数，不足两位自动补0，超过两位小数无效；<br />不可为负数；<br />注：字段为0时表示不使用代金券，默认不使用代金券  */
	LabelList       []*CtecsCreateSameConfigInstancesLabelListRequest       `json:"labelList"`                 /*  标签信息列表，注：单台云主机最多可绑定10个标签；主机创建完成后，云主机变为运行状态，此时标签仍可能未绑定，需等待一段时间（0-10分钟）  */
	GpuDriverKits   string                                                  `json:"gpuDriverKits,omitempty"`   /*  GPU云主机安装驱动的工具包，仅在同时选择NVIDIA显卡、计算加速型、linux公共镜像三个条件下，支持安装驱动  */
	MonitorService  *bool                                                   `json:"monitorService"`            /*  监控参数，支持通过该参数指定云主机在创建后是否开启详细监控，取值范围： <br />false：不开启，<br />true：开启<br />若指定该参数为true或不指定该参数，云主机内默认开启最新详细监控服务<br />若指定该参数为false，默认公共镜像不开启最新监控服务；私有镜像使用镜像中保留的监控服务<br />说明：仅部分资源池支持monitorService参数，详细请参考<a href="https://www.ctyun.cn/document/10026730/10325957">监控Agent概览</a>  */
}

type CtecsCreateSameConfigInstancesNetworkCardListRequest struct {
	NicName  string `json:"nicName,omitempty"`  /*  该参数暂时无法使用。网卡名称，长度2~32，支持拉丁字母、中文、数字、下划线、连字符，中文或英文字母开头，不能以http:或https:开头  */
	FixedIP  string `json:"fixedIP,omitempty"`  /*  内网IPv4地址，注：不可使用已占用IP  */
	IsMaster *bool  `json:"isMaster"`           /*  是否主网卡，取值范围：<br />true：表示主网卡，<br />false：表示扩展网卡<br />注：只能含有一个主网卡  */
	SubnetID string `json:"subnetID,omitempty"` /*  子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94">查询子网列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4812&data=94">创建子网</a><br />注：在多可用区类型资源池下，subnetID通常以“subnet-”开头，非多可用区类型资源池subnetID为uuid格式   */
}

type CtecsCreateSameConfigInstancesDataDiskListRequest struct {
	DiskMode  string `json:"diskMode,omitempty"` /*  云硬盘属性，取值范围：<br />FCSAN：光纤通道协议的SAN网络，<br />ISCSI：小型计算机系统接口，<br />VBD：虚拟块存储设备<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162960">磁盘模式及使用方法</a>  */
	DiskType  string `json:"diskType,omitempty"` /*  云硬盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO，<br />FAST-SSD：极速型SSD<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	DiskSize  int32  `json:"diskSize,omitempty"` /*  磁盘容量，单位为GB，取值范围：[10, 32768]   */
	IsEncrypt *bool  `json:"isEncrypt"`          /*  磁盘是否加密，取值范围：<br />true：加密，<br />false：不加密<br />若该参数为true且不填写cmkID（加密密钥ID），则生成默认密钥（不在密钥列表中）加密进行加密  */
	CmkID     string `json:"cmkID,omitempty"`    /*  加密密钥ID  */
}

type CtecsCreateSameConfigInstancesLabelListRequest struct {
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签键，长度限制1-32字符，注：同一台云主机绑定多个标签时，标签键不可重复  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签值，长度限制1-32字符   */
}

type CtecsCreateSameConfigInstancesResponse struct {
	StatusCode  int32                                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                           `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                           `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                           `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                           `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsCreateSameConfigInstancesReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsCreateSameConfigInstancesReturnObjResponse struct {
	MasterOrderID    string `json:"masterOrderID,omitempty"`    /*  主订单ID。调用方在拿到masterOrderID之后，可以使用masterOrderID进一步确认订单状态及资源状态<br />查询订单状态及资源UUID：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9607&data=87&isNormal=1">根据masterOrderID查询云主机ID</a>  */
	MasterOrderNO    string `json:"masterOrderNO,omitempty"`    /*  订单号  */
	MasterResourceID string `json:"masterResourceID,omitempty"` /*  主资源ID  */
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID  */
}
