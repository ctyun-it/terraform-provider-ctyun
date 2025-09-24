package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsUpdateDataflowApi
/* 更新数据流动策略
 */type HpfsUpdateDataflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsUpdateDataflowApi(client *core.CtyunClient) *HpfsUpdateDataflowApi {
	return &HpfsUpdateDataflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/update-dataflow",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsUpdateDataflowApi) Do(ctx context.Context, credential core.Credential, req *HpfsUpdateDataflowRequest) (*HpfsUpdateDataflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsUpdateDataflowRequest
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
	var resp HpfsUpdateDataflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsUpdateDataflowRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	DataflowID          string `json:"dataflowID,omitempty"`          /*  数据流动策略ID  */
	AutoSync            *bool  `json:"autoSync"`                      /*  是否打开自动同步  */
	DataflowDescription string `json:"dataflowDescription,omitempty"` /*  数据流动策略的描述  */
}

type HpfsUpdateDataflowResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                               `json:"message"`     /*  响应描述  */
	Description string                               `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsUpdateDataflowReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsUpdateDataflowReturnObjResponse struct {
	RegionID string `json:"regionID"` /*  资源所属资源池 ID  */
}
