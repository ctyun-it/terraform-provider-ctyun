package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsAttachEbsApi
/* 将云硬盘挂载至某一云主机，支持非共享云硬盘和共享云硬盘的挂载。
 */type EbsAttachEbsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsAttachEbsApi(client *core.CtyunClient) *EbsAttachEbsApi {
	return &EbsAttachEbsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/attach-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsAttachEbsApi) Do(ctx context.Context, credential core.Credential, req *EbsAttachEbsRequest) (*EbsAttachEbsResponse, error) {
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
	var resp EbsAttachEbsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsAttachEbsRequest struct {
	DiskID     string  `json:"diskID,omitempty"`     /*  待挂载的云硬盘ID。  */
	RegionID   *string `json:"regionID,omitempty"`   /*  如本地语境支持保存regionID，那么建议传递。  */
	InstanceID string  `json:"instanceID,omitempty"` /*  待挂载云主机的ID。  */
}

type EbsAttachEbsResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)。  */
	Message     *string                          `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                          `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsAttachEbsReturnObjResponse   `json:"returnObj"`             /*  返回结构体。  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，请参考错误码。  */
	Error       *string                          `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码，请参考错误码。  */
	ErrorDetail *EbsAttachEbsErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。<br> 其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息。  */
}

type EbsAttachEbsReturnObjResponse struct {
	DiskJobID     *string `json:"diskJobID,omitempty"`     /*  异步任务ID，可通过公共接口/v4/job/info查询该jobID来查看异步任务最终执行结果（该参数即将被弃用，为提高兼容性，请尽量使用diskRequestID参数）。  */
	DiskRequestID *string `json:"diskRequestID,omitempty"` /*  异步任务ID，可通过公共接口/v4/job/info查询该jobID来查看异步任务最终执行结果。  */
}

type EbsAttachEbsErrorDetailResponse struct {
	BssErrCode       *string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中。  */
	BssErrMsg        *string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中。  */
	BssOrigErr       *string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息。  */
	BssErrPrefixHint *string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息。  */
}
