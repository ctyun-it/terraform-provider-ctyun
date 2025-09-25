package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseAttachClusterNodesApi
/* 调用该接口纳管节点。
 */
type CcseAttachClusterNodesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseAttachClusterNodesApi(client *core.CtyunClient) *CcseAttachClusterNodesApi {
	return &CcseAttachClusterNodesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodes/attach",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseAttachClusterNodesApi) Do(ctx context.Context, credential core.Credential, req *CcseAttachClusterNodesRequest) (*CcseAttachClusterNodesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CcseAttachClusterNodesRequest
		RegionId  interface{} `json:"regionId,omitempty"`
		ClusterId interface{} `json:"clusterId,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseAttachClusterNodesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseAttachClusterNodesRequest struct {
	ClusterId string `json:"clusterId,omitempty"` /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId  string `json:"regionId,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Instances []*CcseAttachClusterNodesInstancesRequest `json:"instances"`           /*  云主机ECS 信息  */
	VmType    string                                    `json:"vmType,omitempty"`    /*  默认ecs，值为：弹性云主机则是：ecs ；物理机则是：ebm  */
	Runtime   string                                    `json:"runtime,omitempty"`   /*  容器运行时，目前仅支持 CONTAINERD  */
	ImageType int32                                     `json:"imageType,omitempty"` /*  镜像类型，0-私有镜像，1-公有镜像。
	您可查看<a href="https://www.ctyun.cn/document/10026730/10030151" target="_blank">镜像概述</a>  */
	ImageUuid string `json:"imageUuid,omitempty"` /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004475" target="_blank">节点规格和节点镜像</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&amp;api=4765&amp;data=89" target="_blank">创建私有镜像（云主机系统盘）</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&amp;api=5230&amp;data=89" target="_blank">创建私有镜像（云主机数据盘）</a>
	注：同一镜像名称在不同资源池的镜像ID是不同的，调用前需确认镜像ID是否归属当前资源池  */
	LoginType string `json:"loginType,omitempty"` /*  主机登录方式： password：密码，secretPair：密钥对  */
	Password  string `json:"password,omitempty"`  /*  云主机或物理机，用户登录密码，如果loginType=password，这该项为必填项，满足以下规则：
	物理机：用户密码，满足以下规则：长度在8～30个字符；
	必须包含大小写字母和（至少一个数字或者特殊字符）；
	特殊符号可选：()`~!@#$%&*_-+=\
	云主机：长度在8～30个字符；
	必须包含大写字母、小写字母、数字以及特殊符号中的三项；
	特殊符号可选：()`-!@#$%^&*_-+=｜{}[]:;'<>,.?/且不能以斜线号 / 开头；
	不能包含3个及以上连续字符；
	纳管节点时password字段可选择加密，具体加密方法请参见<a href="https://www.ctyun.cn/document/10083472/11002096" target="_blank">password字段加密的方法</a>  */
	KeyName string `json:"keyName,omitempty"` /*  密钥对名称，如果loginType=secretPair，物理机该项为必填项，您可以查看密钥对来了解密钥对相关内容
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87&vid=81" target="_blank">查询一个或多个密钥对</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8344&data=87&vid=81" target="_blank">创建一对SSH密钥对</a>  */
	KeyPairId string `json:"keyPairId,omitempty"` /*  密钥对ID，如果loginType=secretPair，弹性云主机该项为必填项。您可以查看密钥对来了解密钥对相关内容
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87&vid=81">查询一个或多个密钥对</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8344&data=87&vid=81" target="_blank">创建一对SSH密钥对</a>  */
	Labels                      *CcseAttachClusterNodesLabelsRequest         `json:"labels"`                             /*  K8S节点标签  */
	Annotations                 *CcseAttachClusterNodesAnnotationsRequest    `json:"annotations"`                        /*  K8S节点注解  */
	Taints                      []*CcseAttachClusterNodesTaintsRequest       `json:"taints"`                             /*  节点污点，格式为 [{\"key\":\"{key}\",\"value\":\"{value}\",\"effect\":\"{effect}\"}]，上述的{key}、{value}、{effect}替换为所需字段。effect枚举包括NoSchedule、PreferNoSchedule、NoExecute  */
	VisibilityPostHostScript    string                                       `json:"visibilityPostHostScript,omitempty"` /*  部署后执行自定义脚本 （输入的值需要经过Base64编码，方法如下：echo -n "待编码内容" \  */
	VisibilityHostScript        string                                       `json:"visibilityHostScript,omitempty"`     /*  部署前执行自定义脚本（输入的值需要经过Base64编码，方法如下：echo -n "待编码内容" \  */
	KubeletArgs                 *CcseAttachClusterNodesKubeletArgsRequest    `json:"kubeletArgs"`                        /*  Kubelet自定义参数  */
	IsSyncClusterResourceLabels *bool                                        `json:"isSyncClusterResourceLabels"`        /*  是否同步集群标签。默认是false。如果为true，则同步以当前集群标签为基准的快照。  */
	ResourceLabels              *CcseAttachClusterNodesResourceLabelsRequest `json:"resourceLabels"`                     /*  云主机资源标签  */
	KubeletDirectory            string                                       `json:"kubeletDirectory,omitempty"`         /*  kubelet数据目录。该参数可以自定义指定在/data下的子目录，用于kubelet数据目录。该参数由白名单控制  */
	ContainerDataDirectory      string                                       `json:"containerDataDirectory,omitempty"`   /*  容器数据目录。该参数可以自定义指定在/data下的子目录，用于容器数据目录。该参数由白名单控制  */
}

