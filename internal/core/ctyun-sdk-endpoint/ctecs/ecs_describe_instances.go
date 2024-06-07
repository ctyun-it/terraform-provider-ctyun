package ctecs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EcsDescribeInstancesApi 查询一台或多台云主机信息
// https://www.ctyun.cn/document/10026730/10189239
type EcsDescribeInstancesApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEcsDescribeInstancesApi(client *ctyunsdk.CtyunClient) *EcsDescribeInstancesApi {
	return &EcsDescribeInstancesApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ecs/describe-instances",
		},
	}
}

func (this *EcsDescribeInstancesApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EcsDescribeInstancesRequest) (*EcsDescribeInstancesResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ecsDescribeInstancesRealRequest{
		RegionID:        req.RegionId,
		AzName:          req.AzName,
		ProjectID:       req.ProjectId,
		PageNo:          req.PageNo,
		PageSize:        req.PageSize,
		State:           req.State,
		Keyword:         req.Keyword,
		InstanceName:    req.InstanceName,
		InstanceIDList:  req.InstanceIdList,
		SecurityGroupID: req.SecurityGroupId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtecs, builder)
	if err != nil {
		return nil, err
	}

	var realResponse ecsDescribeInstancesRealResponse
	err = response.ParseByStandardModelWithCheck(&realResponse)
	if err != nil {
		return nil, err
	}
	var results []EcsDescribeInstancesResultsResponse
	for _, result := range realResponse.Results {
		var sgs []EcsDescribeInstancesResultsSecGroupListResponse
		for _, s := range result.SecGroupList {
			sgs = append(sgs, EcsDescribeInstancesResultsSecGroupListResponse{
				SecurityGroupName: s.SecurityGroupName,
				SecurityGroupId:   s.SecurityGroupID,
			})
		}

		var ni []EcsDescribeInstancesResultsNetworkInfoResponse
		for _, n := range result.NetworkInfo {
			ni = append(ni, EcsDescribeInstancesResultsNetworkInfoResponse{
				SubnetId:  n.SubnetID,
				IpAddress: n.IpAddress,
			})
		}

		results = append(results, EcsDescribeInstancesResultsResponse{
			AzName:         result.AzName,
			AzDisplayName:  result.AzDisplayName,
			AttachedVolume: result.AttachedVolume,
			ExpiredTime:    result.ExpiredTime,
			CreatedTime:    result.CreatedTime,
			ProjectId:      result.ProjectID,
			InstanceId:     result.InstanceID,
			DisplayName:    result.DisplayName,
			InstanceName:   result.InstanceName,
			OsType:         result.OsType,
			InstanceStatus: result.InstanceStatus,
			OnDemand:       result.OnDemand,
			KeypairName:    result.KeypairName,
			SecGroupList:   sgs,
			NetworkInfo:    ni,
			Image: EcsDescribeInstancesResultsImageResponse{
				ImageId:   result.Image.ImageID,
				ImageName: result.Image.ImageName,
			},
			Flavor: EcsDescribeInstancesResultsFlavorResponse{
				FlavorId:     result.Flavor.FlavorID,
				FlavorName:   result.Flavor.FlavorName,
				FlavorCpu:    result.Flavor.FlavorCPU,
				FlavorRam:    result.Flavor.FlavorRAM,
				GpuType:      result.Flavor.GpuType,
				GpuCount:     result.Flavor.GpuCount,
				GpuVendor:    result.Flavor.GpuVendor,
				VideoMemSize: result.Flavor.VideoMemSize,
			},
		})
	}
	return &EcsDescribeInstancesResponse{
		CurrentCount: realResponse.CurrentCount,
		TotalCount:   realResponse.TotalCount,
		TotalPage:    realResponse.TotalPage,
		Results:      results,
	}, nil
}

