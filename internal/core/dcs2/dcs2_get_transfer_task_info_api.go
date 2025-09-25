package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2GetTransferTaskInfoApi
/* 查询迁移任务详情。
 */type Dcs2GetTransferTaskInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2GetTransferTaskInfoApi(client *core.CtyunClient) *Dcs2GetTransferTaskInfoApi {
	return &Dcs2GetTransferTaskInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/transfer/getTaskInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2GetTransferTaskInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2GetTransferTaskInfoRequest) (*Dcs2GetTransferTaskInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("taskId", req.TaskId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2GetTransferTaskInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2GetTransferTaskInfoRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	TaskId   string /*  任务ID  */
}

type Dcs2GetTransferTaskInfoResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2GetTransferTaskInfoReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2GetTransferTaskInfoReturnObjResponse struct {
	UserId          string                                          `json:"userId,omitempty"`          /*  用户ID  */
	TenantId        string                                          `json:"tenantId,omitempty"`        /*  租户ID  */
	TaskId          string                                          `json:"taskId,omitempty"`          /*  任务ID  */
	SourceSpuInstId string                                          `json:"sourceSpuInstId,omitempty"` /*  源库实例id  */
	TargetSpuInstId string                                          `json:"targetSpuInstId,omitempty"` /*  目标库实例id  */
	SyncMode        int32                                           `json:"syncMode,omitempty"`        /*  同步模式<li>1：全量同步+增量同步<li>2：全量同步  */
	RawType         int32                                           `json:"type,omitempty"`            /*  任务类型（1：数据迁移任务，其他：未知；本接口该字段固定返回1）  */
	ConflictMode    int32                                           `json:"conflictMode,omitempty"`    /*  数据冲突时的处理办法<li>1：中断迁移<li>2：跳过目标key，继续执行<li>3：覆盖目标key，继续执行  */
	Status          int32                                           `json:"status,omitempty"`          /*  任务状态<li>0：初始态<li>1：运行中<li>2：成功<li>3：失败  */
	RunStep         int32                                           `json:"runStep,omitempty"`         /*  迁移进度<li>1：全量同步中<li>2：增量同步中  */
	CreateTime      int64                                           `json:"createTime,omitempty"`      /*  创建时间(秒)  */
	CompleteTime    int64                                           `json:"completeTime,omitempty"`    /*  结束时间（秒，-1表示时间未知）  */
	Detail          *Dcs2GetTransferTaskInfoReturnObjDetailResponse `json:"detail"`                    /*  详情  */
}

type Dcs2GetTransferTaskInfoReturnObjDetailResponse struct {
	SourceDbInfo *Dcs2GetTransferTaskInfoReturnObjDetailSourceDbInfoResponse `json:"sourceDbInfo"` /*  源库信息  */
	TargetDbInfo *Dcs2GetTransferTaskInfoReturnObjDetailTargetDbInfoResponse `json:"targetDbInfo"` /*  目标库信息  */
}

type Dcs2GetTransferTaskInfoReturnObjDetailSourceDbInfoResponse struct {
	SpuInstId       string `json:"spuInstId,omitempty"`    /*  实例ID  */
	IpAddr          string `json:"ipAddr,omitempty"`       /*  连接地址  */
	OriginalCluster *bool  `json:"originalCluster"`        /*  是否是原生cluster集群  */
	AccountName     string `json:"accountName,omitempty"`  /*  数据库账号  */
	Password        string `json:"password,omitempty"`     /*  数据库密码  */
	InstanceName    string `json:"instanceName,omitempty"` /*  实例名称  */
}

type Dcs2GetTransferTaskInfoReturnObjDetailTargetDbInfoResponse struct {
	SpuInstId       string `json:"spuInstId,omitempty"`    /*  实例ID  */
	IpAddr          string `json:"ipAddr,omitempty"`       /*  连接地址  */
	OriginalCluster *bool  `json:"originalCluster"`        /*  是否是原生cluster集群  */
	AccountName     string `json:"accountName,omitempty"`  /*  数据库账号  */
	Password        string `json:"password,omitempty"`     /*  数据库密码  */
	InstanceName    string `json:"instanceName,omitempty"` /*  实例名称  */
}
