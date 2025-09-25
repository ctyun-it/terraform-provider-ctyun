package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsCountQuotaApi
/* 根据资源池ID查询用户在该资源池的并行文件数量配额总量、已使用数量、剩余数量
 */type HpfsCountQuotaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsCountQuotaApi(client *core.CtyunClient) *HpfsCountQuotaApi {
	return &HpfsCountQuotaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/count-quota-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsCountQuotaApi) Do(ctx context.Context, credential core.Credential, req *HpfsCountQuotaRequest) (*HpfsCountQuotaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsCountQuotaResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsCountQuotaRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
}

type HpfsCountQuotaResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                           `json:"message"`     /*  响应描述  */
	Description string                           `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsCountQuotaReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                           `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsCountQuotaReturnObjResponse struct {
	AvlCountQuota   int32 `json:"avlCountQuota"`   /*  可用并行文件配额数量  */
	UsedCountQuota  int32 `json:"usedCountQuota"`  /*  已用并行文件配额数量  */
	TotalCountQuota int32 `json:"totalCountQuota"` /*  并行文件配额总数量  */
}
