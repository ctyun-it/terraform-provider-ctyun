package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsStorageQuotaSfsApi
/* 查询用户在某地域的容量配额使用情况
 */type SfsStorageQuotaSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsStorageQuotaSfsApi(client *core.CtyunClient) *SfsStorageQuotaSfsApi {
	return &SfsStorageQuotaSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/storage-quota-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsStorageQuotaSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsStorageQuotaSfsRequest) (*SfsStorageQuotaSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsStorageQuotaSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsStorageQuotaSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
}

type SfsStorageQuotaSfsResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                               `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                               `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsStorageQuotaSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsStorageQuotaSfsReturnObjResponse struct {
	TotalQuota int32 `json:"totalQuota"` /*  总容量配额,，单位GB  */
	UsedQuota  int32 `json:"usedQuota"`  /*  已使⽤容量配额，单位GB  */
	AvlQuota   int32 `json:"avlQuota"`   /*  剩余容量配额，单位GB  */
}
