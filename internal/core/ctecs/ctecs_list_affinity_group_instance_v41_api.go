package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListAffinityGroupInstanceV41Api
/* 查询云主机组内的云主机<br />可以根据用户给定的云主机组，查询云主机组内云主机的详细信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />
 */type CtecsListAffinityGroupInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListAffinityGroupInstanceV41Api(client *core.CtyunClient) *CtecsListAffinityGroupInstanceV41Api {
	return &CtecsListAffinityGroupInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/affinity-group/list-instance",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListAffinityGroupInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListAffinityGroupInstanceV41Request) (*CtecsListAffinityGroupInstanceV41Response, error) {
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
	var resp CtecsListAffinityGroupInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListAffinityGroupInstanceV41Request struct {
	RegionID        string `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AffinityGroupID string `json:"affinityGroupID,omitempty"` /*  云主机组ID，获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8324&data=87">查询云主机组列表或者详情</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8316&data=87">创建云主机组</a><br />   */
	PageNo          int32  `json:"pageNo,omitempty"`          /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize        int32  `json:"pageSize,omitempty"`        /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListAffinityGroupInstanceV41Response struct {
	StatusCode  int32                                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                              `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                              `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                              `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                              `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListAffinityGroupInstanceV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResponse struct {
	CurrentCount int32                                                        `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                        `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                        `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListAffinityGroupInstanceV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsResponse struct {
	ProjectID       string                                                                      `json:"projectID,omitempty"`      /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	AzName          string                                                                      `json:"azName,omitempty"`         /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
	AttachedVolume  []string                                                                    `json:"attachedVolume"`           /*  云主机挂载的云硬盘列表  */
	Addresses       []*CtecsListAffinityGroupInstanceV41ReturnObjResultsAddressesResponse       `json:"addresses"`                /*  网络地址信息  */
	ResourceID      string                                                                      `json:"resourceID,omitempty"`     /*  资源ID，非资源的UUID，该ID为订单的资源ID（创建云主机接口为异步接口，订单先返回一个该资源ID方便后续查找）<br />获取：<br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	InstanceID      string                                                                      `json:"instanceID,omitempty"`     /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	DisplayName     string                                                                      `json:"displayName,omitempty"`    /*  云主机显示名称，长度为2-63字符  */
	InstanceName    string                                                                      `json:"instanceName,omitempty"`   /*  云主机名称，不同操作系统下，云主机名称规则有差异<br />Windows：长度为2~15个字符，允许使用大小写字母、数字或连字符（-）。不能以连字符（-）开头或结尾，不能连续使用连字符（-），也不能仅使用数字；<br />其他操作系统：长度为2-64字符，允许使用点（.）分隔字符成多段，每段允许使用大小写字母、数字或连字符（-），但不能连续使用点号（.）或连字符（-），不能以点号（.）或连字符（-）开头或结尾  */
	OsType          int32                                                                       `json:"osType,omitempty"`         /*  操作系统类型，详见枚举值表格  */
	InstanceStatus  string                                                                      `json:"instanceStatus,omitempty"` /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	ExpiredTime     string                                                                      `json:"expiredTime,omitempty"`    /*  到期时间  */
	AvailableDay    int32                                                                       `json:"availableDay,omitempty"`   /*  可用(天)  */
	UpdatedTime     string                                                                      `json:"updatedTime,omitempty"`    /*  更新时间  */
	CreatedTime     string                                                                      `json:"createdTime,omitempty"`    /*  创建时间  */
	ZabbixName      string                                                                      `json:"zabbixName,omitempty"`     /*  监控对象名称  */
	SecGroupList    []*CtecsListAffinityGroupInstanceV41ReturnObjResultsSecGroupListResponse    `json:"secGroupList"`             /*  安全组ID列表，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a><br />注：在多可用区类型资源池下，安全组ID通常以“sg-”开头，非多可用区类型资源池安全组ID为uuid格式  */
	PrivateIP       string                                                                      `json:"privateIP,omitempty"`      /*  内网IPv4地址  */
	PrivateIPv6     string                                                                      `json:"privateIPv6,omitempty"`    /*  内网IPv6址  */
	NetworkCardList []*CtecsListAffinityGroupInstanceV41ReturnObjResultsNetworkCardListResponse `json:"networkCardList"`          /*  网卡信息列表，您可以查看<a href="https://www.ctyun.cn/document/10026730/10225195">弹性网卡概述</a>了解弹性网卡相关信息  */
	VipInfoList     []*CtecsListAffinityGroupInstanceV41ReturnObjResultsVipInfoListResponse     `json:"vipInfoList"`              /*  虚拟IP信息列表  */
	VipCount        int32                                                                       `json:"vipCount,omitempty"`       /*  vip数目  */
	AffinityGroup   *CtecsListAffinityGroupInstanceV41ReturnObjResultsAffinityGroupResponse     `json:"affinityGroup"`            /*  云主机组信息  */
	Image           *CtecsListAffinityGroupInstanceV41ReturnObjResultsImageResponse             `json:"image"`                    /*  镜像信息  */
	Flavor          *CtecsListAffinityGroupInstanceV41ReturnObjResultsFlavorResponse            `json:"flavor"`                   /*  规格信息  */
	OnDemand        *bool                                                                       `json:"onDemand"`                 /*  购买方式，取值范围：<br />false：按周期，<br />true：按需<br />您可以查看<a href="https://www.ctyun.cn/document/10026730/10030877">计费模式</a>了解云主机的计费模式<br />注：按周期（false）创建云主机需要同时指定cycleCount和cycleType参数  */
	VpcName         string                                                                      `json:"vpcName,omitempty"`        /*  vpc名称  */
	VpcID           string                                                                      `json:"vpcID,omitempty"`          /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94">查询VPC列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94">创建VPC</a><br />注：在多可用区类型资源池下，vpcID通常以“vpc-”开头，非多可用区类型资源池vpcID为uuid格式  */
	FixedIPList     []string                                                                    `json:"fixedIPList"`              /*  内网IP  */
	FloatingIP      string                                                                      `json:"floatingIP,omitempty"`     /*  公网IP  */
	SubnetIDList    []string                                                                    `json:"subnetIDList"`             /*  子网ID列表  */
	KeypairName     string                                                                      `json:"keypairName,omitempty"`    /*  密钥对名称。满足以下规则：只能由数字、字母、-组成，不能以数字和-开头、以-结尾，且长度为2-63字符  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsAddressesResponse struct {
	VpcName     string                                                                           `json:"vpcName,omitempty"` /*  vpc名称  */
	AddressList []*CtecsListAffinityGroupInstanceV41ReturnObjResultsAddressesAddressListResponse `json:"addressList"`       /*  网络地址列表  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsSecGroupListResponse struct {
	SecurityGroupID   string `json:"securityGroupID,omitempty"`   /*  安全组ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	SecurityGroupName string `json:"securityGroupName,omitempty"` /*  安全组名称  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsNetworkCardListResponse struct {
	IPv4Address   string   `json:"IPv4Address,omitempty"`   /*  IPv4地址  */
	IPv6Address   []string `json:"IPv6Address"`             /*  IPv6地址列表  */
	SubnetID      string   `json:"subnetID,omitempty"`      /*  子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-子网</a>来了解子网<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94">查询子网列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4812&data=94">创建子网</a>  */
	SubnetCidr    string   `json:"subnetCidr,omitempty"`    /*  子网网段信息  */
	IsMaster      *bool    `json:"isMaster"`                /*  是否主网卡，取值范围：<br />true：表示主网卡，<br />false：表示扩展网卡<br />注：只能含有一个主网卡  */
	Gateway       string   `json:"gateway,omitempty"`       /*  网关地址  */
	NetworkCardID string   `json:"networkCardID,omitempty"` /*  网卡ID  */
	SecurityGroup []string `json:"securityGroup"`           /*  安全组ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsVipInfoListResponse struct {
	VipID          string `json:"vipID,omitempty"`          /*  虚拟IP的ID  */
	VipAddress     string `json:"vipAddress,omitempty"`     /*  虚拟IP地址  */
	VipBindNicIP   string `json:"vipBindNicIP,omitempty"`   /*  虚拟IP绑定的网卡对应IPv4地址  */
	VipBindNicIPv6 string `json:"vipBindNicIPv6,omitempty"` /*  虚拟IP绑定的网卡对应IPv6地址  */
	NicID          string `json:"nicID,omitempty"`          /*  网卡ID  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsAffinityGroupResponse struct {
	AffinityGroupPolicy string `json:"affinityGroupPolicy,omitempty"` /*  云主机组策略  */
	AffinityGroupName   string `json:"affinityGroupName,omitempty"`   /*  云主机组名称，满足以下规则：长度在2～63个字符，只能由数字、英文字母、中划线-、下划线_、点.组成  */
	AffinityGroupID     string `json:"affinityGroupID,omitempty"`     /*  云主机组ID，获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8324&data=87">查询云主机组列表或者详情</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8316&data=87">创建云主机组</a><br />  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsImageResponse struct {
	ImageID   string `json:"imageID,omitempty"`   /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10030151">镜像概述</a>来了解云主机镜像<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89">查询可以使用的镜像资源</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4765&data=89">创建私有镜像（云主机系统盘）</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5230&data=89">创建私有镜像（云主机数据盘）</a>  */
	ImageName string `json:"imageName,omitempty"` /*  镜像名称  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsFlavorResponse struct {
	FlavorID     string `json:"flavorID,omitempty"`     /*  云主机规格ID，您可以调用<a href="https://www.ctyun.cn/document/10026730/10118193">规格说明</a>了解弹性云主机的选型基本信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87">查询一个或多个云主机规格资源</a><br />  */
	FlavorName   string `json:"flavorName,omitempty"`   /*  规格名称  */
	FlavorCPU    int32  `json:"flavorCPU,omitempty"`    /*  VCPU数量  */
	FlavorRAM    int32  `json:"flavorRAM,omitempty"`    /*  内存大小，单位为GB  */
	GpuType      string `json:"gpuType,omitempty"`      /*  GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加  */
	GpuCount     int32  `json:"gpuCount,omitempty"`     /*  GPU数目  */
	GpuVendor    string `json:"gpuVendor,omitempty"`    /*  GPU名称  */
	VideoMemSize int32  `json:"videoMemSize,omitempty"` /*  显存大小  */
}

type CtecsListAffinityGroupInstanceV41ReturnObjResultsAddressesAddressListResponse struct {
	Addr    string `json:"addr,omitempty"`    /*  IP地址  */
	Version int32  `json:"version,omitempty"` /*  IP版本  */
	RawType string `json:"type,omitempty"`    /*  网络类型，取值范围：<br />fixed：内网，<br />floating：公网  */
}
