package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListDataflowtaskApi
/* 查询资源池 ID 下，所有的数据流动任务详情
 */type HpfsListDataflowtaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListDataflowtaskApi(client *core.CtyunClient) *HpfsListDataflowtaskApi {
	return &HpfsListDataflowtaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-dataflowtask",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListDataflowtaskApi) Do(ctx context.Context, credential core.Credential, req *HpfsListDataflowtaskRequest) (*HpfsListDataflowtaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	if req.SfsUID != "" {
		ctReq.AddParam("sfsUID", req.SfsUID)
	}
	if req.DataflowID != "" {
		ctReq.AddParam("dataflowID", req.DataflowID)
	}
	if req.TaskStatus != "" {
		ctReq.AddParam("taskStatus", req.TaskStatus)
	}
	if req.TaskType != "" {
		ctReq.AddParam("taskType", req.TaskType)
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsListDataflowtaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListDataflowtaskRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池 ID  */
	AzName     string `json:"azName,omitempty"`     /*  多可用区下的可用区名字，不传为查询全部  */
	SfsUID     string `json:"sfsUID,omitempty"`     /*  并行文件UID  */
	DataflowID string `json:"dataflowID,omitempty"` /*  数据流动策略ID  */
	TaskStatus string `json:"taskStatus,omitempty"` /*  数据流动任务状态（creating/executing/completed/canceling/fail）  */
	TaskType   string `json:"taskType,omitempty"`   /*  数据流动任务类型（import_data/import_metadata/export_data）  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  每页包含的元素个数范围(1-50)，默认值为10  */
	PageNo     int32  `json:"pageNo,omitempty"`     /*  列表的分页页码，默认值为1  */
}

type HpfsListDataflowtaskResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述  */
	Description string                                 `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListDataflowtaskReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListDataflowtaskReturnObjResponse struct {
	TaskList     []*HpfsListDataflowtaskReturnObjTaskListResponse `json:"taskList"`     /*  返回的数据流动任务列表  */
	TotalCount   int32                                            `json:"totalCount"`   /*  指定条件下用户数据流动任务总数  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页码下查询回来的用户数据流动任务数  */
	PageSize     int32                                            `json:"pageSize"`     /*  每页包含的元素个数范围(1-50)  */
	PageNo       int32                                            `json:"pageNo"`       /*  列表的分页页码  */
}

type HpfsListDataflowtaskReturnObjTaskListResponse struct {
	RegionID        string `json:"regionID"`        /*  资源池id  */
	TaskID          string `json:"taskID"`          /*  数据流动任务id  */
	DataflowID      string `json:"dataflowID"`      /*  数据流动策略id  */
	SfsUID          string `json:"sfsUID"`          /*  并行文件唯一id  */
	SfsDirectory    string `json:"sfsDirectory"`    /*  并行文件下目录  */
	BucketName      string `json:"bucketName"`      /*  对象存储的桶名称  */
	BucketPrefix    string `json:"bucketPrefix"`    /*  对象存储桶的前缀  */
	TaskType        string `json:"taskType"`        /*  任务类型（import_data/import_metadata/export_data）  */
	TaskDescription string `json:"taskDescription"` /*  数据流动任务的描述  */
	TaskStatus      string `json:"taskStatus"`      /*  数据流动任务的状态，creating/executing/completed/canceling/fail。creating：任务创建中；executing：任务执行中；completed：任务已完成；canceling：任务取消中（可能是任务失败正在取消，也可能是策略删除任务正在取消）；fail：任务异常（异常原因可见failMsg，异常的任务不可恢复）  */
	CreateTime      int64  `json:"createTime"`      /*  数据流动任务创建时间  */
	StartTime       int64  `json:"startTime"`       /*  数据流动任务开始时间  */
	CompleteTime    int64  `json:"completeTime"`    /*  数据流动任务完成时间  */
	CancelTime      int64  `json:"cancelTime"`      /*  数据流动任务取消时间  */
	FailTime        int64  `json:"failTime"`        /*  数据流动任务异常发生时间  */
	FailMsg         string `json:"failMsg"`         /*  数据流动任务异常原因  */
	AzName          string `json:"azName"`          /*  多可用区下的可用区名字  */
}
