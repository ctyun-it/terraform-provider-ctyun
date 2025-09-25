package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CcseListNodePoolsApi
/* 调用该接口查看节点池列表。
 */type CcseListNodePoolsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListNodePoolsApi(client *core.CtyunClient) *CcseListNodePoolsApi {
	return &CcseListNodePoolsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepools",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListNodePoolsApi) Do(ctx context.Context, credential core.Credential, req *CcseListNodePoolsRequest) (*CcseListNodePoolsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.NodePoolName != "" {
		ctReq.AddParam("nodePoolName", req.NodePoolName)
	}
	if req.PageNow != 0 {
		ctReq.AddParam("pageNow", strconv.FormatInt(int64(req.PageNow), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseListNodePoolsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListNodePoolsRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	NodePoolName string /*  节点池名称  */
	PageNow      int32  /*  当前页码  */
	PageSize     int32  /*  每页条数  */
}

type CcseListNodePoolsResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                              `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseListNodePoolsReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                              `json:"error,omitempty"`      /*  错误码  */
}

type CcseListNodePoolsReturnObjResponse struct {
	Records []*CcseListNodePoolsReturnObjRecordsResponse `json:"records"`           /*  记录列表  */
	Total   int32                                        `json:"total,omitempty"`   /*  总条数  */
	Size    int32                                        `json:"size,omitempty"`    /*  每页条数  */
	Current int32                                        `json:"current,omitempty"` /*  当前页码  */
	Pages   int32                                        `json:"pages,omitempty"`   /*  总页数  */
}

type CcseListNodePoolsReturnObjRecordsResponse struct {
	Id                       string                                                `json:"id,omitempty"`                       /*  节点池id  */
	NodePoolName             string                                                `json:"nodePoolName,omitempty"`             /*  节点池名称  */
	BillMode                 string                                                `json:"billMode,omitempty"`                 /*  付费类型  */
	Description              string                                                `json:"description,omitempty"`              /*  描述  */
	NodeTotalNum             int32                                                 `json:"nodeTotalNum,omitempty"`             /*  节点总数  */
	NormalNodeNum            int32                                                 `json:"normalNodeNum,omitempty"`            /*  正常节点数  */
	UnNormalNodeNum          int32                                                 `json:"unNormalNodeNum,omitempty"`          /*  异常节点数  */
	Nodes                    []*CcseListNodePoolsReturnObjRecordsNodesResponse     `json:"nodes"`                              /*  节点池节点信息  */
	Runtime                  string                                                `json:"runtime,omitempty"`                  /*  运行时  */
	RuntimeVersion           string                                                `json:"runtimeVersion,omitempty"`           /*  运行时版本  */
	NodeGroup                string                                                `json:"nodeGroup,omitempty"`                /*  伸缩组  */
	VmSpecName               string                                                `json:"vmSpecName,omitempty"`               /*  节点规格  */
	VmType                   string                                                `json:"vmType,omitempty"`                   /*  节点规格类型  */
	Cpu                      int32                                                 `json:"cpu,omitempty"`                      /*  cpu  */
	Memory                   int32                                                 `json:"memory,omitempty"`                   /*  内存  */
	DataDisks                []*CcseListNodePoolsReturnObjRecordsDataDisksResponse `json:"dataDisks"`                          /*  数据盘  */
	Status                   string                                                `json:"status,omitempty"`                   /*  状态  */
	MaxPodNum                int32                                                 `json:"maxPodNum,omitempty"`                /*  最大pod数  */
	Labels                   *CcseListNodePoolsReturnObjRecordsLabelsResponse      `json:"labels"`                             /*  标签  */
	Taints                   []*CcseListNodePoolsReturnObjRecordsTaintsResponse    `json:"taints"`                             /*  污点  */
	IsDefault                *bool                                                 `json:"isDefault"`                          /*  是否默认  */
	EnableAutoScale          *bool                                                 `json:"enableAutoScale"`                    /*  是否自动弹性伸缩  */
	MaxNum                   int32                                                 `json:"maxNum,omitempty"`                   /*  伸缩组最大数量  */
	MinNum                   int32                                                 `json:"minNum,omitempty"`                   /*  伸缩组最小数量  */
	KubeletArgs              *CcseListNodePoolsReturnObjRecordsKubeletArgsResponse `json:"kubeletArgs"`                        /*  kubelet参数  */
	CreatedTime              string                                                `json:"createdTime,omitempty"`              /*  创建时间  */
	UpdatedTime              string                                                `json:"updatedTime,omitempty"`              /*  更新时间  */
	Vpc                      string                                                `json:"vpc,omitempty"`                      /*  vpc名称  */
	Subnetwork               string                                                `json:"subnetwork,omitempty"`               /*  子网  */
	SecruityGroupName        string                                                `json:"secruityGroupName,omitempty"`        /*  安全组名称  */
	ImageName                string                                                `json:"imageName,omitempty"`                /*  操作系统镜像名称  */
	ImageUuid                string                                                `json:"imageUuid,omitempty"`                /*  操作系统镜像ID  */
	EcsPasswd                string                                                `json:"ecsPasswd,omitempty"`                /*  云主机密码  */
	LoginType                string                                                `json:"loginType,omitempty"`                /*  云主机密码类型  */
	KeyName                  string                                                `json:"keyName,omitempty"`                  /*  ebm密钥对使用字段  */
	KeyPairId                string                                                `json:"keyPairId,omitempty"`                /*  ecs密钥对使用字段  */
	SysDiskType              string                                                `json:"sysDiskType,omitempty"`              /*  节点系统盘类型  */
	SysDiskSize              int32                                                 `json:"sysDiskSize,omitempty"`              /*  节点系统盘大小  */
	KubeletDirectory         string                                                `json:"kubeletDirectory,omitempty"`         /*  自定义kubelet目录  */
	ContainerDataDirectory   string                                                `json:"containerDataDirectory,omitempty"`   /*  自定义容器数据目录  */
	VisibilityPostHostScript string                                                `json:"visibilityPostHostScript,omitempty"` /*  部署后执行自定义脚本  */
	VisibilityHostScript     string                                                `json:"visibilityHostScript,omitempty"`     /*  部署前执行自定义脚本  */
	AzInfo                   []*CcseListNodePoolsReturnObjRecordsAzInfoResponse    `json:"azInfo"`                             /*  可用区  */
	SubnetUuid               string                                                `json:"subnetUuid,omitempty"`               /*  子网UUID  */
	NodePoolType             int32                                                 `json:"nodePoolType,omitempty"`             /*  节点池类型  */
	BillingMode              string                                                `json:"billingMode,omitempty"`              /*  订单类型 1-包年包月 2-按需计费  */
	CycleCount               int32                                                 `json:"cycleCount,omitempty"`               /*  订购时长  */
	CycleType                string                                                `json:"cycleType,omitempty"`                /*  订购周期类型 MONTH-月 YEAR-年  */
	AutoRenewStatus          int32                                                 `json:"autoRenewStatus,omitempty"`          /*  是否自动续订 0-否 1-是  */
	DefinedHostnameEnable    int32                                                 `json:"definedHostnameEnable,omitempty"`    /*  是否自定义节点名称 0-否 1-是  */
	HostNamePrefix           string                                                `json:"hostNamePrefix,omitempty"`           /*  自定义主机名前缀  */
	HostNamePostfix          string                                                `json:"hostNamePostfix,omitempty"`          /*  自定义主机名后缀  */
}

type CcseListNodePoolsReturnObjRecordsNodesResponse struct {
	ClusterId            string                                                  `json:"clusterId,omitempty"`        /*  集群ID  */
	NodeName             string                                                  `json:"nodeName,omitempty"`         /*  节点名称  */
	NodeType             int32                                                   `json:"nodeType,omitempty"`         /*  节点类型，取值：<br />1：master <br />2：slave  */
	NodeStatus           string                                                  `json:"nodeStatus,omitempty"`       /*  节点状态，取值：<br/>normal：健康。<br/>abnormal：异常。<br/>expulsion：驱逐中。  */
	IsSchedule           int32                                                   `json:"isSchedule,omitempty"`       /*  是否调度，取值： 1：是。 <br />0：否。  */
	IsEvict              int32                                                   `json:"isEvict,omitempty"`          /*  是否驱逐，取值： 1：是。 <br />0：否。  */
	DockerDataPath       string                                                  `json:"dockerDataPath,omitempty"`   /*  docker数据目录  */
	CreatedTime          string                                                  `json:"createdTime,omitempty"`      /*  创建时间。  */
	HostIp               string                                                  `json:"hostIp,omitempty"`           /*  主机管理ip。  */
	HostIpv6             string                                                  `json:"hostIpv6,omitempty"`         /*  主机管理ipv6。  */
	HostExtraIp          string                                                  `json:"hostExtraIp,omitempty"`      /*  主机业务ip。  */
	HostExtraIpv6        string                                                  `json:"hostExtraIpv6,omitempty"`    /*  主机业务ipv6。  */
	Cpu                  string                                                  `json:"cpu,omitempty"`              /*  cpu核数  */
	CpuUseRate           int32                                                   `json:"cpuUseRate,omitempty"`       /*  cpu使用率%  */
	Memory               string                                                  `json:"memory,omitempty"`           /*  内存  */
	MemoryUseRate        int32                                                   `json:"memoryUseRate,omitempty"`    /*  内存使用率%  */
	Disk                 string                                                  `json:"disk,omitempty"`             /*  磁盘GiB  */
	DiskUseRate          int32                                                   `json:"diskUseRate,omitempty"`      /*  磁盘使用率%  */
	KubeletVersion       string                                                  `json:"kubeletVersion,omitempty"`   /*  Kubelet 版本  */
	PodCidr              string                                                  `json:"podCidr,omitempty"`          /*  Pod CIDR  */
	KernelVersion        string                                                  `json:"kernelVersion,omitempty"`    /*  内核版本  */
	OsImageVersion       string                                                  `json:"osImageVersion,omitempty"`   /*  OS 镜像  */
	KubeProxyVersion     string                                                  `json:"kubeProxyVersion,omitempty"` /*  KubeProxy 版本  */
	DockerVersion        string                                                  `json:"dockerVersion,omitempty"`    /*  容器版本  */
	IsSecure             int32                                                   `json:"isSecure,omitempty"`         /*  是否安全节点，取值：<br />1：是。 <br />0：否。  */
	Taints               []*CcseListNodePoolsReturnObjRecordsNodesTaintsResponse `json:"taints"`                     /*  污点  */
	ChannelLabel         string                                                  `json:"channelLabel,omitempty"`     /*  渠道标签  */
	ZoneName             string                                                  `json:"zoneName,omitempty"`         /*  可用区名称  */
	CloudHostId          string                                                  `json:"cloudHostId,omitempty"`      /*  云主机id  */
	Labels               *CcseListNodePoolsReturnObjRecordsNodesLabelsResponse   `json:"labels"`                     /*  标签  */
	EcsId                string                                                  `json:"ecsId,omitempty"`            /*  paas平台云主机id  */
	VmSpecName           string                                                  `json:"vmSpecName,omitempty"`       /*  节点规格名称  */
	Architecture         string                                                  `json:"architecture,omitempty"`     /*  架构  */
	ZoneCode             string                                                  `json:"zoneCode,omitempty"`         /*  可用区编码  */
	BizState             int32                                                   `json:"bizState,omitempty"`         /*  集群状态  */
	NeedShowHostRoom     *bool                                                   `json:"needShowHostRoom"`           /*  是否显示机房信息  */
	RoomName             string                                                  `json:"roomName,omitempty"`         /*  机房名称  */
	FrameCode            string                                                  `json:"frameCode,omitempty"`        /*  机柜编码  */
	PhysicalPosition     string                                                  `json:"physicalPosition,omitempty"` /*  物理位置  */
	RoomInfo             string                                                  `json:"roomInfo,omitempty"`         /*  机房信息  */
	IsBootstrap          *bool                                                   `json:"isBootstrap"`                /*  是否为bootstrap节点  */
	HostType             string                                                  `json:"hostType,omitempty"`         /*  host类型  */
	ResourceId           string                                                  `json:"resourceId,omitempty"`       /*  资源id  */
	NodeScaleDownProject *bool                                                   `json:"nodeScaleDownProject"`       /*  开启/关闭节点缩容保护  */
	LoginType            string                                                  `json:"loginType,omitempty"`        /*  登陆类型  */
	VirtualNodeId        string                                                  `json:"virtualNodeId,omitempty"`    /*  虚拟节点ID  */
}

type CcseListNodePoolsReturnObjRecordsDataDisksResponse struct {
	DiskSpecName string `json:"diskSpecName,omitempty"` /*  数据盘规格名称  */
	Size         int64  `json:"size,omitempty"`         /*  数据盘大小  */
}

type CcseListNodePoolsReturnObjRecordsLabelsResponse struct{}

type CcseListNodePoolsReturnObjRecordsTaintsResponse struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略  */
}

type CcseListNodePoolsReturnObjRecordsKubeletArgsResponse struct{}

type CcseListNodePoolsReturnObjRecordsAzInfoResponse struct {
	AzId   int64  `json:"azId,omitempty"`   /*  可用区ID  */
	AzName string `json:"azName,omitempty"` /*  可用区Name  */
}

type CcseListNodePoolsReturnObjRecordsNodesTaintsResponse struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略  */
}

type CcseListNodePoolsReturnObjRecordsNodesLabelsResponse struct{}
