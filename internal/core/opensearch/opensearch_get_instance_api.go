package opensearch

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OpensearchGetInstanceApi
/* 查询 OpenSearch 实例详情<br />
<b>准备工作：</b><br />
&emsp;&emsp;构造请求：在调用前需要了解如何构造请求<br />
&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用<br />
*/type OpensearchGetInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOpensearchGetInstanceApi(client *core.CtyunClient) *OpensearchGetInstanceApi {
	return &OpensearchGetInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/os/openapi/v1/cluster/getClusterById",
			ContentType:  "application/json",
		},
	}
}

func (a *OpensearchGetInstanceApi) Do(ctx context.Context, credential core.Credential, req *GetInstanceRequest) (*GetInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	// GET 请求，参数通过 Query 传递
	ctReq.AddParam("clusterId", req.ClusterID)

	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp GetInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetInstanceRequest struct {
	ClusterID string `json:"-"` /*  实例 id（Query参数）*/
}

type GetInstanceResponse struct {
	StatusCode int32          `json:"statusCode"` /*  状态码，成功：200，失败：500  */
	Error      string         `json:"error"`      /*  错误码，请求成功时，不返回该字段  */
	Message    string         `json:"message"`    /*  用来简述当前接口调用状态以及必要提示信息  */
	ReturnObj  InstanceDetail `json:"returnObj"`  /*  返回结果  */
}

type InstanceDetail struct {
	ClusterID                string          `json:"clusterId"`                /*  实例 id  */
	ClusterName              string          `json:"clusterName"`              /*  实例名称  */
	State                    string          `json:"state"`                    /*  健康状态：GREEN/YELLOW/RED  */
	EnableIpv6               string          `json:"enableIpv6"`               /*  OPEN/CLOSE/NOT_DISPLAY  */
	RegionID                 string          `json:"regionId"`                 /*  资源池编码  */
	RegionName               string          `json:"regionName"`               /*  资源池名称  */
	AvailableZoneID          string          `json:"availableZoneId"`          /*  可用区编码  */
	AzName                   string          `json:"azName"`                   /*  可用区名称  */
	VPCName                  string          `json:"vpcName"`                  /*  vpc 名称  */
	VPCID                    string          `json:"vpcId"`                    /*  vpcId  */
	SubnetName               string          `json:"subnetName"`               /*  子网名称  */
	SubnetID                 string          `json:"subnetId"`                 /*  子网 id  */
	SecurityGroupID          string          `json:"securityGroupId"`          /*  安全组 id  */
	SecurityGroupName        string          `json:"securityGroupnName"`       /*  安全组名称  */
	CPUInfo                  string          `json:"cpuInfo"`                  /*  cpu 架构  */
	OSType                   string          `json:"osType"`                   /*  操作系统类型  */
	ClusterType              int32           `json:"clusterType"`              /*  实例类型：1:OpenSearch, 2:Elasticsearch  */
	ClusterTypeName          string          `json:"clusterTypeName"`          /*  类型名称  */
	ClusterTypeVersion       string          `json:"clusterTypeVersion"`       /*  实例版本  */
	PayType                  string          `json:"payType"`                  /*  付费类型  */
	ClusterDueTime           int64           `json:"clusterDueTime"`           /*  实例到期时间  */
	CreateTime               int64           `json:"createTime"`               /*  创建时间  */
	UserName                 string          `json:"userName"`                 /*  访问控制 - 用户名  */
	OsVmSpecName             string          `json:"osVmSpecName"`             /*  映射的主机名称（节点规格名称）*/
	ClusterMessage           *string         `json:"clusterMessage"`           /*  错误原因  */
	CPUNum                   int32           `json:"cpuNum"`                   /*  cpu 大小  */
	Memory                   int32           `json:"memory"`                   /*  内存大小  */
	HostNum                  int32           `json:"hostNum"`                  /*  主机数量  */
	DiskVolumn               int32           `json:"diskVolumn"`               /*  存储空间  */
	ComponentName            string          `json:"componentName"`            /*  组件名称  */
	LoadBalancerName         *string         `json:"loadBalancerName"`         /*  负载均衡器名称  */
	TargetGroupName          *string         `json:"targetGroupName"`          /*  后端主机组名称  */
	ProjectID                string          `json:"projectId"`                /*  企业项目编码  */
	ProjectName              string          `json:"projectName"`              /*  企业名称  */
	RouterHostInfo           *RouterHostInfo `json:"routerHostInfo"`           /*  组件节点信息  */
	DataHostInfos            []DataHostInfo  `json:"dataHostInfos"`            /*  数据节点组类型信息  */
	ExclusiveMasterHostInfos []DataHostInfo  `json:"exclusiveMasterHostInfos"` /*  专属 MASTER 数据节点组类型信息  */
	CoordinateHostInfos      []DataHostInfo  `json:"coordinateHostInfos"`      /*  协调节点组类型信息  */
	ColdHostInfos            []DataHostInfo  `json:"coldHostInfos"`            /*  冷数据节点组类型信息  */
	LogstashHostInfos        []DataHostInfo  `json:"logstashHostInfos"`        /*  Logstash 节点信息  */
}

type RouterHostInfo struct {
	HostIP          string `json:"hostIp"`          /*  内网 ip  */
	AvailableZoneID string `json:"availableZoneId"` /*  区域 id  */
	State           string `json:"state"`           /*  RUNNING、FAILED  */
	StateType       int32  `json:"stateType"`       /*  2:运行中 5:失败  */
	IPv6HostIP      string `json:"ipv6HostIp"`      /*  ipv6 地址  */
	IaasVMSpecCode  string `json:"iaasVmSpecCode"`  /*  云主机规格编码  */
	IOTypeName      string `json:"ioTypeName"`      /*  磁盘 IO 类型名称  */
	CPUNum          int32  `json:"cpuNum"`          /*  cpu 核数  */
	Memory          int32  `json:"memory"`          /*  内存大小  */
	DiskVolumn      int32  `json:"diskVolumn"`      /*  硬盘大小  */
	IaasVmTypeName  string `json:"iaasVmTypeName"`  /*  云主机类型  */
}

type DataHostInfo struct {
	HostIP          string `json:"hostIp"`          /*  内网 ip  */
	AvailableZoneID string `json:"availableZoneId"` /*  区域 id  */
	State           string `json:"state"`           /*  RUNNING、FAILED  */
	StateType       int32  `json:"stateType"`       /*  2:运行中 5:FAILED  */
	IPv6HostIP      string `json:"ipv6HostIp"`      /*  ipv6 地址  */
	IaasVMSpecCode  string `json:"iaasVmSpecCode"`  /*  IAAS 虚机规格编码  */
	IOTypeName      string `json:"ioTypeName"`      /*  磁盘 IO 类型名称  */
	CPUNum          int32  `json:"cpuNum"`          /*  cpu 核数  */
	Memory          int32  `json:"memory"`          /*  内存大小  */
	DiskVolumn      int32  `json:"diskVolumn"`      /*  硬盘大小  */
	IaasVmTypeName  string `json:"iaasVmTypeName"`  /*  云主机类型  */
}
