package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseCreateNodePoolApi
/* 调用该接口新增节点池。
 */type CcseCreateNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseCreateNodePoolApi(client *core.CtyunClient) *CcseCreateNodePoolApi {
	return &CcseCreateNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseCreateNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseCreateNodePoolRequest) (*CcseCreateNodePoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseCreateNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseCreateNodePoolRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	NodePoolName             string `json:"nodePoolName,omitempty"`             /*  节点池名称  */
	Description              string `json:"description,omitempty"`              /*  节点池描述  */
	BillMode                 string `json:"billMode,omitempty"`                 /*  订单类型 1-包年包月 2-按需计费  */
	CycleCount               int32  `json:"cycleCount,omitempty"`               /*  订购时长，billMode为1必传，cycleType为MONTH时，cycleCount为1表示订购1个月  */
	CycleType                string `json:"cycleType,omitempty"`                /*  订购周期类型 MONTH-月 YEAR-年，billMode为1必传  */
	AutoRenewStatus          int32  `json:"autoRenewStatus,omitempty"`          /*  是否自动续订 0-否 1-是，默认为0  */
	VisibilityPostHostScript string `json:"visibilityPostHostScript,omitempty"` /*  部署后执行自定义脚本，base64编码  */
	VisibilityHostScript     string `json:"visibilityHostScript,omitempty"`     /*  部署前执行自定义脚本，base64编码  */
	DefinedHostnameEnable    int32  `json:"definedHostnameEnable,omitempty"`    /*  是否使用自定义节点名称；默认值0：不使用；1：使用  */
	HostNamePrefix           string `json:"hostNamePrefix,omitempty"`           /*  自定义主机名前缀，长度不超过10  */
	HostNamePostfix          string `json:"hostNamePostfix,omitempty"`          /*  自定义主机名后缀，长度不超过10  */
	ImageType                int32  `json:"imageType,omitempty"`                /*  镜像类型，0-私有，1-公有。  */
	ImageName                string `json:"imageName,omitempty"`                /*  镜像名称  */
	ImageUuid                string `json:"imageUuid,omitempty"`                /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004475" target="_blank">节点规格和节点镜像</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&amp;api=4765&amp;data=89" target="_blank">创建私有镜像（云主机系统盘）</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&amp;api=5230&amp;data=89" target="_blank">创建私有镜像（云主机数据盘）</a>
	注：同一镜像名称在不同资源池的镜像ID是不同的，调用前需确认镜像ID是否归属当前资源池  */
	LoginType string `json:"loginType,omitempty"` /*  云主机密码登录类型：secretPair：密钥对；password：密码  */
	EcsPasswd string `json:"ecsPasswd,omitempty"` /*  用户密码，如果loginType=password，这该项为必填项，满足以下规则：
	长度在8～30个字符
	必须包含大写字母、小写字母、数字以及特殊符号中的三项
	特殊符号可选：()`-!@#$%^&*_-+=｜{}[]:;'<>,.?/且不能以斜线号 / 开头
	不能包含3个及以上连续字符
	Linux镜像不能包含镜像用户名（root）、用户名的倒序（toor）、用户名大小写变化（如RoOt、rOot等）
	Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator等）
	*/
	KeyName string `json:"keyName,omitempty"` /*  密钥对名称，如果loginType=secretPair，这该项为必填项，您可以查看密钥对来了解密钥对相关内容
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87&vid=81">查询一个或多个密钥对</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8344&data=87&vid=81" target="_blank">创建一对SSH密钥对</a>  */
	KeyPairId string `json:"keyPairId,omitempty"` /*  密钥对ID，如果loginType=secretPair，这该项为必填项。您可以查看密钥对来了解密钥对相关内容
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87&vid=81">查询一个或多个密钥对</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8344&data=87&vid=81" target="_blank">创建一对SSH密钥对</a>  */
	UseAffinityGroup  *bool  `json:"useAffinityGroup"`            /*  是否启用主机组  */
	AffinityGroupUuid string `json:"affinityGroupUuid,omitempty"` /*  云主机组ID，获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8324&data=87&vid=81" target="_blank">查询云主机组列表或者详情</a>
	<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8316&data=87&vid=81" target="_blank">创建云主机组</a>  */
	ResourceLabels    *CcseCreateNodePoolResourceLabelsRequest `json:"resourceLabels"`       /*  云主机资源标签  */
	SyncNodeLabels    *bool                                    `json:"syncNodeLabels"`       /*  是否同步节点标签  */
	SyncNodeTaints    *bool                                    `json:"syncNodeTaints"`       /*  是否同步节点污点  */
	NodeUnschedulable *bool                                    `json:"nodeUnschedulable"`    /*  节点是否不可调度  */
	Labels            *CcseCreateNodePoolLabelsRequest         `json:"labels"`               /*  标签  */
	Taints            []*CcseCreateNodePoolTaintsRequest       `json:"taints"`               /*  节点污点，格式为 [{\"key\":\"{key}\",\"value\":\"{value}\",\"effect\":\"{effect}\"}]，上述的{key}、{value}、{effect}替换为所需字段。effect枚举包括NoSchedule、PreferNoSchedule、NoExecute  */
	VmSpecName        string                                   `json:"vmSpecName,omitempty"` /*  节点规格，获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87&isNormal=1&vid=81" target="_blank">查询主机规格资源</a>  */
	VmType string `json:"vmType,omitempty"` /*  节点规格类型，获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87&isNormal=1&vid=81" target="_blank">查询主机规格资源</a>  */
	Cpu             int32                                 `json:"cpu,omitempty"`         /*  CPU大小  */
	Memory          int32                                 `json:"memory,omitempty"`      /*  内存大小,单位是G  */
	MaxNum          int32                                 `json:"maxNum,omitempty"`      /*  伸缩组最大数量  */
	MinNum          int32                                 `json:"minNum,omitempty"`      /*  伸缩组最小数量  */
	EnableAutoScale *bool                                 `json:"enableAutoScale"`       /*  是否自动弹性伸缩  */
	DataDisks       []*CcseCreateNodePoolDataDisksRequest `json:"dataDisks"`             /*  数据盘  */
	MaxPodNum       int32                                 `json:"maxPodNum,omitempty"`   /*  最大pod数，默认110  */
	Gpu             int32                                 `json:"gpu,omitempty"`         /*  gpu大小  */
	SysDiskType     string                                `json:"sysDiskType,omitempty"` /*  系统盘规格，云硬盘类型，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	SysDiskSize int32                              `json:"sysDiskSize,omitempty"` /*  系统盘大小，单位是G  */
	AzInfo      []*CcseCreateNodePoolAzInfoRequest `json:"azInfo"`                /*  可用区信息，包括可用区ID，可用区名称
	可用区名称获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&vid=81" target="_blank">资源池可用区查询</a>  */
}

type CcseCreateNodePoolResourceLabelsRequest struct{}

type CcseCreateNodePoolLabelsRequest struct{}

type CcseCreateNodePoolTaintsRequest struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略  */
}

type CcseCreateNodePoolDataDisksRequest struct {
	Size         int32  `json:"size,omitempty"`         /*  数据盘大小，单位G  */
	DiskSpecName string `json:"diskSpecName,omitempty"` /*  数据盘规格名称，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
}

type CcseCreateNodePoolAzInfoRequest struct {
	AzName string `json:"azName,omitempty"` /*  可用区名称  */
}

type CcseCreateNodePoolResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  状态码  */
	Message    string `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码  */
}
