package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsCountQuotaSfsApi
/* 根据资源池ID查询用户在该资源池的弹性文件数量配额、已使用数量、剩余数量配额
 */type SfsCountQuotaSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsCountQuotaSfsApi(client *core.CtyunClient) *SfsCountQuotaSfsApi {
	return &SfsCountQuotaSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/count-quota-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsCountQuotaSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsCountQuotaSfsRequest) (*SfsCountQuotaSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsCountQuotaSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsCountQuotaSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
}

type SfsCountQuotaSfsResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                             `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                             `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsCountQuotaSfsReturnObjResponse `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsCountQuotaSfsReturnObjResponse struct {
	AvailableCount int32 `json:"availableCount"` /*  可用弹性文件配额数量  */
	UsedCount      int32 `json:"usedCount"`      /*  已用弹性文件配额数量  */
	QuotaCount     int32 `json:"quotaCount"`     /*  弹性文件配额总数量  */
}
