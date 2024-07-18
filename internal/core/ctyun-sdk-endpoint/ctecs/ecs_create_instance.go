package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsCreateInstanceApi 创建一台按量付费或包年包月的云主机
// https://www.ctyun.cn/document/10026730/10106580
type EcsCreateInstanceApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsCreateInstanceApi(client *ctyunsdk.CtyunClient) *EcsCreateInstanceApi {
	return &EcsCreateInstanceApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/create-instance",
		},
	}
}

func (this *EcsCreateInstanceApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsCreateInstanceRequest) (*EcsCreateInstanceResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	var networkCardList []ecsCreateInstanceNetworkCardListRealRequest
	for _, request := range req.NetworkCardList {
		networkCardList = append(networkCardList, ecsCreateInstanceNetworkCardListRealRequest{
			NicName:  request.NicName,
			FixedIP:  request.FixedIp,
			IsMaster: request.IsMaster,
			SubnetID: request.SubnetId,
		})
	}
	_, err := builder.WriteJson(&ecsCreateInstanceRealRequest{
		ClientToken:     req.ClientToken,
		RegionID:        req.RegionId,
		AzName:          req.AzName,
		InstanceName:    req.InstanceName,
		DisplayName:     req.DisplayName,
		FlavorID:        req.FlavorId,
		ImageType:       req.ImageType,
		ImageID:         req.ImageId,
		BootDiskType:    req.BootDiskType,
		BootDiskSize:    req.BootDiskSize,
		VpcID:           req.VpcId,
		OnDemand:        req.OnDemand,
		NetworkCardList: networkCardList,
		ExtIP:           req.ExtIp,
		ProjectID:       req.ProjectId,
		SecGroupList:    req.SecGroupList,
		KeyPairId:       req.KeyPairId,
		UserPassword:    req.UserPassword,
		CycleCount:      req.CycleCount,
		CycleType:       req.CycleType,
		AutoRenewStatus: req.AutoRenewStatus,
		UserData:        req.UserData,
		MonitorService:  req.MonitorService,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsCreateInstanceRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	return &EcsCreateInstanceResponse{
		RegionId:         realResponse.RegionID,
		MasterOrderId:    realResponse.MasterOrderID,
		MasterResourceId: realResponse.MasterResourceID,
		MasterOrderNo:    realResponse.MasterOrderNO,
	}, nil
}

type ecsCreateInstanceNetworkCardListRealRequest struct {
	NicName  string `json:"nicName,omitempty"`
	FixedIP  string `json:"fixedIP,omitempty"`
	IsMaster bool   `json:"isMaster"`
	SubnetID string `json:"subnetID"`
}

type ecsCreateInstanceRealRequest struct {
	ClientToken     string                                        `json:"clientToken"`
	RegionID        string                                        `json:"regionID"`
	AzName          string                                        `json:"azName"`
	InstanceName    string                                        `json:"instanceName"`
	DisplayName     string                                        `json:"displayName"`
	FlavorID        string                                        `json:"flavorID"`
	ImageType       *int                                          `json:"imageType"`
	ImageID         string                                        `json:"imageID"`
	BootDiskType    string                                        `json:"bootDiskType"`
	BootDiskSize    *int                                          `json:"bootDiskSize"`
	VpcID           string                                        `json:"vpcID"`
	OnDemand        *bool                                         `json:"onDemand"`
	NetworkCardList []ecsCreateInstanceNetworkCardListRealRequest `json:"networkCardList"`
	ExtIP           string                                        `json:"extIP"`
	ProjectID       string                                        `json:"projectID,omitempty"`
	SecGroupList    []string                                      `json:"secGroupList,omitempty"`
	KeyPairId       string                                        `json:"keyPairId,omitempty"`
	UserPassword    string                                        `json:"userPassword,omitempty"`
	CycleCount      *int                                          `json:"cycleCount,omitempty"`
	CycleType       string                                        `json:"cycleType,omitempty"`
	AutoRenewStatus *int                                          `json:"autoRenewStatus,omitempty"`
	UserData        string                                        `json:"userData,omitempty"`
	MonitorService  *bool                                         `json:"monitorService,omitempty"`
}

type ecsCreateInstanceRealResponse struct {
	RegionID         string `json:"regionID"`
	MasterOrderID    string `json:"masterOrderID"`
	MasterResourceID string `json:"masterResourceID"`
	MasterOrderNO    string `json:"masterOrderNO"`
}

type EcsCreateInstanceRequest struct {
	ClientToken     string                                    // 客户端存根
	RegionId        string                                    // 资源池ID
	AzName          string                                    // 可用区名称
	InstanceName    string                                    // 云主机名称，不可以使用已存在的云主机名称
	DisplayName     string                                    // 云主机显示名称
	FlavorId        string                                    // 云主机规格ID
	ImageType       *int                                      // 镜像类型
	ImageId         string                                    // 镜像ID
	BootDiskType    string                                    // 系统盘类型，取值范围： SATA（普通IO）， SAS（高IO）， SSD（超高IO）， SSD-genric（通用型SSD）， FAST-SSD（极速型SSD），您可以查看磁盘类型及性能介绍来了解磁盘类型及其对应性能指标
	BootDiskSize    *int                                      // 系统盘大小单位为GiB，取值范围：[40, 32768]
	VpcId           string                                    // 虚拟私有云ID
	OnDemand        *bool                                     // 购买方式，取值范围： false（按周期）， true（按需）
	NetworkCardList []EcsCreateInstanceNetworkCardListRequest // 网卡信息列表
	ExtIp           string                                    // 是否使用弹性公网IP，取值范围： 0（不使用）， 1（自动分配）， 2（使用已有）
	ProjectId       string                                    // 企业项目id
	SecGroupList    []string                                  // 安全组ID列表
	KeyPairId       string                                    // 密钥对ID
	UserPassword    string                                    // 用户密码
	CycleCount      *int                                      // 订购时长
	CycleType       string                                    // 订购周期类型
	AutoRenewStatus *int                                      // 是否自动续订，取值范围： 0（不续费）， 1（自动续费）， 注：按月购买，自动续订周期为3个月；按年购买，自动续订周期为1年
	UserData        string                                    // 用户自定义数据，需要以Base64方式编码，Base64编码后的长度限制为1-16384字符
	MonitorService  *bool                                     // 监控参数，支持通过该参数指定云主机在创建后是否开启详细监控，取值范围： false（不开启），true（开启）；注：若指定该参数为true或不指定该参数，云主机内默认开启最新详细监控服务。若指定该参数为false，默认公共镜像不开启最新监控服务；私有镜像使用镜像中保留的监控服务。说明：仅部分资源池支持monitorService参数，
}

type EcsCreateInstanceNetworkCardListRequest struct {
	NicName  string // 长度2~32，支持拉丁字母、中文、数字、下划线、连字符，中文或英文字母开头，不能以http:或https:开头
	FixedIp  string // 内网IPv4地址，注：不可使用已占用IP
	IsMaster bool   // 是否主网卡，取值范围： true（表示主网卡）， false（表示扩展网卡） 注：只能含有一个主网卡
	SubnetId string // 子网ID
}

type EcsCreateInstanceResponse struct {
	RegionId         string
	MasterOrderId    string
	MasterResourceId string
	MasterOrderNo    string
}
