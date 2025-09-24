package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsCapacityQuotaApi
/* 根据资源池ID查询用户在某地域的并行文件容量配额总量，已使用容量，剩余容量
 */type HpfsCapacityQuotaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsCapacityQuotaApi(client *core.CtyunClient) *HpfsCapacityQuotaApi {
	return &HpfsCapacityQuotaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/capacity-quota-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsCapacityQuotaApi) Do(ctx context.Context, credential core.Credential, req *HpfsCapacityQuotaRequest) (*HpfsCapacityQuotaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsCapacityQuotaResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsCapacityQuotaRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
}

type HpfsCapacityQuotaResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                              `json:"message"`     /*  响应描述  */
	Description string                              `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsCapacityQuotaReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsCapacityQuotaReturnObjResponse struct {
	TotalCapacityQuota int32 `json:"totalCapacityQuota"` /*  总容量配额，单位GB，-1表示无限制  */
	UsedCapacityQuota  int32 `json:"usedCapacityQuota"`  /*  已使⽤容量配额，单位GB  */
	AvlCapacityQuota   int32 `json:"avlCapacityQuota"`   /*  剩余容量配额，单位GB，-1表示无限制  */
}
