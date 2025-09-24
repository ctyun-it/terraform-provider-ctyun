package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetFragmentNumApi
/* 查询桶内的碎片数量。
 */type ZosGetFragmentNumApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetFragmentNumApi(client *core.CtyunClient) *ZosGetFragmentNumApi {
	return &ZosGetFragmentNumApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-fragment-num",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetFragmentNumApi) Do(ctx context.Context, credential core.Credential, req *ZosGetFragmentNumRequest) (*ZosGetFragmentNumResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetFragmentNumResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetFragmentNumRequest struct {
	Bucket   string /*  存储桶名  */
	RegionID string /*  区域ID  */
}

type ZosGetFragmentNumResponse struct {
	StatusCode  int64                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                              `json:"message,omitempty"`     /*  状态描述  */
	Description string                              `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetFragmentNumReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                              `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                              `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetFragmentNumReturnObjResponse struct {
	FragmentNum int64 `json:"fragmentNum,omitempty"` /*  碎片数量  */
}
