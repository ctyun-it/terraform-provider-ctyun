package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaChaosListApi
/* 分页查询故障演练列表。
 */type CtgkafkaChaosListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaChaosListApi(client *core.CtyunClient) *CtgkafkaChaosListApi {
	return &CtgkafkaChaosListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/chaos/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaChaosListApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaChaosListRequest) (*CtgkafkaChaosListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(req.PageNum), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaChaosListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaChaosListRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	PageNum    int32  `json:"pageNum,omitempty"`    /*  分页中的页数，默认1，范围1-40000。  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  分页中的每页大小，默认10。  */
}

type CtgkafkaChaosListResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                              `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaChaosListReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaChaosListReturnObjResponse struct {
	Total int32                                     `json:"total,omitempty"` /*  总记录数。  */
	Data  []*CtgkafkaChaosListReturnObjDataResponse `json:"data"`            /*  故障演练信息。  */
}

type CtgkafkaChaosListReturnObjDataResponse struct {
	ExperimentId           string                                                        `json:"experimentId,omitempty"` /*  故障演练ID，故障注入与故障恢复一起构成完整的一次故障演练。使用该ID查询故障执行详情与执行撤销故障。  */
	Action                 *CtgkafkaChaosListReturnObjDataActionResponse                 `json:"action"`                 /*  故障演练动作信息。  */
	CreateDate             int64                                                         `json:"createDate,omitempty"`   /*  演练创建毫秒时间戳。  */
	InjectActionInstance   *CtgkafkaChaosListReturnObjDataInjectActionInstanceResponse   `json:"injectActionInstance"`   /*  故障注入任务。  */
	RecoveryActionInstance *CtgkafkaChaosListReturnObjDataRecoveryActionInstanceResponse `json:"recoveryActionInstance"` /*  故障撤销任务。  */
}

type CtgkafkaChaosListReturnObjDataActionResponse struct {
	ActionCode      string                                                         `json:"actionCode,omitempty"`  /*  故障任务类型。<br><li>node-shutdown：节点关机<br><li>cpu-fullload：CPU高负载<br><li>disk-burn：磁盘IO高负载  */
	ActionScope     string                                                         `json:"actionScope,omitempty"` /*  故障范围。<br><li>host：主机<br><li>cluster：集群  */
	ActionParameter []*CtgkafkaChaosListReturnObjDataActionActionParameterResponse `json:"actionParameter"`       /*  任务参数配置。  */
}

type CtgkafkaChaosListReturnObjDataInjectActionInstanceResponse struct {
	TaskId          string                                                                       `json:"taskId,omitempty"`    /*  任务ID。  */
	State           int32                                                                        `json:"state,omitempty"`     /*  动作执行状态。<br><li>0：任务完成<br><li>1：任务失败<br><li>2：任务进行中<br><li>3：任务未开始  */
	StartTime       int64                                                                        `json:"startTime,omitempty"` /*  开始时间毫秒时间戳。  */
	EndTime         int64                                                                        `json:"endTime,omitempty"`   /*  结束时间毫秒时间戳，当actionCode=disk-burn或cpu-fullload时有效。  */
	ActionParameter []*CtgkafkaChaosListReturnObjDataInjectActionInstanceActionParameterResponse `json:"actionParameter"`     /*  动作执行参数。  */
	Logs            []*CtgkafkaChaosListReturnObjDataInjectActionInstanceLogsResponse            `json:"logs"`                /*  演练任务过程日志。  */
}

type CtgkafkaChaosListReturnObjDataRecoveryActionInstanceResponse struct {
	TaskId          string                                                                         `json:"taskId,omitempty"`    /*  任务ID。  */
	State           int32                                                                          `json:"state,omitempty"`     /*  动作执行状态。<br><li>0：任务完成<br><li>1：任务失败<br><li>2：任务进行中<br><li>3：任务未开始  */
	StartTime       int64                                                                          `json:"startTime,omitempty"` /*  开始时间毫秒时间戳。  */
	EndTime         int64                                                                          `json:"endTime,omitempty"`   /*  结束时间毫秒时间戳，当actionCode=disk-burn或cpu-fullload时有效。  */
	ActionParameter []*CtgkafkaChaosListReturnObjDataRecoveryActionInstanceActionParameterResponse `json:"actionParameter"`     /*  动作执行参数。  */
	Logs            []*CtgkafkaChaosListReturnObjDataRecoveryActionInstanceLogsResponse            `json:"logs"`                /*  演练任务过程日志。  */
}

type CtgkafkaChaosListReturnObjDataActionActionParameterResponse struct {
	Code  string `json:"code,omitempty"`  /*  属性编码。  */
	Value string `json:"value,omitempty"` /*  属性值。  */
}

type CtgkafkaChaosListReturnObjDataInjectActionInstanceActionParameterResponse struct {
	Code  string `json:"code,omitempty"`  /*  属性编码。  */
	Value string `json:"value,omitempty"` /*  属性值。  */
}

type CtgkafkaChaosListReturnObjDataInjectActionInstanceLogsResponse struct {
	StepName    string `json:"stepName,omitempty"`    /*  步骤名称。  */
	SubStepName string `json:"subStepName,omitempty"` /*  子步骤名称。  */
	LogLevel    string `json:"logLevel,omitempty"`    /*  日志级别。<br><li>INFO<br><li>ERROR  */
	StepResult  string `json:"stepResult,omitempty"`  /*  步骤结果。<br><li>PROCESSING：处理中<br><li>SUCCESS：成功<br><li>ERROR：错误<br><li>RETRYING：重试中  */
	Msg         string `json:"msg,omitempty"`         /*  日志内容。  */
	LogTime     string `json:"logTime,omitempty"`     /*  日志时间。  */
}

type CtgkafkaChaosListReturnObjDataRecoveryActionInstanceActionParameterResponse struct {
	Code  string `json:"code,omitempty"`  /*  属性编码。  */
	Value string `json:"value,omitempty"` /*  属性值。  */
}

type CtgkafkaChaosListReturnObjDataRecoveryActionInstanceLogsResponse struct {
	StepName    string `json:"stepName,omitempty"`    /*  步骤名称。  */
	SubStepName string `json:"subStepName,omitempty"` /*  子步骤名称。  */
	LogLevel    string `json:"logLevel,omitempty"`    /*  日志级别。<br><li>INFO<br><li>ERROR  */
	StepResult  string `json:"stepResult,omitempty"`  /*  步骤结果。<br><li>PROCESSING：处理中<br><li>SUCCESS：成功<br><li>ERROR：错误<br><li>RETRYING：重试中  */
	Msg         string `json:"msg,omitempty"`         /*  日志内容。  */
	LogTime     string `json:"logTime,omitempty"`     /*  日志时间。  */
}
