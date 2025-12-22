package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageDeleteImageImportTaskApi
/* 删除 1 个失败的使用镜像文件来创建私有镜像的任务<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：<br />1. 此接口在多可用区资源池已弃用，目前仍可调用。若您欲删除在多可用区资源池中使用镜像文件创建的、镜像状态为 error 的私有镜像，则可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4766&data=89&isNormal=1&vid=83" target="_blank">删除私有镜像</a>接口来删除 1 份私有镜像<br />2. 在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageDeleteImageImportTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageDeleteImageImportTaskApi(client *core.CtyunClient) *CtimageDeleteImageImportTaskApi {
	return &CtimageDeleteImageImportTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/delete-import-task",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageDeleteImageImportTaskApi) Do(ctx context.Context, credential core.Credential, req *CtimageDeleteImageImportTaskRequest) (*CtimageDeleteImageImportTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("taskID", req.TaskID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageDeleteImageImportTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageDeleteImageImportTaskRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
	TaskID   string `json:"taskID,omitempty"`   /*  任务 ID。在非多可用区资源池中，可通过镜像详细信息中的 taskID 参数获取任务 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息。注意：您需自行确保所指定的任务对您而言是存在的，且对应的镜像是镜像状态为 error 的私有镜像。否则，即使接口请求成功，任务也不会被删除  */
}

type CtimageDeleteImageImportTaskResponse struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                                         `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                         `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageDeleteImageImportTaskReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageDeleteImageImportTaskReturnObjResponse struct{}
