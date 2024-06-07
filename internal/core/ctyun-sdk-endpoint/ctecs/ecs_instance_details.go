package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsInstanceDetailsApi 查询一台云主机详细信息
// https://www.ctyun.cn/document/10026730/10106322
type EcsInstanceDetailsApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsInstanceDetailsApi(client *ctyunsdk.CtyunClient) *EcsInstanceDetailsApi {
	return &EcsInstanceDetailsApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/ecs/instance-details",
		},
	}
}

func (this *EcsInstanceDetailsApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsInstanceDetailsRequest) (*EcsInstanceDetailsResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.
		AddParam("regionID", req.RegionId).
		AddParam("instanceID", req.InstanceId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsInstanceDetailsRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}

	var sgs []EcsInstanceDetailsResultsSecGroupListResponse
	for _, s := range realResponse.SecGroupList {
		sgs = append(sgs, EcsInstanceDetailsResultsSecGroupListResponse{
			SecurityGroupName: s.SecurityGroupName,
			SecurityGroupId:   s.SecurityGroupID,
		})
	}

	var ni []EcsInstanceDetailsResultsNetworkCardListResponse
	for _, n := range realResponse.NetworkCardList {
		ni = append(ni, EcsInstanceDetailsResultsNetworkCardListResponse{
			Ipv4Address:   n.IPv4Address,
			Ipv6Address:   n.IPv6Address,
			SubnetId:      n.SubnetID,
			SubnetCidr:    n.SubnetCidr,
			IsMaster:      n.IsMaster,
			Gateway:       n.Gateway,
			NetworkCardId: n.NetworkCardID,
			SecurityGroup: n.SecurityGroup,
		})
	}

	return &EcsInstanceDetailsResponse{
		ProjectId:       realResponse.ProjectID,
		AzName:          realResponse.AzName,
		AttachedVolume:  realResponse.AttachedVolume,
		ResourceId:      realResponse.ResourceID,
		InstanceId:      realResponse.InstanceID,
		DisplayName:     realResponse.DisplayName,
		InstanceName:    realResponse.InstanceName,
		OsType:          realResponse.OsType,
		InstanceStatus:  realResponse.InstanceStatus,
		ExpiredTime:     realResponse.ExpiredTime,
		AvailableDay:    realResponse.AvailableDay,
		UpdatedTime:     realResponse.UpdatedTime,
		CreatedTime:     realResponse.CreatedTime,
		SecGroupList:    sgs,
		PrivateIp:       realResponse.PrivateIP,
		PrivateIPv6:     realResponse.PrivateIPv6,
		NetworkCardList: ni,
		VipCount:        realResponse.VipCount,
		Image: EcsInstanceDetailsResultsImageResponse{
			ImageId:   realResponse.Image.ImageID,
			ImageName: realResponse.Image.ImageName,
		},
		Flavor: EcsInstanceDetailsResultsFlavorResponse{
			FlavorId:     realResponse.Flavor.FlavorID,
			FlavorName:   realResponse.Flavor.FlavorName,
			FlavorCpu:    realResponse.Flavor.FlavorCPU,
			FlavorRam:    realResponse.Flavor.FlavorRAM,
			GpuType:      realResponse.Flavor.GpuType,
			GpuCount:     realResponse.Flavor.GpuCount,
			GpuVendor:    realResponse.Flavor.GpuVendor,
			VideoMemSize: realResponse.Flavor.VideoMemSize,
		},
		OnDemand:     realResponse.OnDemand,
		VpcName:      realResponse.VpcName,
		VpcId:        realResponse.VpcID,
		FixedIpList:  realResponse.FixedIPList,
		FloatingIp:   realResponse.FloatingIP,
		SubnetIdList: realResponse.SubnetIDList,
		KeypairName:  realResponse.KeypairName,
	}, nil
}

