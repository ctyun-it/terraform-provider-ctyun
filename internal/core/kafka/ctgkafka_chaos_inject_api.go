package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaChaosInjectApi
/* 注入故障演练任务。
 */type CtgkafkaChaosInjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaChaosInjectApi(client *core.CtyunClient) *CtgkafkaChaosInjectApi {
	return &CtgkafkaChaosInjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/chaos/inject",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaChaosInjectApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaChaosInjectRequest) (*CtgkafkaChaosInjectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaChaosInjectRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaChaosInjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaChaosInjectRequest struct {
	RegionId        string                                     `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId      string                                     `json:"prodInstId,omitempty"` /*  实例ID。  */
	ActionCode      string                                     `json:"actionCode,omitempty"` /*  故障任务类型。<br><li>node-shutdown：节点关机，Broker宕机，请勿频繁进行开关机演练，可能会因Kafka监控数据不全导致演练失败(若出现失败，可在15分钟之后重试)。<br><li>cpu-fullload：CPU高负载，Broker CPU高负载，该故障动作会随机向Kafka实例的一个Broker节点注入CPU高负载故障，建议实例上有真实消息生产与消费流量，主题分区数大于3。<br><li>disk-burn：磁盘IO高负载，Broker 磁盘I0高负载，该故障动作会随机向Kafka实例的一个Broker节点注入磁盘I0高负载故障，建议实例上有真实消息生产与消费流量，主题分区数大于3。  */
	ActionParameter *CtgkafkaChaosInjectActionParameterRequest `json:"actionParameter"`      /*  故障注入参数。  */
}

type CtgkafkaChaosInjectActionParameterRequest struct {
	CpuPercent   int32  `json:"cpuPercent,omitempty"`   /*  演练CPU负载率，取值范围[1, 100]，当actionCode=cpu-fullload时需传入。  */
	Duration     int32  `json:"duration,omitempty"`     /*  故障持续时间，单位秒，取值范围[60, 3600]，当actionCode=disk-burn或cpu-fullload时需传入。  */
	NodeKillType int32  `json:"nodeKillType,omitempty"` /*  节点宕机模式，当actionCode=node-shutdown时需传入。<br><li>0：随机节点宕机，随机宕机1个节点<br><li>1：随机可用区宕机，宕机1个随机可用区下所有节点<br><li>2：指定节点宕机<br><li>3：指定可用区宕机，宕机该可用区下所有节点</li>  */
	EcsId        string `json:"ecsId,omitempty"`        /*  指定节点宕机的节点ID，当nodeKillType=2时需传入。<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382&data=83&isNormal=1&vid=330">实例信息</a> 查询对应实例下nodeList的ecsId属性值。  */
	AzName       string `json:"azName,omitempty"`       /*  指定可用区宕机的可用区名称，当nodeKillType=3时需传入。<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382&data=83&isNormal=1&vid=33">实例信息</a> 查询实例下nodeList的azName属性值获取实例所在可用区列表。  */
}

type CtgkafkaChaosInjectResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaChaosInjectReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaChaosInjectReturnObjResponse struct {
	ExperimentId string `json:"experimentId,omitempty"` /*  故障演练ID。  */
	TaskId       string `json:"taskId,omitempty"`       /*  故障演练任务ID。  */
}
