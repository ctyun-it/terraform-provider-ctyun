package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2GetTaskProgressDetailInfoApi
/* 查询在线迁移进度明细。
 */type Dcs2GetTaskProgressDetailInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2GetTaskProgressDetailInfoApi(client *core.CtyunClient) *Dcs2GetTaskProgressDetailInfoApi {
	return &Dcs2GetTaskProgressDetailInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/transfer/getTaskProgressDetailInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2GetTaskProgressDetailInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2GetTaskProgressDetailInfoRequest) (*Dcs2GetTaskProgressDetailInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("taskId", req.TaskId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2GetTaskProgressDetailInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2GetTaskProgressDetailInfoRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	TaskId   string /*  任务ID  */
}

type Dcs2GetTaskProgressDetailInfoResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                          `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2GetTaskProgressDetailInfoReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                          `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                          `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2GetTaskProgressDetailInfoReturnObjResponse struct {
	DataSyncCountInfo      *Dcs2GetTaskProgressDetailInfoReturnObjDataSyncCountInfoResponse        `json:"dataSyncCountInfo"`      /*  命令同步信息汇总  */
	SourceProgressInfoList []*Dcs2GetTaskProgressDetailInfoReturnObjSourceProgressInfoListResponse `json:"sourceProgressInfoList"` /*  同步进度信息  */
}

type Dcs2GetTaskProgressDetailInfoReturnObjDataSyncCountInfoResponse struct {
	ReadCount  int64   `json:"readCount,omitempty"`  /*  从源端读取的总命令数  */
	ReadOps    float64 `json:"readOps"`              /*  从源端读取的命令总OPS  */
	WriteCount int64   `json:"writeCount,omitempty"` /*  发送给目标的总命令数  */
	WriteOps   float64 `json:"writeOps"`             /*  发送给目标的命令总OPS  */
}

type Dcs2GetTaskProgressDetailInfoReturnObjSourceProgressInfoListResponse struct {
	Address       string  `json:"address,omitempty"`       /*  源端ip:port  */
	SyncPercent   float64 `json:"syncPercent"`             /*  全量数据的同步百分比（不包含增量）  */
	AofOffsetDiff int64   `json:"aofOffsetDiff,omitempty"` /*  AOF偏移差距（增量同步时有）  */
	State         string  `json:"state,omitempty"`         /*  同步阶段<li>unknown status<li>hand shaking<li>waiting bgsave<li>receiving rdb<li>syncing rdb<li>syncing aof（增量同步时有）  */
}
