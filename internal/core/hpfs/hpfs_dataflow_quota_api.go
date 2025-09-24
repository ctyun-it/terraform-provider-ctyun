package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsDataflowQuotaApi
/* 查询指定并行文件下数据流动策略配额的总量，已使用数量，剩余数量，其中error状态的数据流动策略不占用配额
 */type HpfsDataflowQuotaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsDataflowQuotaApi(client *core.CtyunClient) *HpfsDataflowQuotaApi {
	return &HpfsDataflowQuotaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/dataflow-quota-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsDataflowQuotaApi) Do(ctx context.Context, credential core.Credential, req *HpfsDataflowQuotaRequest) (*HpfsDataflowQuotaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsDataflowQuotaResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsDataflowQuotaRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  并行文件唯一ID  */
}

type HpfsDataflowQuotaResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                              `json:"message"`     /*  响应描述  */
	Description string                              `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsDataflowQuotaReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsDataflowQuotaReturnObjResponse struct {
	TotalDataflowQuota int32 `json:"totalDataflowQuota"` /*  文件系统数据流动策略配额总数  */
	UsedDataflowQuota  int32 `json:"usedDataflowQuota"`  /*  文件系统已用数据流动策略配额数  */
	AvlDataflowQuota   int32 `json:"avlDataflowQuota"`   /*  文件系统剩余可用数据流动策略配额数  */
}
