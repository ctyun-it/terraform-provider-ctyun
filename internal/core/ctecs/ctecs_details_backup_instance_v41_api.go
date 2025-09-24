package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetailsBackupInstanceV41Api
/* 通过虚机ID获取虚拟机最新状态，主要获取虚拟机磁盘挂载信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsDetailsBackupInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetailsBackupInstanceV41Api(client *core.CtyunClient) *CtecsDetailsBackupInstanceV41Api {
	return &CtecsDetailsBackupInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup/instance-details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetailsBackupInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetailsBackupInstanceV41Request) (*CtecsDetailsBackupInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDetailsBackupInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetailsBackupInstanceV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>  */
}

type CtecsDetailsBackupInstanceV41Response struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsDetailsBackupInstanceV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsDetailsBackupInstanceV41ReturnObjResponse struct {
	ProjectID       string                                                           `json:"projectID,omitempty"`      /*  企业项目id  */
	AzName          string                                                           `json:"azName,omitempty"`         /*  可用区名称  */
	AttachedVolume  []string                                                         `json:"attachedVolume"`           /*  云硬盘ID列表  */
	Addresses       []*CtecsDetailsBackupInstanceV41ReturnObjAddressesResponse       `json:"addresses"`                /*  网络地址信息  */
	ResourceID      string                                                           `json:"resourceID,omitempty"`     /*  资源ID  */
	InstanceID      string                                                           `json:"instanceID,omitempty"`     /*  云主机ID  */
	DisplayName     string                                                           `json:"displayName,omitempty"`    /*  云主机显示名称  */
	InstanceName    string                                                           `json:"instanceName,omitempty"`   /*  云主机名称  */
	OsType          int32                                                            `json:"osType,omitempty"`         /*  操作系统类型，取值范围：<br />1: linux，<br />2: windows，<br />3: redhat，<br />4: ubuntu，<br />5: centos，<br />6: oracle  */
	InstanceStatus  string                                                           `json:"instanceStatus,omitempty"` /*  云主机状态，取值范围：<br />backingup: 备份中，<br />creating: 创建中，<br />expired: 已到期，<br />freezing: 冻结中，<br />rebuild: 重装，<br />restarting: 重启中，<br />running: 运行中，<br />starting: 开机中，<br />stopped: 已关机，<br />stopping: 关机中，<br />error: 错误，<br />snapshotting: 快照创建中  */
	ExpiredTime     string                                                           `json:"expiredTime,omitempty"`    /*  到期时间  */
	AvailableDay    int32                                                            `json:"availableDay,omitempty"`   /*  可用(天)  */
	UpdatedTime     string                                                           `json:"updatedTime,omitempty"`    /*  更新时间  */
	CreatedTime     string                                                           `json:"createdTime,omitempty"`    /*  创建时间  */
	ZabbixName      string                                                           `json:"zabbixName,omitempty"`     /*  监控对象名称  */
	SecGroupList    []*CtecsDetailsBackupInstanceV41ReturnObjSecGroupListResponse    `json:"secGroupList"`             /*  安全组信息  */
	PrivateIP       string                                                           `json:"privateIP,omitempty"`      /*  内网IPv4地址  */
	PrivateIPv6     string                                                           `json:"privateIPv6,omitempty"`    /*  内网IPv6址  */
	NetworkCardList []*CtecsDetailsBackupInstanceV41ReturnObjNetworkCardListResponse `json:"networkCardList"`          /*  网卡信息  */
	VipInfoList     []*CtecsDetailsBackupInstanceV41ReturnObjVipInfoListResponse     `json:"vipInfoList"`              /*  虚拟IP信息列表  */
	VipCount        int32                                                            `json:"vipCount,omitempty"`       /*  vip数目  */
	AffinityGroup   *CtecsDetailsBackupInstanceV41ReturnObjAffinityGroupResponse     `json:"affinityGroup"`            /*  云主机组信息  */
	Image           *CtecsDetailsBackupInstanceV41ReturnObjImageResponse             `json:"image"`                    /*  镜像信息  */
	Flavor          *CtecsDetailsBackupInstanceV41ReturnObjFlavorResponse            `json:"flavor"`                   /*  云主机规格信息  */
	OnDemand        *bool                                                            `json:"onDemand"`                 /*  付费方式。取值范围：<br>false：包周期。<br>true：按量付费。  */
	VpcName         string                                                           `json:"vpcName,omitempty"`        /*  vpc名称  */
	VpcID           string                                                           `json:"vpcID,omitempty"`          /*  vpc ID  */
	FixedIP         []string                                                         `json:"fixedIP"`                  /*  内网IP  */
	FloatingIP      string                                                           `json:"floatingIP,omitempty"`     /*  公网IP  */
	SubnetIDList    []string                                                         `json:"subnetIDList"`             /*  子网ID列表  */
	KeypairName     string                                                           `json:"keypairName,omitempty"`    /*  密钥对名称  */
	Volumes         []*CtecsDetailsBackupInstanceV41ReturnObjVolumesResponse         `json:"volumes"`                  /*  云硬盘信息  */
}

