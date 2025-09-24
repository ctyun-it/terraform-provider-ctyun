package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryJobInfoV41Api
/* 查看异步任务job任务状态等
 */type CtecsQueryJobInfoV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryJobInfoV41Api(client *core.CtyunClient) *CtecsQueryJobInfoV41Api {
	return &CtecsQueryJobInfoV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/job/info",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryJobInfoV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryJobInfoV41Request) (*CtecsQueryJobInfoV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("jobID", req.JobID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryJobInfoV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryJobInfoV41Request struct {
	RegionID string /*  资源池ID  */
	JobID    string /*  异步任务ID  */
}

type CtecsQueryJobInfoV41Response struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                 `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                 `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQueryJobInfoV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                 `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryJobInfoV41ReturnObjResponse struct {
	JobID      string                                       `json:"jobID,omitempty"`     /*  异步任务ID  */
	Status     int32                                        `json:"status,omitempty"`    /*  任务状态 (0:执行中 1:执行成功 2:执行失败)  */
	JobStatus  string                                       `json:"jobStatus,omitempty"` /*  job任务状态(executing:执行中, success:执行成功, fail:执行失败)  */
	ResourceId string                                       `json:"resourceId,omitempty"`
	Fields     *CtecsQueryJobInfoV41ReturnObjFieldsResponse `json:"fields"` /*  任务信息  */
	ID         string                                       `json:"ID"`     /*  资源ID  */
}

type CtecsQueryJobInfoV41ReturnObjFieldsResponse struct {
	TaskName string `json:"taskName,omitempty"` /*  任务名  */
}
