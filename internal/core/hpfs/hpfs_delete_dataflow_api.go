package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsDeleteDataflowApi
/* 删除数据流动策略
 */type HpfsDeleteDataflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsDeleteDataflowApi(client *core.CtyunClient) *HpfsDeleteDataflowApi {
	return &HpfsDeleteDataflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/delete-dataflow",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsDeleteDataflowApi) Do(ctx context.Context, credential core.Credential, req *HpfsDeleteDataflowRequest) (*HpfsDeleteDataflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsDeleteDataflowRequest
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
	var resp HpfsDeleteDataflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsDeleteDataflowRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池 ID  */
	DataflowID string `json:"dataflowID,omitempty"` /*  数据流动策略ID  */
}

type HpfsDeleteDataflowResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                               `json:"message"`     /*  响应描述  */
	Description string                               `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsDeleteDataflowReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsDeleteDataflowReturnObjResponse struct {
	RegionID string `json:"regionID"` /*  资源所属资源池 ID  */
}