type ecsInstanceDetailsRealResponse struct {
	ProjectID      string   `json:"projectID,omitempty"`
	AzName         string   `json:"azName"`
	AttachedVolume []string `json:"attachedVolume"`
	Addresses      []struct {
		VpcName     string `json:"vpcName"`
		AddressList []struct {
			Addr    string `json:"addr"`
			Version int    `json:"version"`
			Type    string `json:"type"`
		} `json:"addressList"`
	} `json:"addresses"`
	ResourceID     string `json:"resourceID"`
	InstanceID     string `json:"instanceID"`
	DisplayName    string `json:"displayName"`
	InstanceName   string `json:"instanceName"`
	OsType         int    `json:"osType"`
	InstanceStatus string `json:"instanceStatus"`
	ExpiredTime    string `json:"expiredTime"`
	AvailableDay   int    `json:"availableDay"`
	UpdatedTime    string `json:"updatedTime"`
	CreatedTime    string `json:"createdTime"`
	SecGroupList   []struct {
		SecurityGroupName string `json:"securityGroupName"`
		SecurityGroupID   string `json:"securityGroupID"`
	} `json:"secGroupList"`
	PrivateIP       string `json:"privateIP"`
	PrivateIPv6     string `json:"privateIPv6"`
	NetworkCardList []struct {
		IPv4Address   string   `json:"IPv4Address,omitempty"`
		IPv6Address   []string `json:"IPv6Address,omitempty"`
		SubnetID      string   `json:"subnetID,omitempty"`
		SubnetCidr    string   `json:"subnetCidr,omitempty"`
		IsMaster      bool     `json:"isMaster,omitempty"`
		Gateway       string   `json:"gateway,omitempty"`
		NetworkCardID string   `json:"networkCardID,omitempty"`
		SecurityGroup []string `json:"securityGroup,omitempty"`
	} `json:"networkCardList"`
	VipInfoList []struct {
		VipID          string `json:"vipID"`
		VipAddress     string `json:"vipAddress"`
		VipBindNicIP   string `json:"vipBindNicIP"`
		VipBindNicIPv6 string `json:"vipBindNicIPv6"`
		NicID          string `json:"nicID"`
	} `json:"vipInfoList"`
	VipCount      int `json:"vipCount"`
	AffinityGroup struct {
		Policy            string `json:"policy"`
		AffinityGroupName string `json:"affinityGroupName"`
		AffinityGroupID   string `json:"affinityGroupID"`
	} `json:"affinityGroup"`
	Image struct {
		ImageID   string `json:"imageID"`
		ImageName string `json:"imageName"`
	} `json:"image"`
	Flavor struct {
		FlavorID     string `json:"flavorID"`
		FlavorName   string `json:"flavorName"`
		FlavorCPU    int    `json:"flavorCPU"`
		FlavorRAM    int    `json:"flavorRAM"`
		GpuType      string `json:"gpuType"`
		GpuCount     int    `json:"gpuCount"`
		GpuVendor    string `json:"gpuVendor"`
		VideoMemSize int    `json:"videoMemSize"`
	} `json:"flavor"`
	OnDemand     bool     `json:"onDemand"`
	VpcName      string   `json:"vpcName"`
	VpcID        string   `json:"vpcID"`
	FixedIPList  []string `json:"fixedIPList"`
	FloatingIP   string   `json:"floatingIP"`
	SubnetIDList []string `json:"subnetIDList"`
	KeypairName  string   `json:"keypairName"`
}

type EcsInstanceDetailsRequest struct {
	RegionId   string
	InstanceId string
}

type EcsInstanceDetailsResultsSecGroupListResponse struct {
	SecurityGroupName string
	SecurityGroupId   string
}

type EcsInstanceDetailsResultsNetworkCardListResponse struct {
	Ipv4Address   string
	Ipv6Address   []string
	SubnetId      string
	SubnetCidr    string
	IsMaster      bool
	Gateway       string
	NetworkCardId string
	SecurityGroup []string
}

type EcsInstanceDetailsResultsImageResponse struct {
	ImageId   string
	ImageName string
}

type EcsInstanceDetailsResultsFlavorResponse struct {
	FlavorId     string
	FlavorName   string
	FlavorCpu    int
	FlavorRam    int
	GpuType      string
	GpuCount     int
	GpuVendor    string
	VideoMemSize int
}

type EcsInstanceDetailsResponse struct {
	ProjectId       string
	AzName          string
	AttachedVolume  []string
	ResourceId      string
	InstanceId      string
	DisplayName     string
	InstanceName    string
	OsType          int
	InstanceStatus  string
	ExpiredTime     string
	AvailableDay    int
	UpdatedTime     string
	CreatedTime     string
	SecGroupList    []EcsInstanceDetailsResultsSecGroupListResponse
	PrivateIp       string
	PrivateIPv6     string
	NetworkCardList []EcsInstanceDetailsResultsNetworkCardListResponse
	VipCount        int
	Image           EcsInstanceDetailsResultsImageResponse
	Flavor          EcsInstanceDetailsResultsFlavorResponse
	OnDemand        bool
	VpcName         string
	VpcId           string
	FixedIpList     []string
	FloatingIp      string
	SubnetIdList    []string
	KeypairName     string
}
