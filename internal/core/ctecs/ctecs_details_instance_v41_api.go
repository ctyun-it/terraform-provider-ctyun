package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetailsInstanceV41Api
/* 该接口提供用户一台云主机信息查询功能，用户可以根据此接口的返回值了解自己云主机的详细信息<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;单台查询：当前接口只能查询单台云主机信息，查询多台云主机信息请使用接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>进行查询<br />
 */type CtecsDetailsInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetailsInstanceV41Api(client *core.CtyunClient) *CtecsDetailsInstanceV41Api {
	return &CtecsDetailsInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/instance-details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetailsInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetailsInstanceV41Request) (*CtecsDetailsInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDetailsInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetailsInstanceV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}

type CtecsDetailsInstanceV41Response struct {
	StatusCode  int32                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                    `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                    `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                    `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDetailsInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDetailsInstanceV41ReturnObjResponse struct {
	ProjectID          string                                                     `json:"projectID,omitempty"`      /*  企业项目ID  */
	AzName             string                                                     `json:"azName,omitempty"`         /*  可用区名称  */
	AttachedVolume     []string                                                   `json:"attachedVolume"`           /*  云硬盘ID列表  */
	Addresses          []*CtecsDetailsInstanceV41ReturnObjAddressesResponse       `json:"addresses"`                /*  网络地址信息  */
	ResourceID         string                                                     `json:"resourceID,omitempty"`     /*  云主机资源ID  */
	InstanceID         string                                                     `json:"instanceID,omitempty"`     /*  云主机ID  */
	DisplayName        string                                                     `json:"displayName,omitempty"`    /*  云主机显示名称  */
	InstanceName       string                                                     `json:"instanceName,omitempty"`   /*  云主机名称  */
	OsType             int32                                                      `json:"osType,omitempty"`         /*  操作系统类型，详见枚举值表  */
	InstanceStatus     string                                                     `json:"instanceStatus,omitempty"` /*  云主机状态，，请通过状态<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	ExpiredTime        string                                                     `json:"expiredTime,omitempty"`    /*  到期时间  */
	AvailableDay       int32                                                      `json:"availableDay,omitempty"`   /*  可用天数  */
	UpdatedTime        string                                                     `json:"updatedTime,omitempty"`    /*  更新时间  */
	CreatedTime        string                                                     `json:"createdTime,omitempty"`    /*  创建时间  */
	ZabbixName         string                                                     `json:"zabbixName,omitempty"`     /*  监控对象名称  */
	SecGroupList       []*CtecsDetailsInstanceV41ReturnObjSecGroupListResponse    `json:"secGroupList"`             /*  安全组信息  */
	PrivateIP          string                                                     `json:"privateIP,omitempty"`      /*  内网IPv4地址  */
	PrivateIPv6        string                                                     `json:"privateIPv6,omitempty"`    /*  内网IPv6地址  */
	NetworkCardList    []*CtecsDetailsInstanceV41ReturnObjNetworkCardListResponse `json:"networkCardList"`          /*  网卡信息  */
	VipInfoList        []*CtecsDetailsInstanceV41ReturnObjVipInfoListResponse     `json:"vipInfoList"`              /*  虚拟IP信息列表  */
	VipCount           int32                                                      `json:"vipCount,omitempty"`       /*  vip数目  */
	AffinityGroup      *CtecsDetailsInstanceV41ReturnObjAffinityGroupResponse     `json:"affinityGroup"`            /*  云主机组信息  */
	Image              *CtecsDetailsInstanceV41ReturnObjImageResponse             `json:"image"`                    /*  镜像信息  */
	Flavor             *CtecsDetailsInstanceV41ReturnObjFlavorResponse            `json:"flavor"`                   /*  云主机规格信息  */
	OnDemand           *bool                                                      `json:"onDemand"`                 /*  付费方式，取值范围：<br />true：表示按量付费 <br />false：表示包周期  */
	VpcName            string                                                     `json:"vpcName,omitempty"`        /*  虚拟私有云名称  */
	VpcID              string                                                     `json:"vpcID,omitempty"`          /*  虚拟私有云ID  */
	FixedIPList        []string                                                   `json:"fixedIPList"`              /*  内网IP列表  */
	FloatingIP         string                                                     `json:"floatingIP,omitempty"`     /*  公网IP  */
	SubnetIDList       []string                                                   `json:"subnetIDList"`             /*  子网ID列表  */
	KeypairName        string                                                     `json:"keypairName,omitempty"`    /*  密钥对名称  */
	DeletionProtection *bool                                                      `json:"deletionProtection"`       /*  是否开启实例删除保护  */
	DelegateName       string                                                     `json:"delegateName,omitempty"`   /*  委托名称，注：委托绑定目前仅支持多可用区类型资源池，非可用区资源池为空字符串  */
	RemainingDay       int32                                                      `json:"remainingDay,omitempty"`   /*  距离释放剩余天数  */
	ReleaseTime        string                                                     `json:"releaseTime,omitempty"`    /*  释放时间  */
	Metadata           *CtecsDetailsInstanceV41ReturnObjMetadataResponse          `json:"metadata"`                 /*  用户自定义元数据，注：仅多可用区类型资源池支持返回，非可用区资源池请调用查询元数据接口 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8319&data=87&isNormal=1">云主机元数据查询</a>   */
}

type CtecsDetailsInstanceV41ReturnObjAddressesResponse struct {
	VpcName     string                                                          `json:"vpcName,omitempty"` /*  vpc名称  */
	AddressList []*CtecsDetailsInstanceV41ReturnObjAddressesAddressListResponse `json:"addressList"`       /*  网络地址列表  */
}

type CtecsDetailsInstanceV41ReturnObjSecGroupListResponse struct {
	SecurityGroupID   string `json:"securityGroupID,omitempty"`   /*  安全组ID  */
	SecurityGroupName string `json:"securityGroupName,omitempty"` /*  安全组名称  */
}

type CtecsDetailsInstanceV41ReturnObjNetworkCardListResponse struct {
	IPv4Address   string   `json:"IPv4Address,omitempty"`   /*  IPv4地址  */
	IPv6Address   []string `json:"IPv6Address"`             /*  IPv6地址列表  */
	SubnetID      string   `json:"subnetID,omitempty"`      /*  子网ID  */
	SubnetCidr    string   `json:"subnetCidr,omitempty"`    /*  子网网段信息  */
	IsMaster      *bool    `json:"isMaster"`                /*  是否主网卡，取值范围：<br />true：主网卡，<br />false：扩展网卡  */
	Gateway       string   `json:"gateway,omitempty"`       /*  网关地址  */
	NetworkCardID string   `json:"networkCardID,omitempty"` /*  网卡ID  */
	SecurityGroup []string `json:"securityGroup"`           /*  安全组ID列表  */
}

type CtecsDetailsInstanceV41ReturnObjVipInfoListResponse struct {
	VipID          string `json:"vipID,omitempty"`          /*  虚拟IP的ID  */
	VipAddress     string `json:"vipAddress,omitempty"`     /*  虚拟IP地址  */
	VipBindNicIP   string `json:"vipBindNicIP,omitempty"`   /*  虚拟IP绑定的网卡对应IPv4地址  */
	VipBindNicIPv6 string `json:"vipBindNicIPv6,omitempty"` /*  虚拟IP绑定的网卡对应IPv6地址  */
	NicID          string `json:"nicID,omitempty"`          /*  网卡ID  */
}

type CtecsDetailsInstanceV41ReturnObjAffinityGroupResponse struct {
	Policy            string `json:"policy,omitempty"`            /*  云主机组策略  */
	AffinityGroupName string `json:"affinityGroupName,omitempty"` /*  云主机组名称  */
	AffinityGroupID   string `json:"affinityGroupID,omitempty"`   /*  云主机组ID  */
}

type CtecsDetailsInstanceV41ReturnObjImageResponse struct {
	ImageID   string `json:"imageID,omitempty"`   /*  镜像ID  */
	ImageName string `json:"imageName,omitempty"` /*  镜像名称  */
}

type CtecsDetailsInstanceV41ReturnObjFlavorResponse struct {
	FlavorID     string `json:"flavorID,omitempty"`     /*  规格ID  */
	FlavorName   string `json:"flavorName,omitempty"`   /*  规格名称  */
	FlavorCPU    int32  `json:"flavorCPU,omitempty"`    /*  VCPU  */
	FlavorRAM    int32  `json:"flavorRAM,omitempty"`    /*  内存  */
	GpuType      string `json:"gpuType,omitempty"`      /*  GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加  */
	GpuCount     int32  `json:"gpuCount,omitempty"`     /*  GPU数目  */
	GpuVendor    string `json:"gpuVendor,omitempty"`    /*  GPU名称  */
	VideoMemSize int32  `json:"videoMemSize,omitempty"` /*  GPU显存大小  */
}

type CtecsDetailsInstanceV41ReturnObjMetadataResponse struct{}

type CtecsDetailsInstanceV41ReturnObjAddressesAddressListResponse struct {
	Addr       string `json:"addr,omitempty"`       /*  IP地址  */
	Version    int32  `json:"version,omitempty"`    /*  IP版本  */
	RawType    string `json:"type,omitempty"`       /*  网络类型，取值范围：<br />fixed：内网，<br />floating：弹性公网  */
	IsMaster   *bool  `json:"isMaster"`             /*  网络地址对应网卡是否为主网卡  */
	MacAddress string `json:"macAddress,omitempty"` /*  网络地址对应网卡的mac地址  */
}