type CtecsDetailsBackupInstanceV41ReturnObjAddressesResponse struct {
	VpcName     string                                                                `json:"vpcName,omitempty"` /*  vpc名称  */
	AddressList []*CtecsDetailsBackupInstanceV41ReturnObjAddressesAddressListResponse `json:"addressList"`       /*  网络地址列表  */
}

type CtecsDetailsBackupInstanceV41ReturnObjSecGroupListResponse struct {
	SecurityGroupID   string `json:"securityGroupID,omitempty"`   /*  安全组ID  */
	SecurityGroupName string `json:"securityGroupName,omitempty"` /*  安全组名称  */
}

type CtecsDetailsBackupInstanceV41ReturnObjNetworkCardListResponse struct {
	IPv4Address   string   `json:"IPv4Address,omitempty"`   /*  IPv4地址  */
	IPv6Address   []string `json:"IPv6Address"`             /*  IPv6地址列表  */
	SubnetID      string   `json:"subnetID,omitempty"`      /*  子网ID  */
	SubnetCidr    string   `json:"subnetCidr,omitempty"`    /*  子网CIDR信息  */
	IsMaster      *bool    `json:"isMaster"`                /*  是否主网卡，取值范围：<br />true：主网卡，<br />false：扩展网卡  */
	Gateway       string   `json:"gateway,omitempty"`       /*  网关地址  */
	NetworkCardID string   `json:"networkCardID,omitempty"` /*  网卡ID  */
	SecurityGroup []string `json:"securityGroup"`           /*  安全组ID列表  */
}

type CtecsDetailsBackupInstanceV41ReturnObjVipInfoListResponse struct {
	VipID          string `json:"vipID,omitempty"`          /*  虚拟IP的ID  */
	VipAddress     string `json:"vipAddress,omitempty"`     /*  虚拟IP地址  */
	VipBindNicIP   string `json:"vipBindNicIP,omitempty"`   /*  虚拟IP绑定的网卡对应IPv4地址  */
	VipBindNicIPv6 string `json:"vipBindNicIPv6,omitempty"` /*  虚拟IP绑定的网卡对应IPv6地址  */
	NicID          string `json:"nicID,omitempty"`          /*  网卡ID  */
}

type CtecsDetailsBackupInstanceV41ReturnObjAffinityGroupResponse struct {
	Policy            string `json:"policy,omitempty"`            /*  云主机组策略  */
	AffinityGroupName string `json:"affinityGroupName,omitempty"` /*  云主机组名称  */
	AffinityGroupID   string `json:"affinityGroupID,omitempty"`   /*  云主机组ID  */
}

type CtecsDetailsBackupInstanceV41ReturnObjImageResponse struct {
	ImageID   string `json:"imageID,omitempty"`   /*  镜像ID  */
	ImageName string `json:"imageName,omitempty"` /*  镜像名称  */
}

type CtecsDetailsBackupInstanceV41ReturnObjFlavorResponse struct {
	FlavorID     string `json:"flavorID,omitempty"`     /*  规格ID  */
	FlavorName   string `json:"flavorName,omitempty"`   /*  规格名称  */
	FlavorCPU    int32  `json:"flavorCPU,omitempty"`    /*  VCPU  */
	FlavorRAM    int32  `json:"flavorRAM,omitempty"`    /*  内存  */
	GpuType      string `json:"gpuType,omitempty"`      /*  GPU类型，取值范围：T4\V100\V100S\A10\A100\atlas 300i pro\mlu370-s4，支持类型会随着功能升级增加  */
	GpuCount     int32  `json:"gpuCount,omitempty"`     /*  GPU数目  */
	GpuVendor    string `json:"gpuVendor,omitempty"`    /*  GPU名称  */
	VideoMemSize int32  `json:"videoMemSize,omitempty"` /*  GPU显存大小  */
}

type CtecsDetailsBackupInstanceV41ReturnObjVolumesResponse struct {
	IsBootable *bool  `json:"isBootable"`         /*  是否启动盘  */
	DiskSize   int32  `json:"diskSize,omitempty"` /*  云硬盘大小，单位为GB  */
	DiskType   string `json:"diskType,omitempty"` /*  云硬盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO，<br />FAST-SSD：极速型SSD  */
	DiskID     string `json:"diskID,omitempty"`   /*  云硬盘ID  */
	DiskName   string `json:"diskName,omitempty"` /*  云硬盘名称  */
}

type CtecsDetailsBackupInstanceV41ReturnObjAddressesAddressListResponse struct {
	Addr       string `json:"addr,omitempty"`       /*  IP地址  */
	Version    int32  `json:"version,omitempty"`    /*  IP版本  */
	RawType    string `json:"type,omitempty"`       /*  网络类型，取值范围：<br />fixed：内网，<br />floating：公网  */
	IsMaster   *bool  `json:"isMaster"`             /*  网络地址对应网卡是否为主网卡  */
	MacAddress string `json:"macAddress,omitempty"` /*  网络地址对应网卡的mac地址  */
}
