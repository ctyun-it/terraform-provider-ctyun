package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetNodePoolApi
/* 调用该接口查看节点池详情。
 */type CcseGetNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetNodePoolApi(client *core.CtyunClient) *CcseGetNodePoolApi {
	return &CcseGetNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseGetNodePoolRequest) (*CcseGetNodePoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodePoolId", req.NodePoolId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetNodePoolRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseGetNodePoolResponse struct {
	StatusCode int32                             `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                            `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetNodePoolReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                            `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetNodePoolReturnObjResponse struct {
	Id                       string                                       `json:"id,omitempty"`                       /*  节点池id  */
	NodePoolName             string                                       `json:"nodePoolName,omitempty"`             /*  节点池名称  */
	BillMode                 string                                       `json:"billMode,omitempty"`                 /*  付费类型  */
	Description              string                                       `json:"description,omitempty"`              /*  描述  */
	NodeTotalNum             int32                                        `json:"nodeTotalNum,omitempty"`             /*  节点总数  */
	NormalNodeNum            int32                                        `json:"normalNodeNum,omitempty"`            /*  正常节点数  */
	UnNormalNodeNum          int32                                        `json:"unNormalNodeNum,omitempty"`          /*  异常节点数  */
	Nodes                    []*CcseGetNodePoolReturnObjNodesResponse     `json:"nodes"`                              /*  节点池节点信息  */
	Runtime                  string                                       `json:"runtime,omitempty"`                  /*  运行时  */
	RuntimeVersion           string                                       `json:"runtimeVersion,omitempty"`           /*  运行时版本  */
	NodeGroup                string                                       `json:"nodeGroup,omitempty"`                /*  伸缩组  */
	VmSpecName               string                                       `json:"vmSpecName,omitempty"`               /*  节点规格  */
	VmType                   string                                       `json:"vmType,omitempty"`                   /*  节点规格类型  */
	Cpu                      int32                                        `json:"cpu,omitempty"`                      /*  cpu  */
	Memory                   int32                                        `json:"memory,omitempty"`                   /*  内存  */
	DataDisks                []*CcseGetNodePoolReturnObjDataDisksResponse `json:"dataDisks"`                          /*  数据盘  */
	Status                   string                                       `json:"status,omitempty"`                   /*  状态  */
	MaxPodNum                int32                                        `json:"maxPodNum,omitempty"`                /*  最大pod数  */
	Labels                   *CcseGetNodePoolReturnObjLabelsResponse      `json:"labels"`                             /*  标签  */
	Taints                   []*CcseGetNodePoolReturnObjTaintsResponse    `json:"taints"`                             /*  污点  */
	IsDefault                *bool                                        `json:"isDefault"`                          /*  是否默认  */
	EnableAutoScale          *bool                                        `json:"enableAutoScale"`                    /*  是否自动弹性伸缩  */
	MaxNum                   int32                                        `json:"maxNum,omitempty"`                   /*  伸缩组最大数量  */
	MinNum                   int32                                        `json:"minNum,omitempty"`                   /*  伸缩组最小数量  */
	KubeletArgs              *CcseGetNodePoolReturnObjKubeletArgsResponse `json:"kubeletArgs"`                        /*  kubelet参数  */
	CreatedTime              string                                       `json:"createdTime,omitempty"`              /*  创建时间  */
	UpdatedTime              string                                       `json:"updatedTime,omitempty"`              /*  更新时间  */
	Vpc                      string                                       `json:"vpc,omitempty"`                      /*  vpc名称  */
	Subnetwork               string                                       `json:"subnetwork,omitempty"`               /*  子网  */
	SecruityGroupName        string                                       `json:"secruityGroupName,omitempty"`        /*  安全组名称  */
	ImageType                int32                                        `json:"imageType,omitempty"`                /*  操作系统镜像类型  */
	ImageName                string                                       `json:"imageName,omitempty"`                /*  操作系统镜像名称  */
	KeyName                  string                                       `json:"keyName,omitempty"`                  /*  ebm密钥对使用字段  */
	KeyPairId                string                                       `json:"keyPairId,omitempty"`                /*  ecs密钥对使用字段  */
	VisibilityPostHostScript string                                       `json:"visibilityPostHostScript,omitempty"` /*  部署后执行自定义脚本  */
	AzInfo                   []*CcseGetNodePoolReturnObjAzInfoResponse    `json:"azInfo"`                             /*  可用区  */
	SubnetUuid               string                                       `json:"subnetUuid,omitempty"`               /*  子网UUID  */
}

type CcseGetNodePoolReturnObjNodesResponse struct {
	ClusterId            string                                         `json:"clusterId,omitempty"`            /*  集群ID  */
	NodeName             string                                         `json:"nodeName,omitempty"`             /*  节点名称  */
	NodeType             int32                                          `json:"nodeType,omitempty"`             /*  节点类型，取值：<br />1：master <br />2：slave  */
	NodeStatus           string                                         `json:"nodeStatus,omitempty"`           /*  节点状态，取值：<br/>normal：健康。<br/>abnormal：异常。<br/>expulsion：驱逐中。  */
	IsSchedule           int32                                          `json:"isSchedule,omitempty"`           /*  是否调度，取值： 1：是。 <br />0：否。  */
	IsEvict              int32                                          `json:"isEvict,omitempty"`              /*  是否驱逐，取值： 1：是。 <br />0：否。  */
	DockerDataPath       string                                         `json:"dockerDataPath,omitempty"`       /*  docker数据目录  */
	CreatedTime          string                                         `json:"createdTime,omitempty"`          /*  创建时间。  */
	HostIp               string                                         `json:"hostIp,omitempty"`               /*  主机管理ip。  */
	HostIpv6             string                                         `json:"hostIpv6,omitempty"`             /*  主机管理ipv6。  */
	HostExtraIp          string                                         `json:"hostExtraIp,omitempty"`          /*  主机业务ip。  */
	HostExtraIpv6        string                                         `json:"hostExtraIpv6,omitempty"`        /*  主机业务ipv6。  */
	Cpu                  string                                         `json:"cpu,omitempty"`                  /*  cpu核数  */
	CpuUseRate           int32                                          `json:"cpuUseRate,omitempty"`           /*  cpu使用率%  */
	Memory               string                                         `json:"memory,omitempty"`               /*  内存  */
	MemoryUseRate        int32                                          `json:"memoryUseRate,omitempty"`        /*  内存使用率%  */
	Disk                 string                                         `json:"disk,omitempty"`                 /*  磁盘GiB  */
	DiskUseRate          int32                                          `json:"diskUseRate,omitempty"`          /*  磁盘使用率%  */
	KubeletVersion       string                                         `json:"kubeletVersion,omitempty"`       /*  Kubelet 版本  */
	PodCidr              string                                         `json:"podCidr,omitempty"`              /*  Pod CIDR  */
	KernelVersion        string                                         `json:"kernelVersion,omitempty"`        /*  内核版本  */
	OsImageVersion       string                                         `json:"osImageVersion,omitempty"`       /*  OS 镜像  */
	KubeProxyVersion     string                                         `json:"kubeProxyVersion,omitempty"`     /*  KubeProxy 版本  */
	DockerVersion        string                                         `json:"dockerVersion,omitempty"`        /*  容器版本  */
	IsSecure             int32                                          `json:"isSecure,omitempty"`             /*  是否安全节点，取值：<br />1：是。 <br />0：否。  */
	Taints               []*CcseGetNodePoolReturnObjNodesTaintsResponse `json:"taints"`                         /*  污点  */
	ChannelLabel         string                                         `json:"channelLabel,omitempty"`         /*  渠道标签  */
	ZoneName             string                                         `json:"zoneName,omitempty"`             /*  可用区名称  */
	CloudHostId          string                                         `json:"cloudHostId,omitempty"`          /*  云主机id  */
	Labels               *CcseGetNodePoolReturnObjNodesLabelsResponse   `json:"labels"`                         /*  标签  */
	EcsId                string                                         `json:"ecsId,omitempty"`                /*  paas平台云主机id  */
	VmSpecName           string                                         `json:"vmSpecName,omitempty"`           /*  节点规格名称  */
	Architecture         string                                         `json:"architecture,omitempty"`         /*  架构  */
	NeedShowHostRoom     string                                         `json:"needShowHostRoom,omitempty"`     /*  是否显示机房信息  */
	RoomName             string                                         `json:"roomName,omitempty"`             /*  机房名称  */
	FrameCode            string                                         `json:"frameCode,omitempty"`            /*  机柜编码  */
	PhysicalPosition     string                                         `json:"physicalPosition,omitempty"`     /*  物理位置  */
	RoomInfo             string                                         `json:"roomInfo,omitempty"`             /*  机房信息  */
	HostType             string                                         `json:"hostType,omitempty"`             /*  host类型  */
	ResourceId           string                                         `json:"resourceId,omitempty"`           /*  资源id  */
	NodeScaleDownProject string                                         `json:"nodeScaleDownProject,omitempty"` /*  开启/关闭缩容保护  */
	NodePoolName         string                                         `json:"nodePoolName,omitempty"`         /*  节点池名称  */
	LoginType            string                                         `json:"loginType,omitempty"`            /*  登陆类型  */
	VirtualNodeId        string                                         `json:"virtualNodeId,omitempty"`        /*  虚拟节点ID  */
}

type CcseGetNodePoolReturnObjDataDisksResponse struct {
	DiskSpecName string `json:"diskSpecName,omitempty"` /*  数据盘规格名称  */
	Size         int64  `json:"size,omitempty"`         /*  数据盘大小  */
}

type CcseGetNodePoolReturnObjLabelsResponse struct{}

type CcseGetNodePoolReturnObjTaintsResponse struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略，字典值：K8S_TAINTS_EFFECT  */
}

type CcseGetNodePoolReturnObjKubeletArgsResponse struct{}

type CcseGetNodePoolReturnObjAzInfoResponse struct {
	AzName string `json:"azName,omitempty"` /*  azName  */
}

type CcseGetNodePoolReturnObjNodesTaintsResponse struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略，字典值：K8S_TAINTS_EFFECT  */
}

type CcseGetNodePoolReturnObjNodesLabelsResponse struct{}