type ecsDescribeInstancesRealResponse struct {
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
	Results      []struct {
		AzName         string   `json:"azName"`
		AzDisplayName  string   `json:"azDisplayName"`
		AttachedVolume []string `json:"attachedVolume"`
		ExpiredTime    string   `json:"expiredTime"`
		CreatedTime    string   `json:"createdTime"`
		ProjectID      string   `json:"projectID,omitempty"`
		InstanceID     string   `json:"instanceID"`
		DisplayName    string   `json:"displayName"`
		InstanceName   string   `json:"instanceName"`
		OsType         int      `json:"osType"`
		InstanceStatus string   `json:"instanceStatus"`
		OnDemand       bool     `json:"onDemand"`
		KeypairName    string   `json:"keypairName"`
		Addresses      []struct {
			VpcName     string `json:"vpcName"`
			AddressList []struct {
				Addr    string `json:"addr"`
				Version int    `json:"version"`
				Type    string `json:"type"`
			} `json:"addressList"`
		} `json:"addresses"`
		SecGroupList []struct {
			SecurityGroupName string `json:"securityGroupName"`
			SecurityGroupID   string `json:"securityGroupID"`
		} `json:"secGroupList"`
		VipInfoList []struct {
			VipID          string `json:"vipID"`
			VipAddress     string `json:"vipAddress"`
			VipBindNicIP   string `json:"vipBindNicIP"`
			VipBindNicIPv6 string `json:"vipBindNicIPv6"`
			NicID          string `json:"nicID"`
		} `json:"vipInfoList"`
		NetworkInfo []struct {
			SubnetID  string `json:"subnetID"`
			IpAddress string `json:"ipAddress"`
		} `json:"networkInfo"`
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
	} `json:"results"`
}

type ecsDescribeInstancesRealRequest struct {
	RegionID        string `json:"regionID"`
	AzName          string `json:"azName"`
	ProjectID       string `json:"projectID,omitempty"`
	PageNo          int    `json:"pageNo"`
	PageSize        int    `json:"pageSize"`
	State           string `json:"state"`
	Keyword         string `json:"keyword"`
	InstanceName    string `json:"instanceName"`
	InstanceIDList  string `json:"instanceIDList"`
	SecurityGroupID string `json:"securityGroupID"`
}

type EcsDescribeInstancesRequest struct {
	RegionId        string
	AzName          string
	ProjectId       string
	PageNo          int
	PageSize        int
	State           string
	Keyword         string
	InstanceName    string
	InstanceIdList  string
	SecurityGroupId string
}

type EcsDescribeInstancesResultsSecGroupListResponse struct {
	SecurityGroupName string
	SecurityGroupId   string
}

type EcsDescribeInstancesResultsNetworkInfoResponse struct {
	SubnetId  string
	IpAddress string
}

type EcsDescribeInstancesResultsImageResponse struct {
	ImageId   string
	ImageName string
}

type EcsDescribeInstancesResultsFlavorResponse struct {
	FlavorId     string
	FlavorName   string
	FlavorCpu    int
	FlavorRam    int
	GpuType      string
	GpuCount     int
	GpuVendor    string
	VideoMemSize int
}

type EcsDescribeInstancesResultsResponse struct {
	AzName         string
	AzDisplayName  string
	AttachedVolume []string
	ExpiredTime    string
	CreatedTime    string
	ProjectId      string
	InstanceId     string
	DisplayName    string
	InstanceName   string
	OsType         int
	InstanceStatus string
	OnDemand       bool
	KeypairName    string
	SecGroupList   []EcsDescribeInstancesResultsSecGroupListResponse
	NetworkInfo    []EcsDescribeInstancesResultsNetworkInfoResponse
	Image          EcsDescribeInstancesResultsImageResponse
	Flavor         EcsDescribeInstancesResultsFlavorResponse
}

type EcsDescribeInstancesResponse struct {
	CurrentCount int
	TotalCount   int
	TotalPage    int
	Results      []EcsDescribeInstancesResultsResponse
}
