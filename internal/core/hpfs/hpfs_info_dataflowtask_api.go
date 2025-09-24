package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsInfoDataflowtaskApi
/* 查询资源池 ID 下，指定数据流动任务详情
 */type HpfsInfoDataflowtaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsInfoDataflowtaskApi(client *core.CtyunClient) *HpfsInfoDataflowtaskApi {
	return &HpfsInfoDataflowtaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/info-dataflowtask",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsInfoDataflowtaskApi) Do(ctx context.Context, credential core.Credential, req *HpfsInfoDataflowtaskRequest) (*HpfsInfoDataflowtaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("taskID", req.TaskID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsInfoDataflowtaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsInfoDataflowtaskRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	TaskID   string `json:"taskID,omitempty"`   /*  数据流动任务ID  */
}

type HpfsInfoDataflowtaskResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述  */
	Description string                                 `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsInfoDataflowtaskReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsInfoDataflowtaskReturnObjResponse struct {
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
