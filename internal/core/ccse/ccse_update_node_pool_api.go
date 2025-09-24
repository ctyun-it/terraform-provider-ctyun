package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateNodePoolApi
/* 调用该接口更新节点池。
 */type CcseUpdateNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateNodePoolApi(client *core.CtyunClient) *CcseUpdateNodePoolApi {
	return &CcseUpdateNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateNodePoolRequest) (*CcseUpdateNodePoolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodePoolId", req.NodePoolId)
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
	var resp CcseUpdateNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateNodePoolRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	NodePoolName    string                                `json:"nodePoolName,omitempty"`    /*  节点池名称  */
	BillMode        string                                `json:"billMode,omitempty"`        /*  订单类型 1-包年包月 2-按需计费  */
	CycleCount      int32                                 `json:"cycleCount,omitempty"`      /*  订购时长，billMode为1必传，cycleType为MONTH时，cycleCount为1表示订购1个月  */
	CycleType       string                                `json:"cycleType,omitempty"`       /*  订购周期类型 MONTH-月 YEAR-年，billMode为1必传  */
	AutoRenewStatus int32                                 `json:"autoRenewStatus,omitempty"` /*  是否自动续订 0-否 1-是，默认为0  */
	Description     string                                `json:"description,omitempty"`     /*  描述  */
	DataDisks       []*CcseUpdateNodePoolDataDisksRequest `json:"dataDisks"`                 /*  数据盘  */
	Labels          *CcseUpdateNodePoolLabelsRequest      `json:"labels"`                    /*  标签  */
	Taints          []*CcseUpdateNodePoolTaintsRequest    `json:"taints"`                    /*  节点污点，格式为 [{\"key\":\"{key}\",\"value\":\"{value}\",\"effect\":\"{effect}\"}]，上述的{key}、{value}、{effect}替换为所需字段。effect枚举包括NoSchedule、PreferNoSchedule、NoExecute  */
	EnableAutoScale *bool                                 `json:"enableAutoScale"`           /*  是否自动弹性伸缩  */
	MaxNum          int32                                 `json:"maxNum,omitempty"`          /*  伸缩组最大数量0-20  */
	MinNum          int32                                 `json:"minNum,omitempty"`          /*  伸缩组最小数量0-20  */
	SysDiskSize     int32                                 `json:"sysDiskSize,omitempty"`     /*  系统盘大小  */
	SysDiskType     string                                `json:"sysDiskType,omitempty"`     /*  系统盘规格，云硬盘类型，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	VisibilityPostHostScript string `json:"visibilityPostHostScript,omitempty"` /*  部署后执行自定义脚本，base64编码  */
	VisibilityHostScript     string `json:"visibilityHostScript,omitempty"`     /*  部署前执行自定义脚本，base64编码  */
}

type CcseUpdateNodePoolDataDisksRequest struct {
	Size         int32  `json:"size,omitempty"`         /*  数据盘大小，单位G  */
	DiskSpecName string `json:"diskSpecName,omitempty"` /*  数据盘规格名称，取值范围：
	SATA：普通IO，
	SAS：高IO，
	SSD：超高IO
	您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
}

type CcseUpdateNodePoolLabelsRequest struct{}

type CcseUpdateNodePoolTaintsRequest struct {
	Key    string `json:"key,omitempty"`    /*  键  */
	Value  string `json:"value,omitempty"`  /*  值  */
	Effect string `json:"effect,omitempty"` /*  策略  */
}

type CcseUpdateNodePoolResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求ID  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *bool  `json:"returnObj"`            /*  返回对象  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
