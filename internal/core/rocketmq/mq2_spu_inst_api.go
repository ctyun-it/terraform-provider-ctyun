package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2SpuInstApi
/* 提供用户查询已订购的PaaS产品列表信息的能力 */
type Mq2SpuInstApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2SpuInstApi(client *core.CtyunClient) *Mq2SpuInstApi {
	return &Mq2SpuInstApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/tenants/ctyun/spuInst",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2SpuInstApi) Do(ctx context.Context, credential core.Credential, req *Mq2SpuInstRequest) (*Mq2SpuInstResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddHeader("accountId", req.AccountId)
	ctReq.AddParam("pageNow", req.PageNow)
	if req.PageSize != nil && *req.PageSize != "" {
		ctReq.AddParam("pageSize", *req.PageSize)
	}
	ctReq.AddParam("spuCode", req.SpuCode)
	if req.SpuInstId != nil && *req.SpuInstId != "" {
		ctReq.AddParam("spuInstId", *req.SpuInstId)
	}
	if req.SpuInstName != nil && *req.SpuInstName != "" {
		ctReq.AddParam("spuInstName", *req.SpuInstName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2SpuInstResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2SpuInstRequest struct {
	RegionId    string  `json:"regionId"`              /*  资源池编码  */
	AccountId   string  `json:"accountId"`             /*  账户  */
	PageNow     string  `json:"pageNow"`               /*  当前页码  */
	PageSize    *string `json:"pageSize,omitempty"`    /*  页面大小  */
	SpuCode     string  `json:"spuCode"`               /*  产品编码  */
	SpuInstId   *string `json:"spuInstId,omitempty"`   /*  实例ID  */
	SpuInstName *string `json:"spuInstName,omitempty"` /*  实例名称/别称  */
}

type Mq2SpuInstResponse struct {
	StatusCode *string `json:"statusCode"` /*  返回码
	取值范围：800 成功  */
	ReturnObj *Mq2SpuInstReturnObjResponse `json:"returnObj"` /*  返回对象  */
}

type Mq2SpuInstReturnObjResponse struct {
	PageNum  *int32                             `json:"pageNum"`  /*  当前页  */
	PageSize *int32                             `json:"pageSize"` /*  每页行数  */
	Size     *int32                             `json:"size"`     /*  页数  */
	List     []*Mq2SpuInstReturnObjListResponse `json:"list"`     /*  分页内容  */
}

type Mq2SpuInstReturnObjListResponse struct{}