type CcseAttachClusterNodesInstancesRequest struct {
	InstanceId string `json:"instanceId,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs" target="_blank">弹性云主机</a>了解云主机的相关信息
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=8309&amp;data=87">查询云主机列表</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=8281&amp;data=87" target="_blank">创建一台按量付费或包年包月的云主机</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=8282&amp;data=87">批量创建按量付费或包年包月云主机</a>

	物理机 instanceUUID，获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6941&data=97&isNormal=1&vid=91">批量查询物理机</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=16&api=6942&data=97&isNormal=1&vid=91">物理机创建</a>  */
	AzName string `json:"azName,omitempty"` /*  可用区名称，纳管是物理机，此项必填，可用区信息可用区名称获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&vid=81" target="_blank">资源池可用区查询</a>  */
}

type CcseAttachClusterNodesLabelsRequest struct{}

type CcseAttachClusterNodesAnnotationsRequest struct{}

type CcseAttachClusterNodesTaintsRequest struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略  */
}

type CcseAttachClusterNodesKubeletArgsRequest struct {
	KubeAPIQPS           int32  `json:"kubeAPIQPS,omitempty"`           /*  请求至kube-apiserver的QPS配置  */
	KubeAPIBurst         int32  `json:"kubeAPIBurst,omitempty"`         /*  请求至kube-apiserver的Burst配置  */
	MaxPods              int32  `json:"maxPods,omitempty"`              /*  kubelet管理的pod上限  */
	RegistryPullQPS      int32  `json:"registryPullQPS,omitempty"`      /*  每秒钟可以执行的镜像仓库拉取操作限值  */
	RegistryBurst        int32  `json:"registryBurst,omitempty"`        /*  突发性镜像拉取的上限值  */
	PodPidsLimit         int32  `json:"podPidsLimit,omitempty"`         /*  限制Pod中的进程数  */
	EventRecordQPS       int32  `json:"eventRecordQPS,omitempty"`       /*  事件创建QPS限制  */
	EventBurst           int32  `json:"eventBurst,omitempty"`           /*  事件创建的个数的突发峰值上限  */
	TopologyManagerScope string `json:"topologyManagerScope,omitempty"` /*   拓扑管理策略的资源对齐粒度  */
	CpuCFSQuota          *bool  `json:"cpuCFSQuota"`                    /*  是否为设置了CPU限制的容器实施CPU CFS配额约束,默认值为true  */
}

type CcseAttachClusterNodesResourceLabelsRequest struct{}

type CcseAttachClusterNodesResponse struct {
	StatusCode int32  `json:"statusCode"` /*  响应状态码  */
	RequestId  string `json:"requestId"`  /*  请求ID  */
	Message    string `json:"message"`    /*  响应信息  */
	ReturnObj  *bool  `json:"returnObj"`  /*  响应对象  */
	Error      string `json:"error"`      /*  错误码，参见错误码说明  */
}
