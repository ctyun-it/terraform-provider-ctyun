package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsEcsSnapshotTaskListApi
/* 该接口提供用户查询云主机快照任务列表的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsEcsSnapshotTaskListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsEcsSnapshotTaskListApi(client *core.CtyunClient) *CtecsEcsSnapshotTaskListApi {
	return &CtecsEcsSnapshotTaskListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot-task/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsEcsSnapshotTaskListApi) Do(ctx context.Context, credential core.Credential, req *CtecsEcsSnapshotTaskListRequest) (*CtecsEcsSnapshotTaskListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsEcsSnapshotTaskListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsEcsSnapshotTaskListRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	TaskID       string `json:"taskID,omitempty"`       /*  任务ID  */
	TaskType     string `json:"taskType,omitempty"`     /*  任务类型，取值范围：<br />create：生成快照任务，<br />restore：恢复快照数据任务，<br />apply：申请云主机任务，<br />delete：删除快照任务<br />注：不传默认全部  */
	InstanceID   string `json:"instanceID,omitempty"`   /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	InstanceName string `json:"instanceName,omitempty"` /*  云主机名称  */
	SnapshotID   string `json:"snapshotID,omitempty"`   /*  云主机快照ID，<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8349&data=87">查询云主机快照列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8352&data=87">创建云主机快照</a>  */
	SnapshotName string `json:"snapshotName,omitempty"` /*  快照名称  */
	StrategyID   string `json:"strategyID,omitempty"`   /*  快照策略ID  */
	QueryContent string `json:"queryContent,omitempty"` /*  模糊查询，可匹配查询字段：任务ID、快照ID、快照名称、云主机ID、云主机名称  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsEcsSnapshotTaskListResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                     `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                     `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                     `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsEcsSnapshotTaskListReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsEcsSnapshotTaskListReturnObjResponse struct {
	CurrentCount int32                                               `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                               `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                               `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsEcsSnapshotTaskListReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsEcsSnapshotTaskListReturnObjResultsResponse struct {
	TaskID         string `json:"taskID,omitempty"`         /*  任务ID  */
	TaskStatus     string `json:"taskStatus,omitempty"`     /*  任务状态  */
	TaskType       string `json:"taskType,omitempty"`       /*  任务类型  */
	InstanceID     string `json:"instanceID,omitempty"`     /*  云主机ID  */
	InstanceName   string `json:"instanceName,omitempty"`   /*  云主机名称  */
	SnapshotID     string `json:"snapshotID,omitempty"`     /*  快照ID  */
	SnapshotName   string `json:"snapshotName,omitempty"`   /*  快照名称  */
	StrategyID     string `json:"strategyID,omitempty"`     /*  策略的ID，非策略触发执行时为空  */
	StrategyName   string `json:"strategyName,omitempty"`   /*  策略的名称，非策略触发执行时为空  */
	StartTime      string `json:"startTime,omitempty"`      /*  开始时间  */
	FinishTime     string `json:"finishTime,omitempty"`     /*  完成时间  */
	TaskDetailDesc string `json:"taskDetailDesc,omitempty"` /*  任务详情  */
}
