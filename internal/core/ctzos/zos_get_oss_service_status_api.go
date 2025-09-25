package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetOssServiceStatusApi
/* 查询对象存储是否开通。
 */type ZosGetOssServiceStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetOssServiceStatusApi(client *core.CtyunClient) *ZosGetOssServiceStatusApi {
	return &ZosGetOssServiceStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-oss-service-status",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetOssServiceStatusApi) Do(ctx context.Context, credential core.Credential, req *ZosGetOssServiceStatusRequest) (*ZosGetOssServiceStatusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetOssServiceStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetOssServiceStatusRequest struct {
	RegionID string /*  区域 ID  */
}

type ZosGetOssServiceStatusResponse struct {
	StatusCode  int64                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosGetOssServiceStatusReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	Description string                                   `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetOssServiceStatusReturnObjResponse struct {
	State string `json:"state,omitempty"` /*  开通状态，①true，已开通；②false，开通失败；③processing，开通中；④frozen，冻结；  */
}
