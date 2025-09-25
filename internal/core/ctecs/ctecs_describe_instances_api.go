package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDescribeInstancesApi
/* 该接口提供用户多台云主机信息查询功能，用户可以根据此接口的返回值得到多台云主机信息。该接口相较于/v4/ecs/list-instances提供更精简的云主机信息，拥有更高的查找效率<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />&emsp;&emsp;匹配查找：可以通过部分字段进行匹配筛选数据，无符合条件的为空，在指定多台云主机ID的情况下，只返回匹配到的云主机信息。推荐每次使用单个条件查找
 */type CtecsDescribeInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDescribeInstancesApi(client *core.CtyunClient) *CtecsDescribeInstancesApi {
	return &CtecsDescribeInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/describe-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDescribeInstancesApi) Do(ctx context.Context, credential core.Credential, req *CtecsDescribeInstancesRequest) (*CtecsDescribeInstancesResponse, error) {
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
	var resp CtecsDescribeInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDescribeInstancesRequest struct {
	RegionID        string                                    `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName          string                                    `json:"azName,omitempty"`          /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
	ProjectID       string                                    `json:"projectID,omitempty"`       /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	PageNo          int32                                     `json:"pageNo,omitempty"`          /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize        int32                                     `json:"pageSize,omitempty"`        /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	State           string                                    `json:"state,omitempty"`           /*  云主机状态，详见枚举值表<br />注：该参数大小写不敏感（如active可填写为ACTIVE）  */
	Keyword         string                                    `json:"keyword,omitempty"`         /*  关键字，对部分参数进行模糊查询，包含：instanceName、displayName、instanceID、privateIP  */
	InstanceName    string                                    `json:"instanceName,omitempty"`    /*  云主机名称，精准匹配  */
	InstanceIDList  string                                    `json:"instanceIDList,omitempty"`  /*  云主机ID列表，多台使用英文逗号分割，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br/><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	SecurityGroupID string                                    `json:"securityGroupID,omitempty"` /*  安全组ID，模糊匹配，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	LabelList       []*CtecsDescribeInstancesLabelListRequest `json:"labelList"`                 /*  标签信息列表  */
}

type CtecsDescribeInstancesLabelListRequest struct {
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签键，长度限制1~32字符  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签值，长度限制1~32字符  */
}

type CtecsDescribeInstancesResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDescribeInstancesReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDescribeInstancesReturnObjResponse struct {
	CurrentCount int32                                             `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                             `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                             `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsDescribeInstancesReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsDescribeInstancesReturnObjResultsResponse struct {
	ProjectID           string                                                        `json:"projectID,omitempty"`           /*  企业项目ID  */
	AzName              string                                                        `json:"azName,omitempty"`              /*  可用区名称  */
	AzDisplayName       string                                                        `json:"azDisplayName,omitempty"`       /*  可用区展示名称  */
	AttachedVolume      []string                                                      `json:"attachedVolume"`                /*  云硬盘ID列表  */
	Addresses           []*CtecsDescribeInstancesReturnObjResultsAddressesResponse    `json:"addresses"`                     /*  网络地址信息  */
	InstanceID          string                                                        `json:"instanceID,omitempty"`          /*  云主机ID  */
	DisplayName         string                                                        `json:"displayName,omitempty"`         /*  云主机显示名称  */
	InstanceName        string                                                        `json:"instanceName,omitempty"`        /*  云主机名称  */
	OsType              int32                                                         `json:"osType,omitempty"`              /*  操作系统类型，详见枚举值表  */
	InstanceStatus      string                                                        `json:"instanceStatus,omitempty"`      /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	ExpiredTime         string                                                        `json:"expiredTime,omitempty"`         /*  到期时间  */
	CreatedTime         string                                                        `json:"createdTime,omitempty"`         /*  创建时间  */
	SecGroupList        []*CtecsDescribeInstancesReturnObjResultsSecGroupListResponse `json:"secGroupList"`                  /*  安全组信息列表  */
	VipInfoList         []*CtecsDescribeInstancesReturnObjResultsVipInfoListResponse  `json:"vipInfoList"`                   /*  虚拟IP信息列表  */
	AffinityGroup       *CtecsDescribeInstancesReturnObjResultsAffinityGroupResponse  `json:"affinityGroup"`                 /*  云主机组信息  */
	Image               *CtecsDescribeInstancesReturnObjResultsImageResponse          `json:"image"`                         /*  镜像信息  */
	Flavor              *CtecsDescribeInstancesReturnObjResultsFlavorResponse         `json:"flavor"`                        /*  云主机规格信息  */
	OnDemand            *bool                                                         `json:"onDemand"`                      /*  付费方式，取值范围：<br />true（按量付费），<br />false（包周期）  */
	KeypairName         string                                                        `json:"keypairName,omitempty"`         /*  密钥对名称  */
	NetworkInfo         []*CtecsDescribeInstancesReturnObjResultsNetworkInfoResponse  `json:"networkInfo"`                   /*  网络信息  */
	DelegateName        string                                                        `json:"delegateName,omitempty"`        /*  委托名称，注：委托绑定目前仅支持多可用区类型资源池，非可用区资源池为空字符串  */
	DeletionProtection  *bool                                                         `json:"deletionProtection"`            /*  是否开启实例删除保护  */
	InstanceDescription string                                                        `json:"instanceDescription,omitempty"` /*  云主机描述信息  */
}

type CtecsDescribeInstancesReturnObjResultsAddressesResponse struct {
	VpcName     string                                                                `json:"vpcName,omitempty"` /*  vpc名称  */
	AddressList []*CtecsDescribeInstancesReturnObjResultsAddressesAddressListResponse `json:"addressList"`       /*  网络地址列表  */
}

type CtecsDescribeInstancesReturnObjResultsSecGroupListResponse struct {
	SecurityGroupID   string `json:"securityGroupID,omitempty"`   /*  安全组ID  */
	SecurityGroupName string `json:"securityGroupName,omitempty"` /*  安全组名称  */
}

type CtecsDescribeInstancesReturnObjResultsVipInfoListResponse struct {
	VipID          string `json:"vipID,omitempty"`          /*  虚拟IP的ID  */
	VipAddress     string `json:"vipAddress,omitempty"`     /*  虚拟IP地址  */
	VipBindNicIP   string `json:"vipBindNicIP,omitempty"`   /*  虚拟IP绑定的网卡对应IPv4地址  */
	VipBindNicIPv6 string `json:"vipBindNicIPv6,omitempty"` /*  虚拟IP绑定的网卡对应IPv6地址  */
	NicID          string `json:"nicID,omitempty"`          /*  网卡ID  */
}

type CtecsDescribeInstancesReturnObjResultsAffinityGroupResponse struct {
	Policy            string `json:"policy,omitempty"`            /*  云主机组策略  */
	AffinityGroupName string `json:"affinityGroupName,omitempty"` /*  云主机组名称  */
	AffinityGroupID   string `json:"affinityGroupID,omitempty"`   /*  云主机组ID  */
}

type CtecsDescribeInstancesReturnObjResultsImageResponse struct {
	ImageID   string `json:"imageID,omitempty"`   /*  镜像ID  */
	ImageName string `json:"imageName,omitempty"` /*  镜像名称  */
}

type CtecsDescribeInstancesReturnObjResultsFlavorResponse struct {
	FlavorID     string `json:"flavorID,omitempty"`     /*  规格ID  */
	FlavorName   string `json:"flavorName,omitempty"`   /*  规格名称  */
	FlavorCPU    int32  `json:"flavorCPU,omitempty"`    /*  VCPU  */
	FlavorRAM    int32  `json:"flavorRAM,omitempty"`    /*  内存  */
	GpuType      string `json:"gpuType,omitempty"`      /*  GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加  */
	GpuCount     int32  `json:"gpuCount,omitempty"`     /*  GPU数目  */
	GpuVendor    string `json:"gpuVendor,omitempty"`    /*  GPU名称  */
	VideoMemSize int32  `json:"videoMemSize,omitempty"` /*  显存大小  */
}

type CtecsDescribeInstancesReturnObjResultsNetworkInfoResponse struct {
	SubnetID  string                 `json:"subnetID,omitempty"`  /*  子网ID  */
	IpAddress string                 `json:"ipAddress,omitempty"` /*  IP地址  */
	BoundType map[string]interface{} `json:"boundType"`           /*  绑定类型  */
}

type CtecsDescribeInstancesReturnObjResultsAddressesAddressListResponse struct {
	Addr       string `json:"addr,omitempty"`       /*  IP地址  */
	Version    int32  `json:"version,omitempty"`    /*  IP版本  */
	RawType    string `json:"type,omitempty"`       /*  网络类型，取值范围：<br />fixed（内网），<br />floating（弹性公网）  */
	IsMaster   *bool  `json:"isMaster"`             /*  是否为主网卡  */
	MacAddress string `json:"macAddress,omitempty"` /*  mac地址  */
}

type CtecsDescribeInstancesReturnObjResultsNetworkInfoBoundTypeResponse struct{}
