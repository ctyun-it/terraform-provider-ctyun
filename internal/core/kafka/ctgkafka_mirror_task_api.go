package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaMirrorTaskApi
/* 跨集群迁移，用于自建实例或跨云云实例与分布式消息Kafka实例之间的数据同步。
 */type CtgkafkaMirrorTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaMirrorTaskApi(client *core.CtyunClient) *CtgkafkaMirrorTaskApi {
	return &CtgkafkaMirrorTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/mirrorTask",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaMirrorTaskApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaMirrorTaskRequest) (*CtgkafkaMirrorTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaMirrorTaskRequest
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
	var resp CtgkafkaMirrorTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaMirrorTaskRequest struct {
	RegionId            string `json:"regionId,omitempty"`            /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId          string `json:"prodInstId,omitempty"`          /*  目标实例ID。  */
	TaskName            string `json:"taskName,omitempty"`            /*  迁移任务名称。  */
	SourceAddr          string `json:"sourceAddr,omitempty"`          /*  源集群的接入点地址，迁移上云时为必填。  */
	SourceProtocol      string `json:"sourceProtocol,omitempty"`      /*  源集群的接入协议，迁移上云时为必填。<br><li>PLAINTEXT<br><li>SASL_PLAINTEXT  */
	SourceSaslMechanism string `json:"sourceSaslMechanism,omitempty"` /*  sasl认证机制，sasl连接时为必填。<br><li>PLAIN<br><li>SCRAM-SHA-256<br><li>SCRAM-SHA-512  */
	SourceSaslUser      string `json:"sourceSaslUser,omitempty"`      /*  sasl用户名称，sasl连接时为必填。  */
	SourceSaslPwd       string `json:"sourceSaslPwd,omitempty"`       /*  sasl用户密码，sasl连接时为必填。  */
	TaskNum             int32  `json:"taskNum,omitempty"`             /*  任务数，默认为1。  */
	SyncAcl             string `json:"syncAcl,omitempty"`             /*  是否同步acl信息。<br><li>1：是<br><li>2：否（默认值）  */
	SyncGroup           string `json:"syncGroup,omitempty"`           /*  是否同步消费组。<br><li>1：是<br><li>2：否（默认值）  */
	Topics              string `json:"topics,omitempty"`              /*  迁移的Topic名称，多个Topic用逗号隔开，不填默认所有Topic。  */
	Groups              string `json:"groups,omitempty"`              /*  迁移的消费组名称，多个消费组用逗号隔开，不填默认所有消费组。  */
	AutoStopTask        string `json:"autoStopTask,omitempty"`        /*  是否自动停止任务。<br><li>1：是<br><li>2：否（默认值）  */
	RawType             string `json:"type,omitempty"`                /*  迁移类型。<br><li>2：本地实例迁移上云<br><li>3：云间实例迁移  */
	SourceClusterId     string `json:"sourceClusterId,omitempty"`     /*  源集群ID，云实例间迁移时必填。  */
	DefaultReplica      string `json:"defaultReplica,omitempty"`      /*  是否保持和源集群副本一致。<br><li>1：是（默认值）<br><li>2：否，3副本  */
}

type CtgkafkaMirrorTaskResponse struct {
	StatusCode string                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                               `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaMirrorTaskReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaMirrorTaskReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回提升。  */
}
