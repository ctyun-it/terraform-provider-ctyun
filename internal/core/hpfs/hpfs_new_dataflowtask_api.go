package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsNewDataflowtaskApi
/* 创建数据流动任务
 */type HpfsNewDataflowtaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsNewDataflowtaskApi(client *core.CtyunClient) *HpfsNewDataflowtaskApi {
	return &HpfsNewDataflowtaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/new-dataflowtask",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsNewDataflowtaskApi) Do(ctx context.Context, credential core.Credential, req *HpfsNewDataflowtaskRequest) (*HpfsNewDataflowtaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsNewDataflowtaskRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsNewDataflowtaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsNewDataflowtaskRequest struct {
	RegionID        string `json:"regionID,omitempty"`        /*  资源池 ID  */
	SfsUID          string `json:"sfsUID,omitempty"`          /*  并行文件唯一ID  */
	DataflowID      string `json:"dataflowID,omitempty"`      /*  数据流动策略ID  */
	TaskType        string `json:"taskType,omitempty"`        /*  数据流动任务类型（目前支持import_data/import_metadata/export_data）  */
	TaskDescription string `json:"taskDescription,omitempty"` /*  数据流动策略的描述，最高支持128字符  */
}

type HpfsNewDataflowtaskResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                `json:"message"`     /*  响应描述  */
	Description string                                `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsNewDataflowtaskReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsNewDataflowtaskReturnObjResponse struct {
	RegionID  string                                           `json:"regionID"`  /*  资源所属资源池 ID  */
	Resources []*HpfsNewDataflowtaskReturnObjResourcesResponse `json:"resources"` /*  资源明细  */
}

type HpfsNewDataflowtaskReturnObjResourcesResponse struct {
	DataflowID string `json:"dataflowID"` /*  数据流动策略ID  */
	SfsUID     string `json:"sfsUID"`     /*  并行文件ID  */
	TaskID     string `json:"taskID"`     /*  数据流动任务ID  */
}
